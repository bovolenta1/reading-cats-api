# Reading Cats API — Setup & Development Guide (Windows + AWS SAM + Neon Postgres)

This repo is the backend API for **Reading Cats**, built with **Go**, deployed on **AWS Lambda + API Gateway (HTTP API)** using **AWS SAM**, and backed by **Postgres (Neon)**.
It uses **golang-migrate** for database migrations.

---

## Table of Contents

- [Reading Cats API — Setup \& Development Guide (Windows + AWS SAM + Neon Postgres)](#reading-cats-api--setup--development-guide-windows--aws-sam--neon-postgres)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
    - [Required](#required)
    - [Recommended](#recommended)
  - [Project layout (high level)](#project-layout-high-level)
  - [Environment variables](#environment-variables)
    - [`env.local` (local only)](#envlocal-local-only)
  - [Local development (SAM)](#local-development-sam)
    - [Build](#build)
    - [Start the API locally (HTTP)](#start-the-api-locally-http)
    - [Clean SAM artifacts](#clean-sam-artifacts)
  - [Database migrations](#database-migrations)
    - [Install migrate CLI](#install-migrate-cli)
    - [Create a new migration](#create-a-new-migration)
    - [Run migrations locally](#run-migrations-locally)
    - [Migration files example (create\_user)](#migration-files-example-create_user)
  - [Production: how to run migrations](#production-how-to-run-migrations)
    - [Recommended: GitHub Actions step](#recommended-github-actions-step)
    - [Concurrency / safety](#concurrency--safety)
  - [Troubleshooting](#troubleshooting)
    - [`make` commands fail with `'test' is not recognized`](#make-commands-fail-with-test-is-not-recognized)
    - [`migrate` not found](#migrate-not-found)
    - [Driver errors with `postgresql://`](#driver-errors-with-postgresql)
    - [Neon SSL errors](#neon-ssl-errors)
  - [Makefile reference (Windows / cmd-friendly)](#makefile-reference-windows--cmd-friendly)

---

## Prerequisites

### Required
- **Go** installed and available in PATH:
  - `go version`
- **AWS SAM CLI** installed:
  - `sam --version`
- **Docker Desktop** installed (required by `sam local`)
- **Postgres database** (Neon) and a connection string

### Recommended
- **GNU Make** installed (Windows):
  - Installed via Chocolatey: `choco install make -y`
  - Check: `make --version`

---

## Project layout (high level)

Typical structure (may evolve as features grow):

- `main.go` — Lambda entry point
- `internal/` — application code (domain/app/infra/presentation)
- `migrations/` — SQL migrations (golang-migrate format)
- `template.yaml` — AWS SAM template
- `Makefile` — local dev + migration commands

---

## Environment variables

### `env.local` (local only)

Create `env.local` at the repo root (do **not** commit it):

```
DATABASE_URL=postgres://USER:PASSWORD@HOST/DBNAME?sslmode=require&channel_binding=require
```

Notes:
- Some tools prefer `postgres://` over `postgresql://`. If you have issues, use `postgres://`.
- Keep credentials out of Git. Add `env.local` to `.gitignore`.

Example `.gitignore` entry:

```
env.local
```

---

## Local development (SAM)

### Build
```
make build
```

### Start the API locally (HTTP)
```
make start
```

Default: `http://localhost:3001`

### Clean SAM artifacts
```
make clean
```

---

## Database migrations

This project uses **golang-migrate/migrate**:
- Each migration has two files:
  - `NNNNNN_name.up.sql`
  - `NNNNNN_name.down.sql`

### Install migrate CLI

Install the CLI with Postgres support:

```
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Ensure your Go bin is in PATH (usually `C:\Users\<you>\go\bin`).
Check:

```
migrate -version
```

---

### Create a new migration

Using Make:

```
make migrate-create name=create_groups
```

This generates:
- `migrations/000002_create_groups.up.sql`
- `migrations/000002_create_groups.down.sql`

---

### Run migrations locally

Run all pending migrations:

```
make migrate-up
```

Show current version:

```
make migrate-version
```

Rollback a single migration:

```
make migrate-down
```

---

### Migration files example (create_user)

**`migrations/000001_create_user.up.sql`**

```sql
CREATE TABLE IF NOT EXISTS users (
  id             UUID PRIMARY KEY,
  cognito_sub    TEXT NOT NULL UNIQUE,
  email          TEXT,
  display_name   TEXT,
  avatar_url     TEXT,
  profile_source TEXT NOT NULL DEFAULT 'idp', -- 'idp' | 'user'
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Email optional, but unique when present (case-insensitive)
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique
ON users (lower(email))
WHERE email IS NOT NULL;
```

**`migrations/000001_create_user.down.sql`**

```sql
DROP INDEX IF EXISTS users_email_unique;
DROP TABLE IF EXISTS users;
```

---

## Production: how to run migrations

### Recommended: GitHub Actions step

Run migrations as part of your CI/CD pipeline **before** `sam deploy`.

Example (Ubuntu runner):

```yaml
- name: Install migrate
  run: |
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    echo "$(go env GOPATH)/bin" >> $GITHUB_PATH

- name: Run migrations
  env:
    DATABASE_URL: ${{ secrets.DATABASE_URL_PROD }}
  run: |
    migrate -path migrations -database "$DATABASE_URL" up

- name: Build
  run: sam build

- name: Deploy
  run: sam deploy --no-confirm-changeset --no-fail-on-empty-changeset
```

### Concurrency / safety

- Ensure only **one deployment** runs at a time (to avoid two pipelines migrating concurrently).
- GitHub Actions tip: use workflow `concurrency` to serialize deploys.
- Migrations should be **forward-only** for production (down is mainly for local/dev).
- Avoid editing old migrations that have already run in production.

---

## Troubleshooting

### `make` commands fail with `'test' is not recognized`
Your Make recipes are being executed by `cmd.exe` on Windows. Use **cmd-friendly recipes** (no `test`, no bash syntax), or run `make` inside Git Bash/WSL.

### `migrate` not found
Your Go bin folder is not in PATH. Add:
- `C:\Users\<you>\go\bin`

Check:
```
where.exe migrate
migrate -version
```

### Driver errors with `postgresql://`
Some tooling expects `postgres://`. Try switching the scheme in `env.local`:
- from `postgresql://...` to `postgres://...`

### Neon SSL errors
Ensure your URL includes:
- `sslmode=require`
- (optional) `channel_binding=require` (ok to keep)

---

## Makefile reference (Windows / cmd-friendly)

A minimal Makefile for this setup:

```makefile
.PHONY: build start clean migrate-create migrate-up migrate-down migrate-version

MIGRATIONS_DIR=migrations
MIGRATE=migrate

build:
  sam build

start: build
  sam local start-api -p 3001

clean:
  @if exist .aws-sam rmdir /s /q .aws-sam

define LOAD_DB_URL
  @if not defined DATABASE_URL ( \
    if exist env.local ( \
      for /f "usebackq tokens=1,* delims==" %%A in ("env.local") do ( \
        if "%%A"=="DATABASE_URL" set "DATABASE_URL=%%B" \
      ) \
    ) \
  )
  @if not defined DATABASE_URL (echo DATABASE_URL not set. Put it in env.local or set it in the environment. & exit /b 1)
endef

migrate-create:
  @if "$(name)"=="" (echo use: make migrate-create name=create_user & exit /b 1)
  $(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
  $(LOAD_DB_URL)
  $(MIGRATE) -path $(MIGRATIONS_DIR) -database "%DATABASE_URL%" up

migrate-down:
  $(LOAD_DB_URL)
  $(MIGRATE) -path $(MIGRATIONS_DIR) -database "%DATABASE_URL%" down 1

migrate-version:
  $(LOAD_DB_URL)
  $(MIGRATE) -path $(MIGRATIONS_DIR) -database "%DATABASE_URL%" version
```

---

If you want, add a dedicated section later for:
- `/v1/me` endpoint contract
- auth flow (Cognito JWT authorizer)
- local testing with sample JWT / authorizer simulation
