# Deploy guide — ChatServer on VPS (k3s + GitHub Actions)

End-to-end walkthrough: get a fresh Debian 13 VPS serving `https://zloypes.com` with CI/CD from `master`. The pipeline is: push to master → GitHub Actions builds and pushes an image to GHCR → Actions SSHes to the VPS → `kubectl set image` rolls the deployment.

---

## 0. Prerequisites checklist

Before you start, confirm each of these:

- [ ] VPS is Debian 13 (trixie), with root SSH access. Verify from your laptop:
  ```bash
  ssh root@62.83.87.38 "cat /etc/os-release | head -2"
  ```
- [ ] DNS A record `zloypes.com → 62.83.87.38` is in place. Verify:
  ```bash
  dig +short zloypes.com      # should print 62.83.87.38
  ```
  If it doesn't, fix DNS at your registrar first and wait for propagation (usually minutes, up to 24h).
- [ ] You have push access to <https://github.com/VladimirTakeda/ChatServer>.
- [ ] The VPS has at least ~4 GB RAM. HPA is set to 3–10 replicas; add Postgres/Redis/MinIO on top of k3s overhead and 2 GB is tight.

---

## 1. Bootstrap the VPS (one-time)

SSH into the VPS as root and run the setup script. It installs k3s (with Traefik and ServiceLB disabled), Helm, ingress-nginx, cert-manager, the Let's Encrypt `ClusterIssuer`, creates the `enrollment` namespace, and applies `k8s/`.

```bash
ssh root@62.83.87.38
curl -sSL https://raw.githubusercontent.com/VladimirTakeda/ChatServer/master/scripts/vps-setup.sh | bash
```

Expected duration: 3–5 minutes. You should see green `==>` log lines for each phase.

**If something fails mid-run** — the script is idempotent-ish: re-run it with the same command. It will no-op the parts that already succeeded (k3s is already installed, etc.).

At the end it prints a "VPS bootstrap complete" banner with the next steps. The sections below expand them.

### Sanity-check after the script

```bash
kubectl get nodes                                                    # Ready
kubectl -n enrollment get pods                                       # all Running (eventually)
kubectl -n ingress-nginx get svc ingress-nginx-controller            # EXTERNAL-IP = 62.83.87.38
kubectl get clusterissuer letsencrypt-prod                           # Ready=True
```

If `ingress-nginx-controller` has no EXTERNAL-IP, k3s's built-in `ServiceLB` is probably still running and conflicting. The script disables it — if you installed k3s before, re-install with `--disable traefik --disable servicelb` or delete the stale LoadBalancer manually.

---

## 2. Create a deploy SSH key for CI

GitHub Actions needs its own key — don't reuse your personal one.

On the VPS:

```bash
ssh-keygen -t ed25519 -f /root/.ssh/gha_deploy -N ''
cat /root/.ssh/gha_deploy.pub >> /root/.ssh/authorized_keys
chmod 600 /root/.ssh/authorized_keys
```

Now print the **private** key — you'll paste this into GitHub next:

```bash
cat /root/.ssh/gha_deploy
```

Copy the entire block, including `-----BEGIN OPENSSH PRIVATE KEY-----` / `-----END OPENSSH PRIVATE KEY-----` lines.

Quickly verify the key works from your laptop (optional but worth it):

```bash
# From your laptop, paste the same private key into a temp file:
pbpaste > /tmp/test_deploy_key && chmod 600 /tmp/test_deploy_key
ssh -i /tmp/test_deploy_key root@62.83.87.38 "echo OK"
rm /tmp/test_deploy_key
```

Should print `OK`. If it asks for a password, something's wrong — fix before continuing.

---

## 3. Add GitHub Actions Secrets

Go to **<https://github.com/VladimirTakeda/ChatServer/settings/secrets/actions>** and add three repository secrets:

| Name | Value |
|---|---|
| `VPS_HOST` | `62.83.87.38` |
| `VPS_USER` | `root` |
| `VPS_SSH_KEY` | Contents of `/root/.ssh/gha_deploy` (private key, full block) |

`GITHUB_TOKEN` is provided automatically — no action needed for GHCR auth.

---

## 4. First deploy — push to master

From your laptop:

