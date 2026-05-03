# Repository Guidelines

## Project Structure & Module Organization

This is a Go 1.24 REST API starter using Fiber, GORM, MySQL/PostgreSQL, and Prometheus. The executable entry point is `cmd/server/main.go`. Application internals live under `internal/`: `config`, `db`, `http`, `middleware`, `metrics`, `logger`, and domain packages such as `internal/domain/user`. Shared packages belong in `pkg/`, currently `pkg/resp`. SQL migrations are in `migrations/`, scripts are in `scripts/`, and CI is in `.github/workflows/ci.yml`.

## Build, Test, and Development Commands

Use the Makefile as the primary workflow:

- `make run` starts the API from `cmd/server/main.go`.
- `make dev` starts with `air` auto-reload when installed, otherwise falls back to `make run`.
- `make build` builds the local binary `fiber-gorm-starter`.
- `make test` runs `go test -v -race ./...`.
- `make e2e` runs black-box user API checks against a running server.
- `make lint` runs `golangci-lint run --timeout=5m`; CI and `make install-tools` use v1.64.2.
- `make check` runs format, vet, lint, tests, and build.
- `make docker-up` or `make docker-up-pg` starts MySQL or PostgreSQL stacks; the PostgreSQL target passes the compose DB host/driver overrides for the app container.
- `make migrate-create name=add_example`, `make migrate-up`, and `make migrate-down` manage migrations.

## Coding Style & Naming Conventions

Follow `.editorconfig`: Go files use tabs and gofmt; YAML, SQL, shell, JSON, and Markdown use two-space indentation where applicable. Run `make fmt` before committing. Keep package names short, lowercase, and domain-oriented. Place business logic in `internal/domain/<name>` and avoid exports unless used across package boundaries.

Imports are organized by `goimports`/`gci` with local imports under `github.com/kyungseok-lee/go-fiber-gorm-starter`. Respect `.golangci.yml`, including the 120-character line limit.

## Testing Guidelines

Tests use Go’s testing package with `stretchr/testify`. Keep tests next to implementation files using `*_test.go`, as in `internal/domain/user/service_test.go`. Prefer table-driven tests for service and repository behavior; repository tests use in-memory SQLite through GORM. Run `make test` normally and `make coverage` when changing shared logic. For runtime API changes, start Docker Compose and run `make e2e`.

## Commit & Pull Request Guidelines

Git history uses Conventional Commit-style messages such as `feat: ...`, `fix: ...`, `refactor: ...`, and scoped forms like `feat(init): ...`. Keep commits focused and imperative.

Pull requests should include a summary, linked issue when available, migration or environment changes, and test evidence such as `make check`. Include API examples only for externally visible behavior or docs changes.

## Security & Configuration Tips

Do not commit `.env`, credentials, generated binaries, coverage artifacts, `.omc/`, or `.serena/`. Use `.env.example` for new configuration keys, and express duration values with Go duration units such as `300s`. Production config must not rely on placeholder `API_KEY`, default `DB_PASS=password`, or wildcard CORS origins. When changing database behavior, update migrations and verify both MySQL and PostgreSQL paths when practical.
