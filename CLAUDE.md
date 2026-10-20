# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Go chat server (Gin HTTP + Gorilla WebSockets) with Postgres, two Redis instances, and MinIO. Designed to run horizontally scaled behind an ingress — the three app pods coordinate real-time delivery via Redis Pub/Sub, so a message produced on one pod reaches sockets served by another. Current work is on the `feature/group_chats` branch, extending the single-channel-per-user Pub/Sub model to also subscribe each socket to its group chat topics.

## Build and run

Local (docker-compose):
- `docker compose up --build` — builds the app image and starts db / cache / redis_pubsub / minio-server / application on `:8080`.
- `go build ./...` / `go run ./cmd` — straight Go build; must be run with the config/dependencies reachable (see "Config" below).
- `go mod tidy` — module graph is in `go.mod` (Go 1.21).

Kubernetes (minikube) — this is the deployment target used in `README.md`:
- `make create-app-docker-image` — builds `myapp:latest` inside minikube's docker daemon.
- `make setup-minikube` → `make deploy-all` → `minikube tunnel` — one-time bring-up (namespace `enrollment`).
- `make redeploy` / `make deploy-update` — roll a new image / re-apply manifests.
- `make logs` — tail the application deployment logs.
- `make kuber-gui` — opens the minikube dashboard.

Reset data (compose uses persistent volumes so flushes are required between runs):
- `make flush-all` — truncate Postgres tables + FLUSHALL both Redis instances in docker-compose.
- `make flush-all-kuber` — same, against the k8s pods.

Testing: there are **no `*_test.go` files in the repo**. If you add tests, there is no established pattern yet — don't invent conventions without checking with the user.

## Config

Two sources, both loaded in `cmd/main.go`:
- `cmd/internal/config.yml` via viper — hostnames (`db`, `cache`, `pubsub`, `minio-server`) default to docker-compose service names; the commented `localhost` lines are for running `go run ./cmd` against port-forwarded containers.
- `.env` via `gotenv` — only `DB_PASSWORD` is read from env; the Redis/MinIO credentials are inline in `config.yml` and the compose file.

Migrations run automatically on startup via `postgres.RunMigrations("file://migrations", db)`, so new SQL files in `migrations/` apply the next time the app boots.

## Architecture

The code is layered strictly: `handlers → service → repository`. Wiring happens once in `cmd/main.go`, which constructs each layer top-down and injects it into the next.

**Repository layer** (`pkg/repository/`) — one subpackage per backend:
- `postgres/` — pgx pool. Each entity (`Chat`, `Message`, `User`, `Info`, `Device`) is an interface declared in `repository.go`; implementations live under `postgres/internal/`. The `Repository` struct embeds all of them so callers see one flat surface.
- `redis/` — split into **two separate clients** by role: `Cache` (keyed data, currently chat members) and `PubSub` (`MessageBus` interface over `redis.Publish`/`Subscribe`). They are configured independently (`cache.*` vs `pubsub.*` in `config.yml`) and in k8s run as two separate deployments.
- `minio/` — S3Storage wrapper for attachment uploads/downloads.

**Service layer** (`pkg/service/service.go`) — mirrors the repository pattern: one interface per domain, a composite `Service` struct that embeds all of them. Message send fans out: `MessageService.AddMessage` writes to Postgres, then publishes the serialized message to **one Pub/Sub channel per recipient** (`user:<id>`).

**Handlers** are split by transport:
- `pkg/handlers/http/` — Gin routes for account/chat/file endpoints, with `gin-rate-limit` (80 req/s per IP) applied per group.
- `pkg/handlers/ws/` — single `GET /create` endpoint that upgrades to WebSocket.

**WebSocket flow** — the lifecycle in `pkg/handlers/ws/websocket.go` + `pkg/connection/websocket/`:
1. Upgrade; register the `Client` with the `ServerHub` (map keyed by `userId`, so one connection per user).
2. Subscribe to `user:<id>` (personal) **and** to each of the user's group chat topics (group feature in progress — see `GetGroupsListIDs`).
3. Flush missed messages from Postgres (`GetMissedMessages` uses `user_last_seen` + `last_message` to compute the gap).
4. `go ConsumeMessagesFromPubSub(...)` — forwards Pub/Sub payloads to the socket.
5. `client.ReadMessage(hub)` runs in the request goroutine — reads incoming JSON, looks up chat members, hands off to `MessageService.AddMessage` which persists + publishes to every member's `user:<id>` channel.
6. On disconnect, `hub.Unregister` saves `lastSeen` for that `(userId, deviceId)` in the `devices` table.

The horizontal-scaling invariant: because every pod subscribes recipients to Pub/Sub, a producer pod never needs to know which pod the recipient's socket is on. The Pub/Sub Redis is shared; the cache Redis is also shared but is per-entity key/value only.

## Current in-flight work (feature/group_chats)

The dirty files on this branch add a `GroupChatIds` column (migration `01_init_db.up.sql`), a `GetGroupsListIDs` lookup, and group-topic subscriptions in the WS handler. Note that **`pkg/service/info.go` currently has a syntax error** (`members, err = s.repo.(ctx, chatId)` with a dangling `}`) — this is uncompiled WIP, not a repo-wide pattern. Don't model other code on it.

## Conventions worth matching

- Error handling is almost entirely `log.Printf` / `logrus.Infof` + bare `return`. The repository layer returns errors; handlers mostly log and drop them. Match the surrounding style — don't introduce a new error-wrapping scheme without asking.
- Repository implementations live in a non-exported `internal/` subpackage and are constructed via `NewXxxPostgres` / `NewXxxCache` factories. Public types are the interfaces in the parent package.
- SQL table names come from constants in `pkg/repository/postgres/internal/utils.go` — reference them via those constants rather than hardcoding strings.
- The composite `Service` and `Repository` structs embed sub-interfaces to flatten the API. When adding a new domain, add an interface in `service.go` / `repository.go`, embed it, and wire the implementation in the matching `NewService` / `NewRepository`.