```bash
cd ~/Desktop/GoProjects/ChatServer
git add scripts/ .github/ docs/ k8s/
git commit -m "ci: add CI/CD to VPS via GHCR and SSH deploy"
git push origin master
```

Watch the workflow:

```bash
gh run watch
# or open https://github.com/VladimirTakeda/ChatServer/actions
```

The workflow has two jobs:

1. **build-and-push** — builds `./docker/backend_docker/Dockerfile`, pushes `ghcr.io/vladimirtakeda/chatserver:latest` and `:sha-<short>` to GHCR.
2. **deploy** — SSHes to the VPS and runs `kubectl set image deployment/application application=ghcr.io/vladimirtakeda/chatserver:sha-<short> -n enrollment`.

Expected total runtime: 2–4 minutes, longer on the first run (no build cache yet).

---

## 5. Make the GHCR package public

After the **first** successful push, the package appears here:

<https://github.com/users/VladimirTakeda/packages/container/chatserver>

Click **Package settings** (right sidebar) → scroll to **Danger Zone** → **Change visibility** → **Public** → confirm.

Without this step, k3s can't pull the image (it has no GHCR credentials). You'll see `ImagePullBackOff` on the pods until you either make the package public or create a pull secret.

**If you want to keep the package private**, create a pull secret on the VPS instead:

```bash
# Generate a PAT with 'read:packages' scope at:
# https://github.com/settings/tokens/new?scopes=read:packages
export GHCR_PAT='<your PAT>'

kubectl -n enrollment create secret docker-registry ghcr \
  --docker-server=ghcr.io \
  --docker-username=VladimirTakeda \
  --docker-password="$GHCR_PAT" \
  --docker-email=ooornament@gmail.com
```

Then add to `k8s/application.yaml` under `spec.template.spec`:

```yaml
      imagePullSecrets:
        - name: ghcr
```

Re-apply: `kubectl apply -f k8s/application.yaml -n enrollment`.

---

## 6. Wait for the TLS cert

cert-manager will request a Let's Encrypt cert via HTTP-01 on first ingress reconcile. Watch it:

```bash
kubectl -n enrollment get certificate -w
# status goes:  (none) → False → True
```

Describe it if it's stuck:

```bash
kubectl -n enrollment describe certificate zloypes-com-tls
kubectl -n enrollment describe order
kubectl -n enrollment describe challenge
```

