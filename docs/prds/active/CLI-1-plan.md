# CLI-1: Full CLI + TUI Implementation — Plan

## Issue
- **ID:** CLI-1
- **Title:** Full CLI + TUI implementation
- **URL:** https://github.com/mindyournow/myn-cli/issues/1
- **Spec:** `docs/MYN-TUI-CLI_SPEC.md` (4709 lines, single source of truth)

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Scope | Full implementation | All 29 command groups, 20 TUI screens, auth, integration tests, packaging |
| YNAB plugin | Out of scope | Proprietary; separate repo/issue |
| Device auth | Stub only | Backend doesn't support it yet (spec note) |
| Bead granularity | Medium (~30 beads) | Group related commands per domain |
| Priority order | Daily workflow first | Auth → Task → Inbox → Now → Compass → Habit → Calendar → Timer → rest |
| API types | Hand-written per domain | Types in each `api/*.go` file, not a central types package |
| Keyring | zalando/go-keyring | Abstracts GNOME/KDE Secret Service; add encrypted file fallback |
| TUI stack | Full Charm stack | Bubble Tea + Lip Gloss + Glamour + Bubbles |
| Markdown | Glamour | For rendering descriptions, briefings, goals, help |

## Architecture Overview

Greenfield Go project. Existing skeleton has stubs for `cmd/mynow/main.go`, `internal/{api,app,auth,config,output,tui}`, and `test/integration/`.

### Layers (bottom-up build order)
1. **Foundation** — config, auth, API client base, output formatter
2. **API Client** — one file per domain (`api/tasks.go`, `api/habits.go`, etc.)
3. **App Layer** — business logic shared by CLI and TUI (`app/tasks.go`, etc.)
4. **CLI Commands** — Cobra commands in `cmd/mynow/`
5. **TUI Screens** — Bubble Tea screens in `internal/tui/screens/`
6. **TUI Components** — Reusable widgets in `internal/tui/components/`
7. **Infrastructure** — shell completions, man pages, plugin system, packaging

### Key Dependencies
```
go-keyring         — credential storage
cobra              — CLI framework (already in go.mod)
bubbletea          — TUI framework
lipgloss           — TUI styling
glamour            — markdown rendering
bubbles            — TUI components (textinput, list, viewport, etc.)
```

## Implementation Plan (30 Beads)

### Phase 1: Foundation (beads 1-4, 34)
These must be built first — everything else depends on them.

| # | Bead | Difficulty | Dependencies | Files |
|---|------|-----------|--------------|-------|
| 1 | Config system | medium | — | `internal/config/config.go` (rewrite) |
| 2 | Auth (OAuth PKCE + API key + keyring) | complex | 1 | `internal/auth/*.go` |
| 3 | API client base (HTTP, retry, auth injection, SSE) | complex | 1,2 | `internal/api/client.go` (rewrite) |
| 4 | Output formatting (text, JSON, table, color, markdown) | medium | 1 | `internal/output/*.go` |
| 34 | Shared utilities (date/duration/recurrence parsing, global flags, Cobra root) | medium | 1 | `internal/parse/*.go`, `cmd/mynow/root.go` |

### Phase 2: Core CLI Commands (beads 5-14)
Daily workflow commands first, then secondary commands.

| # | Bead | Difficulty | Dependencies | Spec Sections |
|---|------|-----------|--------------|---------------|
| 5 | Task commands (list, add, show, edit, done, delete, restore, archive, snooze, batch, move) | complex | 3,4,34 | 4.2 |
| 6 | Inbox commands (list, add, process, clear, count) | medium | 5 | 4.3 |
| 7 | Now/Focus commands (now, focus, complete, snooze) | medium | 5 | 4.4 |
| 8 | Compass commands (show, generate, correct, complete, status, history) | medium | 3,4 | 4.5 |
| 9 | Habit commands (list, done, skip, streak, chains, schedule, reminders) | complex | 5 | 4.6 |
| 10 | Chore commands (list, done, schedule, stats, rotation) | medium | 5 | 4.7, 4.22 |
| 11 | Calendar commands (list, add, delete, decline, skip) | medium | 3,4 | 4.8 |
| 12 | Timer commands (start, alarm, pause, resume, complete, dismiss, count) | medium | 3,4 | 4.9 |
| 13 | Grocery commands (list, add, add-bulk, check, delete, clear, convert) | medium | 3,4 | 4.10 |
| 14 | Project commands (list, show, create) | medium | 3,4 | 4.11 |

### Phase 3: Secondary CLI Commands (beads 15-22)

