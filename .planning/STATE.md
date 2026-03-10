# Agent State: CLI-1

## Current Status: IN PROGRESS — Implementation Barely Begun

The previous agent (Kimi K2.5) built a solid scaffold and fixed 16 reviewer blockers,
but only implemented OAuth PKCE. Everything else is stubs or missing.
A comprehensive Opus 4.6 audit estimated ~5% of spec coverage.

**Model switch**: Implementation agent changed from kimi-k2.5 to claude-sonnet-4-6.

---

## Spec & PRD

- **Spec**: `docs/MYN-TUI-CLI_SPEC.md` (4,500+ lines) — single source of truth
- **Issue**: CLI-1 on GitHub (`github.com/mindyournow/myn-cli/issues/1`)
- **Review**: `~/Projects/ReviewOfKimiMYNTUICode.md` — full Opus 4.6 code review

---

## What EXISTS (keep and build on)

| File | Status | Notes |
|---|---|---|
| `cmd/mynow/main.go` | Partial | Cobra root + 8 command groups wired. Missing: 25 command groups, --api-url, --debug flags |
| `internal/api/client.go` | Good | Generic HTTP client with retry, rate limiting, context. Needs: 401 auto-refresh, response body size limit, Retry-After cap |
| `internal/auth/oauth.go` | Good | Full OAuth PKCE flow, correct endpoints. Needs: server timeouts, 127.0.0.1 binding |
| `internal/auth/keyring.go` | Needs Fix | AES-GCM + PBKDF2. Needs: fix weak key fallback (HIGH-1) |
| `internal/auth/auth.go` | Good | TokenStore interface (SaveRefreshToken, LoadRefreshToken, Clear) |
| `internal/config/config.go` | Partial | URL validation good. Needs: YAML config file, additional env vars |
| `internal/output/output.go` | Partial | Text/JSON formatter. Needs: table, color, markdown (Glamour), progress |
| `test/integration/setup.go` | Good | Docker Compose lifecycle |
| `test/integration/docker-compose.yml` | Needs Fix | Hardcoded secrets (HIGH-2) |
| Tests (7 files) | Mixed | client_test, oauth_test, keyring_test, config_test, output_test are solid. app_test is worthless stubs. |

## What is MISSING (must build)

### Architectural Bugs (do first — these affect correctness)
- [ ] BUG-1: `app.go` methods return `Formatter.Error()` which returns nil — commands exit 0 on failure. Must return actual errors.
- [ ] BUG-2: `Formatter.Error()` writes to stdout, not stderr (Spec §13.2 says "Errors go to stderr")
- [ ] BUG-3: `app.New()` called before flag parsing in main.go — prevents `--api-url` flag from working. Must defer config loading to PersistentPreRunE.

### Security Fixes (do first)
- [ ] HIGH-1: Fix weak key derivation fallback in keyring.go
- [ ] HIGH-2: Replace hardcoded secrets in docker-compose.yml with env var refs
- [ ] MED-1: Add ReadTimeout/WriteTimeout to OAuth callback server
- [ ] MED-2: Bind callback to 127.0.0.1:0 instead of localhost:0
- [ ] MED-3: Sanitize errorParam in OAuth callback before printing
- [ ] MED-4: Cap Retry-After to 60 seconds in client.go
- [ ] MED-5: Add timeout to integration test HTTP client
- [ ] LOW-1: Add io.LimitReader for response bodies in client.go

### Dependencies to Add
- [ ] `go get github.com/charmbracelet/bubbletea` (TUI framework)
- [ ] `go get github.com/charmbracelet/lipgloss` (TUI styling)
- [ ] `go get github.com/charmbracelet/bubbles` (TUI components)
- [ ] `go get github.com/charmbracelet/glamour` (markdown rendering)
- [ ] `go get github.com/zalando/go-keyring` (OS keychain)
- [ ] `go get gopkg.in/yaml.v3` (YAML config)
- [ ] `go mod tidy` (fix x/crypto indirect annotation)

