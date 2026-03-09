# MYN CLI — Development Guide

## Project
- **Repo**: github.com/mindyournow/myn-cli
- **Binary**: `mynow`
- **Language**: Go 1.24+
- **TUI**: Bubble Tea (charmbracelet)
- **CLI**: Cobra
- **License**: MIT (open source)

## Build & Test
```bash
make build              # Build to bin/mynow
make test               # Unit tests
make test-integration   # Integration tests (needs Docker + MYN backend source)
make lint               # golangci-lint
```

## Architecture
- `cmd/mynow/` — CLI entry point, command definitions
- `internal/app/` — Application logic shared by CLI and TUI
- `internal/api/` — HTTP client for MYN backend
- `internal/auth/` — OAuth PKCE + credential storage
- `internal/config/` — Configuration loading (XDG)
- `internal/output/` — Text/JSON output formatting
- `internal/tui/` — Bubble Tea TUI screens
- `plugins/` — Plugin interface (proprietary plugins stay external)
- `test/integration/` — Integration tests with Docker Compose

## Key Rules
- **No proprietary code** in this repo — it's open source
- **Backend is a black box** — interact only via HTTPS API
- **Plugins must be optional** — core client works without any plugins
- All output must support `--json` for scripting
- Prefer `internal/` packages to keep API surface minimal
- Use `MYN_API_URL` env var to override backend URL

## Integration Testing
- Tests spin up MYN backend via Docker Compose (PostgreSQL + Redis + Spring Boot)
- Set `MYN_INTEGRATION_TEST=1` to enable
- Backend source path: `MYN_BACKEND_PATH` (defaults to `~/Projects/myn/api`)
- Or point `MYN_TEST_BACKEND_URL` at an already-running instance
- Demo account flow (`/api/v1/admin/demo/recreate-account`) bootstraps test data

## Issue Tracking
- GitHub Issues on this repo (not Linear)
