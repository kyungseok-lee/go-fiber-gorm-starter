# Repository Guidelines

## Project Structure & Module Organization

This is a Go 1.24 REST API starter using Fiber, GORM, MySQL/PostgreSQL, and Prometheus middleware. The executable entry point is `cmd/server/main.go`. Application internals live under `internal/`: `config`, `db`, `http`, `middleware`, `metrics`, `logger`, and domain packages such as `internal/domain/user`. Shared packages belong in `pkg/`, currently `pkg/resp`. SQL migrations are in `migrations/`, scripts are in `scripts/`, and CI is in `.github/workflows/ci.yml`.

## Build, Test, and Development Commands

Use the Makefile as the primary workflow:

- `make run` starts the API from `cmd/server/main.go`.
- `make dev` starts with `air` auto-reload when installed, otherwise falls back to `make run`.
- `make build` builds the local binary `fiber-gorm-starter`.
- `make test` runs `go test -v -race ./...`.
- `make coverage` writes `coverage.out` and `coverage.html`.
- `make lint` runs `golangci-lint run --timeout=5m`; CI and `make install-tools` use v1.64.2.
- `make check` runs format, vet, lint, tests, and build.
- `make docker-up` or `make docker-up-pg` starts MySQL or PostgreSQL stacks.
- `make migrate-create name=add_example`, `make migrate-up`, and `make migrate-down` manage migrations.

Copy `.env.example` to `.env` before local database-backed runs.

## Coding Style & Naming Conventions

Follow `.editorconfig`: Go files use tabs and gofmt formatting; YAML, SQL, shell, JSON, and Markdown use two-space indentation where applicable. Run `make fmt` before committing. Keep package names short, lowercase, and domain-oriented. Place business logic in `internal/domain/<name>` and avoid exports unless used across package boundaries.

Imports are organized by `goimports`/`gci` with local imports under `github.com/kyungseok-lee/go-fiber-gorm-starter`. Respect `.golangci.yml`, including the 120-character line limit and security/static analysis checks.

## Testing Guidelines

Tests use Go’s testing package with `stretchr/testify`. Keep tests next to implementation files using the `*_test.go` suffix, as in `internal/domain/user/service_test.go`. Prefer table-driven tests for service and repository behavior. Run `make test` normally and `make coverage` when changing shared logic.

## Commit & Pull Request Guidelines

Git history uses Conventional Commit-style messages such as `feat: ...`, `fix: ...`, `refactor: ...`, and scoped forms like `feat(init): ...`. Keep commits focused and imperative.

Pull requests should include a summary, linked issue when available, migration or environment changes, and test evidence such as `make check`. Include API examples or screenshots only for externally visible behavior or docs changes.

## Security & Configuration Tips

Do not commit `.env`, credentials, generated binaries, coverage artifacts, `.omc/`, or `.serena/`. Use `.env.example` for new configuration keys. When changing database behavior, update migrations and verify both MySQL and PostgreSQL paths when practical.