### Auth (Spec §2)
- [ ] `internal/auth/apikey.go` — API key storage + validation via GET /api/v1/customers
- [ ] `internal/auth/device.go` — Device authorization flow (stub; backend doesn't support yet)
- [ ] `internal/auth/tokens.go` — Token refresh, in-memory access token cache, 401 auto-refresh
- [ ] GNOME Keyring / KDE Wallet integration via go-keyring
- [ ] `login --api-key` command
- [ ] `login --device` command
- [ ] `whoami` command
- [ ] `auth status` / `auth refresh` commands
- [ ] `logout` — revoke via POST /api/mcp/oauth/logout

### Config (Spec §3)
- [ ] YAML config file support (`~/.config/mynow/config.yaml`)
- [ ] All env vars: MYN_API_KEY, MYNOW_CONFIG, MYNOW_KEYRING, NO_COLOR, MYNOW_DEBUG
- [ ] `config show/set/get/reset/path` commands

### API Domain Layer (Spec §1.2, Appendix A — 22 files)
- [ ] `internal/api/tasks.go` — /api/v2/unified-tasks (CRUD, complete, archive, batch, move)
- [ ] `internal/api/habits.go` — /api/habits/chains, /api/habits/reminders, /api/v2/scheduling/habits
- [ ] `internal/api/chores.go` — /api/v2/chores
- [ ] `internal/api/compass.go` — /api/v2/compass
- [ ] `internal/api/calendar.go` — /api/v2/calendar
- [ ] `internal/api/timers.go` — /api/v2/timers
- [ ] `internal/api/pomodoro.go` — /api/v1/pomodoro
- [ ] `internal/api/lists.go` — /api/v1/households/.../grocery-list
- [ ] `internal/api/projects.go` — /api/project
- [ ] `internal/api/planning.go` — /api/ai/chat/stream + /api/v2/unified-tasks (client-side scheduling)
- [ ] `internal/api/search.go` — /api/v2/search
- [ ] `internal/api/profile.go` — /api/v1/customers
- [ ] `internal/api/memory.go` — /api/v1/customers/memories
- [ ] `internal/api/household.go` — /api/v1/households
- [ ] `internal/api/comments.go` — /api/v2/unified-tasks/{id}/comments
- [ ] `internal/api/sharing.go` — /api/v2/unified-tasks/{id}/share
- [ ] `internal/api/notifications.go` — /api/v1/notifications
- [ ] `internal/api/gamification.go` — /api/v1/gamification
- [ ] `internal/api/export.go` — /api/v1/customers/exports
- [ ] `internal/api/account.go` — /api/v1/account-deletion, /api/payments, /api/v1/usage
- [ ] `internal/api/apikeys.go` — /api/v1/api-keys
- [ ] `internal/api/ai.go` — /api/ai/chat/stream (SSE), /api/v1/ai/conversations

### App Domain Layer (Spec §1.2, Appendix H — 24 files)
- [ ] One file per domain (tasks.go, habits.go, etc.) replacing stubs in app.go
- [ ] Each file calls domain API layer methods and formats output

### CLI Commands (Spec §4.2-4.29 — 29 groups)
The 4 existing stub groups (task, inbox, now, review) need real implementations.
25 new groups must be created:
- [ ] compass (show, generate, correct, complete, status, history)
- [ ] habit (list, done, skip, streak, chains, schedule, reminders)
- [ ] chore (list, done, schedule, rotation)
- [ ] calendar (show, add, delete, decline, skip)
- [ ] timer (list, start, pomodoro, alarm, cancel, snooze, pause, resume, complete, dismiss, count)
- [ ] grocery (list, add, add-bulk, check, delete, clear, convert)
- [ ] project (list, show, create)
- [ ] plan / schedule / reschedule
- [ ] search
- [ ] whoami / goals / prefs (coaching, notifications, timers)
- [ ] memory (list, show, add, update, search, delete, delete-all, export)
- [ ] household (info, members, invite, leaderboard, challenges)
- [ ] review weekly
- [ ] task comment (list, add, edit, delete, count)
- [ ] task share + shared-inbox
- [ ] chore rotation (show, advance, reset, order)
- [ ] notifications (list, unread, read, read-all, delete)
- [ ] stats / stats pomodoro / stats usage
- [ ] achievements (list, streaks, points, challenges)
- [ ] export (request, list, download, delete, delete-batch)
- [ ] account (info, usage, subscription, billing, delete, mcp-sessions)
- [ ] apikey (list, create, show, update, revoke)
- [ ] ai (chat, conversations CRUD)
- [ ] timer pomodoro (smart, current, pause, resume, interrupt, suggest, stop, complete, history, settings)
- [ ] config (show, set, get, reset, path)
- [ ] completion (bash, zsh, fish) + man

### Global Flags (Spec §4.1)
- [ ] `--api-url` flag on root command
- [ ] `--debug` flag on root command

### Output (Spec §9 — 4 files)
- [ ] `internal/output/table.go` — column-aligned tables
- [ ] `internal/output/color.go` — ANSI color by priority zone (red/yellow/blue/gray), with --no-color
- [ ] `internal/output/markdown.go` — Glamour rendering for compass, goals, comments
- [ ] `internal/output/progress.go` — progress bars, SSE streaming display

### TUI (Spec §5 — 20 screens, 16 components)
- [ ] `internal/tui/app.go` — Root Bubble Tea model, screen router, tab management
- [ ] 17 screen files in `internal/tui/screens/`
- [ ] 16 component files in `internal/tui/components/`

### Plugin System (Spec §10)
- [ ] `plugins/plugin.go` — interface, loading, command injection
- [ ] Plugin YAML state file support

### Error Handling (Spec §13)
- [ ] 6 exit codes (0=success, 1=general, 2=usage, 3=auth, 4=network, 5=API, 6=rate-limited)
- [ ] JSON error wrapping
- [ ] 401 auto-refresh middleware

### Integration Tests (Spec §14 — 22 test files)
- [ ] 21 missing test files (auth, tasks, habits, chores, calendar, compass, timers, pomodoro, grocery, projects, planning, search, profile, memory, household, review, comments, sharing, notifications, stats, achievements, export, account, apikeys, ai, cli_output)

### Distribution (Spec §15)
- [ ] `.goreleaser.yaml`
- [ ] `.github/workflows/` CI pipeline
- [ ] `go mod vendor` for reproducible builds

---

## Specialist Feedback History

- **[2026-03-09T22:56Z] review-agent → CHANGES-REQUESTED** — 16 blockers in initial code
- **[2026-03-10T01:30Z] review-agent → CHANGES-REQUESTED** — Incomplete scope, wrong OAuth paths
- **[2026-03-10T01:32Z] review-agent → CHANGES-REQUESTED** — Security review findings
- **[2026-03-10T02:00Z] Opus 4.6 audit** — Full spec audit: ~5% coverage, detailed gap analysis

## Previous Fixes Applied (by Kimi agent)

All 16 initial review blockers were addressed: B1-B16 (tests, crypto/rand, client_id,
redirect URI, server shutdown, PBKDF2, dir perms, URL parsing, retry POST, token save,
Clear dead code, global flags wiring, help text, SilenceUsage, context.Context, error returns).

## Architecture Decisions (preserve these)

1. `TokenStore` interface for swappable keyring backends
2. `Formatter` interface for text/JSON output switching
3. Retry-only-idempotent pattern in HTTP client
4. PKCE listener-first-then-register flow (port consistency)
5. Single binary providing both CLI and TUI modes
6. `internal/` packages to minimize public API surface
