# MyApp

A fullstack SSR web application built with [Bifrost](https://github.com/3-lines-studio/bifrost) — a Go framework that server-renders React (TSX) pages with seamless hydration. Uses SQLite/Turso for persistence, Tailwind CSS for styling, and supports i18n and dark mode out of the box.

## Tech Stack

| Layer       | Technology                                           |
| ----------- | ---------------------------------------------------- |
| Backend     | Go 1.25, `net/http` stdlib router                    |
| SSR Engine  | Bifrost (Go + React SSR with hydration)              |
| Frontend    | React 19, TypeScript, Tailwind CSS 4                 |
| Database    | SQLite locally / Turso (libSQL) in production, GORM  |
| Migrations  | Atlas + atlas-provider-gorm (schema from GORM models)|
| Auth        | Stateless JWT in `session` cookie (bcrypt + HS256)   |
| i18n        | JSON locale files (English, Spanish)                 |
| Dev tooling | Air (hot reload), Bun (JS dependencies)              |

## Project Structure

```
.
├── main.go              # Entry point: wires DI graph, registers routes
├── config/
│   └── env.go           # Environment config (DB_DSN, JWT_SECRET)
├── model/
│   └── user.go          # User GORM model + UserRepository (CRUD)
├── services/
│   └── auth.go          # AuthService: signup, login, session resolution
├── handlers/
│   └── auth.go          # AuthHandler: HTTP layer, form parsing, redirects
├── util/
│   ├── db.go            # Database connection + Entity base struct (UUID PK, soft delete)
│   ├── jwt.go           # JWT sign/parse helpers
│   ├── error.go         # AppError type
│   └── uuid.go          # UUID generation helper
├── testutil/
│   └── db.go            # Test helper: in-memory SQLite DB with AutoMigrate
├── i18n/
│   ├── i18n.go          # Locale detection, translation loader, T() helper
│   └── locales/
│       ├── en.json      # English translations
│       └── es.json      # Spanish translations
├── pages/               # React TSX pages (server-rendered by Bifrost)
│   ├── app.css          # Global styles, Tailwind theme, dark mode tokens
│   ├── layout.tsx       # Shared layout (navbar, footer)
│   ├── home.tsx         # Home page
│   ├── login.tsx        # Login page
│   ├── signup.tsx       # Signup page
│   ├── theme-toggle.tsx # Dark/light mode toggle (client-side hydrated)
│   ├── theme-script.tsx # Inline script to prevent theme flash (FOUC)
│   └── i18n.ts          # Client-side t() helper with {{param}} interpolation
├── migrations/          # Atlas-generated SQL migration files
├── atlas.hcl            # Atlas config (reads schema from GORM models)
├── .air.toml            # Air hot-reload config (app :8080, proxy :3000)
├── Makefile             # Dev, build, and migration commands
├── go.mod
├── package.json
└── tsconfig.json
```

## Architecture

The app follows a **dependency injection** pattern with three distinct layers wired in `main()`:

```
handlers.AuthHandler
    └── services.AuthService
            └── model.UserRepository
                    └── *gorm.DB
```

Each layer is constructed explicitly and passed down — no globals, no service locator.

**Repository** (`model/`) — thin GORM wrappers that speak to the database. All queries are context-aware and respect soft deletes (`deleted_at IS NULL`).

**Service** (`services/`) — business logic. `AuthService` hashes passwords with bcrypt, signs JWT tokens, and resolves the current user from a request's session cookie.

**Handler** (`handlers/`) — HTTP boundary. Parses form values, calls the service, and redirects with localised error messages on failure.

## Getting Started

### Prerequisites

- **Go** 1.25+
- **Bun** (for JS dependencies)
- **Air** (optional, for hot reload): `go install github.com/air-verse/air@latest`
- **Atlas** (for migrations): see [Installation](#database-migrations)

### Setup

```bash
# Install JS dependencies
bun install

# Copy and fill in environment variables
cp .env.example .env.local

# Apply migrations to local DB
make migrations-apply-local

# Run in development mode (hot reload via Air, proxy on :3000)
make dev

# Or build and run directly
make build
make start
```

The dev server runs on `http://localhost:3000` (Air proxy) with the app on port `8080`.

## Features

### Server-Side Rendering with Bifrost

Pages are React TSX components in `pages/` rendered on the server by Bifrost. Each page uses `bifrost.WithPropsLoader` to inject server-side data (user session, locale, translations) as React props before render. The `.bifrost/` build artifacts are embedded into the Go binary at compile time via `//go:embed`.

Each page exports:
- A **default component** — the page content
- A **`Head` component** — for `<title>`, `<meta>`, and the theme-flash prevention script

### Authentication

Stateless cookie-based auth using bcrypt + JWT (HS256):

- `POST /api/signup` — validate form, hash password, create user, issue JWT
- `POST /api/login` — verify credentials, issue JWT
- `POST /api/logout` — clear the session cookie

The JWT is stored in an `HttpOnly`, `SameSite=Lax` cookie named `session`. On each request, `AuthService.GetUserFromRequest` parses the cookie, validates the token, and fetches the user from the database. There is no server-side session table.

### Internationalization (i18n)

Supports English (`en`) and Spanish (`es`). Locale detection order:

1. `lang` cookie (set via `POST /api/set-lang`)
2. `Accept-Language` request header
3. Falls back to `en`

Translations are loaded at startup from embedded JSON files (`i18n/locales/*.json`) and passed as props to every page. The client-side `t()` helper supports `{{param}}` interpolation.

### Dark Mode

Theme preference is persisted in `localStorage` and applied via a `.dark` CSS class. An inline script injected via `ThemeScript` reads the preference before first paint to prevent flash of unstyled content (FOUC). `ThemeToggle` is a hydrated client-side component.

## Database Migrations

Migrations are managed by **[Atlas](https://atlasgo.io)** using the `atlas-provider-gorm` plugin, which introspects GORM model structs to derive the target schema automatically.

### Install Atlas

```bash
# macOS (Homebrew)
brew install ariga/tap/atlas

# Linux / manual
curl -sSf https://atlasgo.sh | sh
```

### How it works

`atlas.hcl` configures Atlas to run `atlas-provider-gorm` against the `./model` package to produce a schema, then diff it against the current migration directory:

```hcl
data "external_schema" "gorm" {
  program = ["go", "run", "-mod=mod", "ariga.io/atlas-provider-gorm",
             "load", "--path", "./model", "--dialect", "sqlite"]
}
```

### Workflow

```bash
# 1. Generate a new migration after changing a GORM model
make migrations-generate name=add_bio_to_users

# 2. Apply to local SQLite DB
make migrations-apply-local

# 3. Apply to production (Turso)
make migrations-apply-prod
```

To add a new table, create a GORM model struct in `model/`, then run `make migrations-generate name=<description>`. Atlas will diff the new schema against the existing migrations and write a new SQL file to `migrations/`.

### Reset local DB

```bash
make dev-reset   # deletes dev.db and re-applies all migrations from scratch
```

## Testing

Tests use the standard `testing` package — no third-party test framework.

Each package has its own `_test.go` file covering its layer in isolation:

| File | What it tests |
|---|---|
| `model/user_test.go` | Repository CRUD: Create, GetByID, GetByEmail, Update, Delete, ExistsByEmail |
| `services/auth_test.go` | Service logic: Signup, Login (wrong password / user not found), GetUserFromRequest |
| `handlers/auth_test.go` | HTTP flows: form validation, redirect targets, session cookie set/cleared |

### Test database

`testutil.NewTestDB` opens an **in-memory SQLite** database and runs `AutoMigrate` on the provided models:

```go
func NewTestDB(t *testing.T, models ...any) *gorm.DB
```

Every test (or sub-test) that calls `newTestHandler(t)` / `newTestService(t)` gets a fresh, isolated database. `t.Cleanup` closes the connection after each test.

### Run tests

```bash
go test ./...
```

## Routes

| Method | Path            | Description                   |
| ------ | --------------- | ----------------------------- |
| GET    | `/`             | Home page (SSR)               |
| GET    | `/login`        | Login page (SSR)              |
| GET    | `/signup`       | Signup page (SSR)             |
| POST   | `/api/signup`   | Create account                |
| POST   | `/api/login`    | Authenticate                  |
| POST   | `/api/logout`   | Destroy session               |
| POST   | `/api/set-lang` | Switch language (en / es)     |

## Environment Variables

| Variable          | Default                  | Description                        |
| ----------------- | ------------------------ | ---------------------------------- |
| `DB_DSN`          | `file:dev.db`            | GORM data source name              |
| `JWT_SECRET`      | `dev-secret-change-me`   | HMAC secret for JWT signing        |
| `TURSO_DB_URL`    | —                        | Turso host (used by `migrations-apply-prod`) |
| `TURSO_AUTH_TOKEN`| —                        | Turso auth token                   |

Copy `.env.example` to `.env.local` and fill in the values. The `Makefile` loads `.env.$(ENV)` automatically (`ENV` defaults to `local`).

## Makefile Commands

| Command                          | Description                                              |
| -------------------------------- | -------------------------------------------------------- |
| `make dev`                       | Start dev server with hot reload (Air proxy on :3000)   |
| `make build`                     | Build Bifrost assets + Go binary                        |
| `make start`                     | Build and run the production binary                     |
| `make doctor`                    | Run Bifrost environment diagnostics                     |
| `make migrations-generate name=` | Diff GORM models → new SQL file in `migrations/`        |
| `make migrations-apply-local`    | Apply pending migrations to local SQLite DB             |
| `make migrations-apply-prod`     | Apply pending migrations to Turso (production)          |
| `make inspect-models`            | Inspect the Atlas schema derived from GORM models       |
| `make dev-reset`                 | Delete local DB and re-apply all migrations             |
