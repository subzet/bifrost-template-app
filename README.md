# MyApp

A fullstack SSR web application built with [Bifrost](https://github.com/3-lines-studio/bifrost) — a Go framework that server-renders React (TSX) pages with seamless hydration. Uses SQLite for persistence, Tailwind CSS for styling, and supports i18n and dark mode out of the box.

## Tech Stack

| Layer       | Technology                                      |
| ----------- | ----------------------------------------------- |
| Backend     | Go 1.25, `net/http` stdlib router               |
| SSR Engine  | Bifrost (Go + React SSR with hydration)         |
| Frontend    | React 19, TypeScript, Tailwind CSS 4             |
| Database    | SQLite (via `modernc.org/sqlite`, pure Go)      |
| Migrations  | Goose (SQL files auto-generated from Go models) |
| Auth        | Session-based (bcrypt + secure cookies)         |
| i18n        | JSON locale files (English, Spanish)            |
| Dev tooling | Air (hot reload), Bun (JS dependencies)         |

## Project Structure

```
.
├── main.go              # HTTP server, routes, and API handlers
├── auth/
│   └── auth.go          # HTTP-level session helpers (cookies, request user)
├── db/
│   ├── db.go            # SQLite connection (WAL mode, foreign keys), goose migration
│   └── migrations/      # Auto-generated goose SQL files (gitignored)
├── cmd/
│   └── dbgen/
│       └── main.go      # Migration file generator (reads model schemas)
├── model/
│   ├── schema.go        # Migration registry (ordered list of migrations)
│   ├── user.go          # User model, CRUD, migration definition
│   └── session.go       # Session model, CRUD, migration definition
├── i18n/
│   ├── i18n.go          # Locale detection, translation loader
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
│   └── i18n.ts          # Client-side translation helper with interpolation
├── .bifrost/            # Bifrost build artifacts (auto-generated)
├── .air.toml            # Air hot-reload config
├── Makefile             # Dev, build, and migration commands
├── go.mod
├── package.json
└── tsconfig.json
```

## Getting Started

### Prerequisites

- **Go** 1.25+
- **Bun** (for installing JS dependencies)
- **Air** (optional, for hot reload during development)

### Setup

```bash
# Install JS dependencies
bun install

# Run in development mode (hot reload via Air)
make dev

# Or build and run for production
make build
make start
```

The dev server runs on `http://localhost:3000` (Air proxy) with the app on port `8080`.

## Features

### Server-Side Rendering with Bifrost

Pages are defined as React TSX components in `pages/` and rendered on the server by Bifrost. Each page uses `bifrost.WithPropsLoader` to inject server-side data (user session, locale, translations) as React props. The `.bifrost/` directory is embedded into the Go binary at compile time.

Each page exports:
- A **default component** — the page content
- A **`Head` component** — for `<title>`, `<meta>`, and the theme-flash prevention script

### Authentication

Cookie-based session auth with bcrypt password hashing:

- **POST /api/signup** — Create account (email + password, min 8 chars, confirmation)
- **POST /api/login** — Authenticate and set session cookie
- **POST /api/logout** — Destroy session and clear cookie

Sessions are stored in SQLite with a 32-byte random hex token. Cookies are `HttpOnly` and `SameSite=Lax`.

### Internationalization (i18n)

Supports English (`en`) and Spanish (`es`). Locale is detected from:

1. `lang` cookie (set via **POST /api/set-lang**)
2. `Accept-Language` header
3. Falls back to `en`

Translations are loaded from `i18n/locales/*.json` and passed as props to every page. The client-side `t()` helper supports `{{param}}` interpolation.

### Dark Mode

Theme preference is persisted in `localStorage` and applied via a `.dark` CSS class. An inline script in `<head>` (via `ThemeScript`) reads the preference before paint to prevent flash of unstyled content. The `ThemeToggle` component is a hydrated client-side widget.

### Database & Migrations

SQLite with WAL journal mode and foreign keys enabled. Table schemas are defined as `Migration` structs alongside their models in the `model/` package. The migration generator (`cmd/dbgen`) reads the ordered list from `model.Migrations()` and writes numbered goose SQL files to `db/migrations/`. On startup, `db.Migrate()` runs `goose.Up` to apply pending migrations.

To add a new table:
1. Define the schema as a `Migration` var in the model file
2. Add it to the slice in `model/schema.go` → `Migrations()`
3. Run `make migrate-generate` to regenerate SQL files

### Styling

Tailwind CSS 4 with a custom shadcn/ui-style design token system using oklch colors for both light and dark themes, defined in `pages/app.css`.

## Routes

| Method | Path           | Description                  |
| ------ | -------------- | ---------------------------- |
| GET    | `/`            | Home page                    |
| GET    | `/login`       | Login page                   |
| GET    | `/signup`      | Signup page                  |
| POST   | `/api/signup`  | Create account               |
| POST   | `/api/login`   | Authenticate                 |
| POST   | `/api/logout`  | Destroy session              |
| POST   | `/api/set-lang`| Switch language (en/es)      |

## Makefile Commands

| Command                | Description                                      |
| ---------------------- | ------------------------------------------------ |
| `make dev`             | Generate migrations, start dev server with hot reload |
| `make build`           | Generate migrations, build Bifrost assets + Go binary |
| `make start`           | Build and run production binary                  |
| `make doctor`          | Run Bifrost diagnostics                          |
| `make migrate-generate`| Generate goose SQL files from model schemas      |
| `make migrate`         | Generate + apply migrations with goose           |
| `make migrate-reset`   | Delete DB and re-run all migrations              |
