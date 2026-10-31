# Chirpy 🐦

A Twitter-style microblogging **REST API** written in Go, backed by PostgreSQL. Users sign up, log in, and post short messages ("chirps") of up to 140 characters.

## Why it's worth a look

Chirpy is a compact, dependency-light example of how to build a real-world backend in Go using mostly the standard library:

- **Authentication done properly**: JWT access tokens plus revocable refresh tokens
- **Type-safe database access**: SQL queries compiled to Go with [sqlc](https://sqlc.dev/)
- **Clean REST design**: Go 1.22+ `net/http` method-and-path routing, no web framework
- **Real-world features**: content moderation, ownership checks, webhooks, and server metrics

## Features

- User registration and profile updates
- Login with JWT access tokens, refresh tokens, and token revocation
- Create, read, list, and delete chirps (only the author can delete)
- Filter chirps by author and sort by creation date
- Profanity filter: banned words are replaced with `****`
- Webhook endpoint for third-party payment upgrades (Polka), protected by an API key
- Health check, request metrics, and an admin reset endpoint
- Static file server with hit counting

## Tech Stack

| Layer      | Technology                              |
|------------|-----------------------------------------|
| Language   | Go                                      |
| Database   | PostgreSQL (`lib/pq` driver)            |
| Query gen  | sqlc                                    |
| Auth       | JWT + refresh tokens                    |
| Config     | `.env` via `godotenv`                   |

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.22 or newer
- [PostgreSQL](https://www.postgresql.org/download/)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) (only if you change the SQL queries)
- A migration tool such as [goose](https://github.com/pressly/goose) to apply the schema in `sql/`

### 1. Clone the repo

```bash
git clone https://github.com/junaidmd0980/chirpy.git
cd chirpy
```

### 2. Create the database

```bash
psql postgres -c "CREATE DATABASE chirpy;"
```

### 3. Configure environment variables

Create a `.env` file in the project root:

```env
DB_URL="postgres://USER:PASSWORD@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
JWT_SECRET="<a long random string>"
POLKA_KEY="<your webhook API key>"
```

Generate a strong JWT secret with:

```bash
openssl rand -base64 64
```

| Variable     | Description                                                      |
|--------------|------------------------------------------------------------------|
| `DB_URL`     | PostgreSQL connection string (**required**)                      |
| `PLATFORM`   | Environment name; `dev` enables the admin reset endpoint         |
| `JWT_SECRET` | Secret used to sign and verify JWTs                              |
| `POLKA_KEY`  | API key that authorizes incoming Polka webhook requests          |

### 4. Run database migrations

```bash
goose -dir sql/schema postgres "$DB_URL" up
```

### 5. Install dependencies and run

```bash
go mod download
go run .
```

The server starts on **http://localhost:8080**.

Quick check:

```bash
curl http://localhost:8080/api/healthz
```

## API Reference

### Users & Auth

| Method | Endpoint        | Description                          |
|--------|-----------------|--------------------------------------|
| POST   | `/api/users`    | Register a new user                  |
| PUT    | `/api/users`    | Update the authenticated user        |
| POST   | `/api/login`    | Log in; returns access + refresh tokens |
| POST   | `/api/refresh`  | Exchange a refresh token for a new access token |
| POST   | `/api/revoke`   | Revoke a refresh token               |

### Chirps

| Method | Endpoint                | Description                              |
|--------|-------------------------|------------------------------------------|
| POST   | `/api/chirps`           | Create a chirp (auth required, max 140 chars) |
| GET    | `/api/chirps`           | List chirps                              |
| GET    | `/api/chirps/{chirpID}` | Get a single chirp                       |
| DELETE | `/api/chirps/{chirpID}` | Delete your own chirp (auth required)    |

**Query parameters for `GET /api/chirps`**

- `author_id=<uuid>`: only return chirps by this user
- `sort=asc|desc`: sort by creation time (default `asc`)

### Webhooks

| Method | Endpoint               | Description                                   |
|--------|------------------------|-----------------------------------------------|
| POST   | `/api/polka/webhooks`  | Receives payment events to upgrade users (requires API key) |

### System

| Method | Endpoint         | Description                                    |
|--------|------------------|------------------------------------------------|
| GET    | `/api/healthz`   | Readiness check                                |
| GET    | `/admin/metrics` | Number of hits to the static file server       |
| POST   | `/admin/reset`   | Reset metrics and delete all users (dev only)  |
| GET    | `/app/`          | Static file server                             |

Protected endpoints expect a bearer token:

```
Authorization: Bearer <access_token>
```

## Example Usage

```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "s3cret"}'

# Log in
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "s3cret"}'

# Post a chirp
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"body": "Hello, Chirpy!"}'

# List newest chirps first
curl "http://localhost:8080/api/chirps?sort=desc"
```

## Project Structure

```
chirpy/
├── main.go               # Server setup, config, and routes
├── handler_users.go      # User registration, updates, Polka upgrades
├── handler_login.go      # Login
├── handlerRefresh.go     # Refresh-token exchange
├── handlerRevoke.go      # Refresh-token revocation
├── handler_chirps.go     # Chirp CRUD, validation, filtering
├── json_helper.go        # JSON response helpers
├── metrics.go            # Hit counter and metrics page
├── readiness.go          # Health check
├── reset.go              # Admin reset
├── internal/             # Auth helpers and sqlc-generated database code
├── sql/                  # Schema migrations and SQL queries
├── sqlc.yaml             # sqlc configuration
└── assets/, index.html   # Static content
```

## Development

Regenerate the Go database code after editing SQL queries:

```bash
sqlc generate
```