Common failure causes:
- DNS doesn't actually point to this VPS (Let's Encrypt can't reach `/.well-known/acme-challenge/...`).
- Port 80 blocked by cloud firewall (check your VPS provider's firewall, separate from `ufw`).
- Rate limit hit — Let's Encrypt limits are 5 duplicate certs per week; if you're testing, use the staging ClusterIssuer first (edit `vps-setup.sh` to point at `acme-staging-v02.api.letsencrypt.org`).

Once `Ready=True`, open **<https://zloypes.com>** in a browser. You should see whatever the app's root endpoint returns.

---

## 7. Everyday ops

### Trigger a deploy

Any push to `master` triggers CI. You can also kick a deploy manually without pushing code:

- GitHub UI: Actions tab → "Build and deploy" → Run workflow → master
- CLI: `gh workflow run "Build and deploy" --ref master`

### Check what's running

```bash
ssh root@62.83.87.38
kubectl -n enrollment get pods -o wide
kubectl -n enrollment get deployment application -o jsonpath='{.spec.template.spec.containers[0].image}'
kubectl -n enrollment logs -f deployment/application
```

### Rollback

Last-known-good:

```bash
kubectl -n enrollment rollout undo deployment/application
kubectl -n enrollment rollout status deployment/application
```

Pin to a specific previous image (tags are `sha-<12 chars of SHA>`):

```bash
kubectl -n enrollment set image deployment/application \
  application=ghcr.io/vladimirtakeda/chatserver:sha-abcdef123456
```

See deployment history:

```bash
kubectl -n enrollment rollout history deployment/application
```

### Scale

HPA is set to 3–10 based on CPU. To temporarily override:

```bash
kubectl -n enrollment scale deployment/application --replicas=1   # HPA will push back up
# to actually override, delete the HPA first:
kubectl -n enrollment delete hpa application-hpa
```

---

## 8. Troubleshooting

### `ImagePullBackOff` on application pods

```bash
kubectl -n enrollment describe pod <pod-name> | tail -20
```

- `unauthorized`: GHCR package is private and no pull secret — see §5.
- `manifest unknown`: CI pushed but the tag doesn't exist. Check `kubectl get deployment application -o yaml | grep image:` matches a tag that actually exists at <https://github.com/VladimirTakeda/ChatServer/pkgs/container/chatserver>.

### Pods `CrashLoopBackOff`

```bash
kubectl -n enrollment logs deployment/application --previous
```

Likely causes right now, flagged in `scripts/vps-setup.sh` output:
- **Redis pub/sub not resolving** — the configmap has `REDIS_PUBSUB_HOST: "redis-pubsub"` but the Service is named `pubsub`. Fix:
  ```bash
  kubectl -n enrollment set env configmap/app-config REDIS_PUBSUB_HOST=pubsub
  kubectl -n enrollment rollout restart deployment/application
  ```
  Or edit `k8s/configmap.yaml` and re-apply.
- **Postgres migrations failing** — migrations run on startup from `postgres.RunMigrations("file://migrations", db)`. Check that `migrations/` was included in the Docker image (it should be, via `COPY . ./`).

### Ingress returns 404 for all paths

```bash
kubectl -n enrollment get ingress
kubectl -n ingress-nginx logs deployment/ingress-nginx-controller
```

- Make sure `ingressClassName: nginx` is present (it is, in our `k8s/ingress.yaml`).
- Confirm the Service `application` exists and selects the running pods:
  ```bash
  kubectl -n enrollment get svc application -o wide
  kubectl -n enrollment get endpoints application
  ```

### WebSocket drops after 60s

Shouldn't happen — our ingress annotations set `proxy-read-timeout: 3600`. If it does, check that ingress-nginx is the active controller (`kubectl get ingressclass`) and not some leftover Traefik from a k3s install that predated our setup script.

### CI deploy step fails with `permission denied (publickey)`

- `VPS_SSH_KEY` secret is wrong or truncated. Re-copy the private key from `/root/.ssh/gha_deploy` making sure both `BEGIN` and `END` markers are included.
- `VPS_USER` doesn't match the user `authorized_keys` lives under. If you used `root`, the key must be in `/root/.ssh/authorized_keys`.

### CI deploy step fails with `error: no server found`

The SSH worked, but `kubectl` on the VPS can't reach the k3s API. Make sure `/etc/rancher/k3s/k3s.yaml` exists (the script's `--write-kubeconfig-mode 644` makes it readable) and `KUBECONFIG` is exported in `/root/.bashrc`.

---

## 9. Known issues to fix soon

- **`k8s/configmap.yaml` has `REDIS_PUBSUB_HOST: "redis-pubsub"`** but the Service in `k8s/redis-pubsub.yaml:38` is named `pubsub`. The app won't reach pub/sub until this is corrected. See §8 for a hot-fix; the lasting fix is to edit the configmap and commit.
- **`docker/backend_docker/Dockerfile` uses `COPY ../../go.mod ./`** — paths above the build context aren't valid Docker. If the CI build fails with `forbidden path outside the build context`, change to:
  ```dockerfile
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . ./
  RUN go build -v -o chatServer ./cmd/main.go
  ```
- **`k8s/secrets.yaml` holds base64-encoded plaintext passwords** committed to git. Any time you want to tighten this: delete the file from the next `kubectl apply`, rotate the passwords, and create the secret out-of-band:
  ```bash
  kubectl -n enrollment create secret generic app-secrets \
    --from-literal=POSTGRES_PASSWORD='...' \
    --from-literal=REDIS_PASSWORD='...' \
    --from-literal=REDIS_PUBSUB_PASSWORD='...' \
    --from-literal=MINIO_ROOT_USER='...' \
    --from-literal=MINIO_ROOT_PASSWORD='...'
  ```

---

## 10. Teardown (if you ever need to start over)

On the VPS:

```bash
# Nuke the whole app
kubectl delete namespace enrollment

# Nuke k3s entirely
/usr/local/bin/k3s-uninstall.sh
rm -rf /etc/rancher /var/lib/rancher /opt/chatserver
```

Then re-run `scripts/vps-setup.sh` from §1.
