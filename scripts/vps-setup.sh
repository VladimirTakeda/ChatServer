#!/usr/bin/env bash
# One-time VPS bootstrap for ChatServer on Debian 13 (trixie).
# Run as root. Idempotent-ish — safe to re-run; it will no-op steps that are already done.
#
#   curl -sSL https://raw.githubusercontent.com/VladimirTakeda/ChatServer/master/scripts/vps-setup.sh | bash
# or, after cloning the repo on the VPS:
#   sudo bash scripts/vps-setup.sh

set -euo pipefail

DOMAIN="${DOMAIN:-zloypes.com}"
LE_EMAIL="${LE_EMAIL:-ooornament@gmail.com}"
NAMESPACE="${NAMESPACE:-enrollment}"
REPO_URL="${REPO_URL:-https://github.com/VladimirTakeda/ChatServer.git}"
REPO_DIR="${REPO_DIR:-/opt/chatserver}"

if [[ $EUID -ne 0 ]]; then
  echo "Run as root (sudo bash $0)." >&2
  exit 1
fi

log() { printf '\n\033[1;32m==>\033[0m %s\n' "$*"; }

log "Installing base packages"
export DEBIAN_FRONTEND=noninteractive
apt-get update -y
apt-get install -y curl ca-certificates git jq apt-transport-https gnupg

log "Cloning/updating repo at ${REPO_DIR}"
if [[ -d "${REPO_DIR}/.git" ]]; then
  git -C "${REPO_DIR}" fetch --all --quiet
  git -C "${REPO_DIR}" reset --hard origin/master
else
  git clone --depth 1 "${REPO_URL}" "${REPO_DIR}"
fi

log "Installing k3s (Traefik disabled — we use ingress-nginx)"
if ! command -v k3s >/dev/null 2>&1; then
  curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--disable traefik --disable servicelb --write-kubeconfig-mode 644" sh -
else
  echo "k3s already installed, skipping."
fi

# Make kubectl + KUBECONFIG usable for the current shell and for root's profile.
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
if ! grep -q 'KUBECONFIG=/etc/rancher/k3s/k3s.yaml' /root/.bashrc; then
  echo 'export KUBECONFIG=/etc/rancher/k3s/k3s.yaml' >> /root/.bashrc
fi
ln -sf /usr/local/bin/k3s /usr/local/bin/kubectl

log "Waiting for k3s node to be Ready"
for _ in $(seq 1 60); do
  if kubectl get nodes 2>/dev/null | grep -q ' Ready '; then break; fi
  sleep 2
done
kubectl get nodes

log "Installing Helm (official installer)"
if ! command -v helm >/dev/null 2>&1; then
  curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 -o /tmp/get-helm-3.sh
  chmod +x /tmp/get-helm-3.sh
  /tmp/get-helm-3.sh
  rm -f /tmp/get-helm-3.sh
fi
helm version

log "Installing ingress-nginx (hostNetwork mode — single-node VPS)"
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx >/dev/null 2>&1 || true
helm repo update >/dev/null
helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx --create-namespace \
  --set controller.hostNetwork=true \
  --set controller.dnsPolicy=ClusterFirstWithHostNet \
  --set controller.kind=DaemonSet \
  --set controller.service.type=ClusterIP \
  --set controller.ingressClassResource.default=true \
  --wait --timeout 5m

log "Installing cert-manager"
helm repo add jetstack https://charts.jetstack.io >/dev/null 2>&1 || true
helm repo update >/dev/null
helm upgrade --install cert-manager jetstack/cert-manager \
  --namespace cert-manager --create-namespace \
  --set crds.enabled=true \
  --wait --timeout 5m

log "Creating Let's Encrypt ClusterIssuer"
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: ${LE_EMAIL}
    privateKeySecretRef:
      name: letsencrypt-prod-account-key
    solvers:
      - http01:
          ingress:
            ingressClassName: nginx
EOF

log "Creating namespace ${NAMESPACE}"
kubectl create namespace "${NAMESPACE}" --dry-run=client -o yaml | kubectl apply -f -

log "Applying ChatServer manifests from ${REPO_DIR}/k8s"
kubectl apply -f "${REPO_DIR}/k8s" -n "${NAMESPACE}"

log "Current pods"
kubectl get pods -n "${NAMESPACE}"

log "Opening firewall (ufw if present)"
if command -v ufw >/dev/null 2>&1; then
  ufw allow 22/tcp || true
  ufw allow 80/tcp || true
  ufw allow 443/tcp || true
fi

cat <<EOF

================================================================
VPS bootstrap complete.

Next steps (one-time):

1. Create an SSH key for GitHub Actions to log in with:
     ssh-keygen -t ed25519 -f /root/.ssh/gha_deploy -N ''
     cat /root/.ssh/gha_deploy.pub >> /root/.ssh/authorized_keys
     cat /root/.ssh/gha_deploy       # <-- PRIVATE KEY, copy into GitHub Secrets

2. In GitHub → repo settings → Secrets and variables → Actions, add:
     VPS_HOST    = $(curl -s ifconfig.me || echo '62.83.87.38')
     VPS_USER    = root
     VPS_SSH_KEY = <contents of /root/.ssh/gha_deploy>

3. Make the GHCR package 'chatserver' public after the first CI push:
     https://github.com/users/VladimirTakeda/packages/container/chatserver/settings
     → Change visibility → Public
   (Otherwise you need an imagePullSecret — see docs/deploy.md.)

4. Point DNS: ${DOMAIN} A → this VPS's public IP.
   Current A record resolves to: $(getent hosts ${DOMAIN} | awk '{print $1}' | head -n1 || echo 'unknown')

5. Watch the first TLS cert issue:
     kubectl -n ${NAMESPACE} describe certificate
     kubectl -n ${NAMESPACE} get ingress
================================================================
EOF