| # | Bead | Difficulty | Dependencies | Spec Sections |
|---|------|-----------|--------------|---------------|
| 15 | Planning commands (plan, schedule, reschedule) | medium | 5,8 | 4.12 |
| 16 | Search command | medium | 3,4 | 4.13 |
| 17 | Profile commands (whoami, goals, prefs, coaching, notifications, timers, mcp-sessions) | medium | 3,4 | 4.14 |
| 18 | Memory commands (list, add, show, update, search, delete, export) | medium | 3,4 | 4.15 |
| 19 | Household commands (info, members, invite, leaderboard, challenges) | medium | 3,4 | 4.16 |
| 20 | Task comments + sharing (comment CRUD, share, respond, revoke, shared-inbox) | medium | 5 | 4.20, 4.21 |
| 21 | Notifications + Stats + Achievements | medium | 3,4 | 4.23, 4.24 |
| 22 | Account + API keys + Export + AI conversations | complex | 3,4 | 4.25-4.28 |

### Phase 4: Extended Features (bead 23)

| # | Bead | Difficulty | Dependencies | Spec Sections |
|---|------|-----------|--------------|---------------|
| 23 | Extended Pomodoro (smart, current, pause, resume, interrupt, suggest, stop, complete, history, settings) | medium | 12 | 4.29 |

### Phase 5: TUI (beads 24-29)

| # | Bead | Difficulty | Dependencies | Spec Sections |
|---|------|-----------|--------------|---------------|
| 24 | TUI framework (app shell, tab bar, status bar, keybindings, command palette) | complex | 4 | 5.1-5.2, 5.16, 6 |
| 25 | TUI core screens: Now, Inbox, Tasks, Task Detail | complex | 24,5,6,7 | 5.3-5.5, 5.11 |
| 26 | TUI screens: Habits, Chores, Calendar | complex | 24,9,10,11 | 5.6-5.8 |
| 27 | TUI screens: Timers, Pomodoro Focus, Grocery | complex | 24,12,13,23 | 5.9-5.10, 5.17 |
| 28 | TUI screens: Compass, Search, Notifications, Stats, AI Chat | complex | 24,8,16,21,22 | 5.12-5.13, 5.18-5.20 |
| 29 | TUI screens: Settings, Help + TUI components library | medium | 24 | 5.14-5.15, components |

### Phase 6: Infrastructure (beads 30-33)

| # | Bead | Difficulty | Dependencies | Spec Sections |
|---|------|-----------|--------------|---------------|
| 30 | Review commands (daily, weekly interactive workflow) | medium | 5,8,9,11 | 4.17 |
| 31 | Plugin system + shell completions + man page + utility commands | medium | all CLI | 4.18-4.19, 10-12 |
| 32 | Error handling, retry logic, exit codes | medium | 3 | 13 |
| 33 | Integration tests + CI + packaging (GoReleaser, deb, rpm) | complex | all | 14-15 |

## Dependency Graph (simplified)

```
1 (Config)
├── 2 (Auth) ──┐
│              ├── 3 (API Client) ──┐
├── 4 (Output) ─────────────────────┤
├── 34 (Shared Utils: date/duration/recurrence/global flags) ──┐
│                                   │                          │
│                                   ├── 5 (Task) ──┬── 6 (Inbox)
│                                   │              ├── 7 (Now)
│                                   │              ├── 9 (Habit)
│                                   │              ├── 10 (Chore)
│                                   │              ├── 15 (Planning)
│                                   │              └── 20 (Comments/Sharing)
│                                   ├── 8 (Compass)
│                                   ├── 11 (Calendar)
│                                   ├── 12 (Timer) ── 23 (Extended Pomodoro)
│                                   ├── 13 (Grocery)
│                                   ├── 14 (Project)
│                                   ├── 16 (Search)
│                                   ├── 17 (Profile + MCP sessions)
│                                   ├── 18 (Memory)
│                                   ├── 19 (Household)
│                                   ├── 21 (Notifications/Stats)
│                                   └── 22 (Account/Keys/Export/AI)

24 (TUI Framework) ── 25-29 (TUI Screens)
30-33 (Infrastructure) — after all CLI beads
```

## Risk Assessment

| Risk | Mitigation |
|------|-----------|
| SSE streaming for AI chat is complex | Implement in API client base (bead 3), test early |
| OAuth PKCE browser flow needs local HTTP server | Well-documented pattern; go-keyring handles storage |
| 20 TUI screens is a lot of UI code | Reusable components (bead 29) reduce duplication |
| Integration tests need Docker + backend source | CI runs weekly; local dev can use `MYN_TEST_BACKEND_URL` |
| Spec may have gaps vs actual backend | VERIFIED_CONTROLLERS.md tracks known gaps; test against real backend |

## Out of Scope
- YNAB plugin (Appendix F) — proprietary, separate repo
- Device auth flow — backend doesn't support it yet (stub only)
- Mobile/desktop notifications — Linux only
- Any backend changes
