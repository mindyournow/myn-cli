# MYN CLI + TUI — Complete Specification

> Binary name: `mynow`
> Project: `myn-cli` (github.com/mindyournow/myn-cli)
> License: MIT

---

## Table of Contents

1. [Architecture](#1-architecture)
2. [Authentication](#2-authentication)
3. [Configuration](#3-configuration)
4. [CLI Command Reference](#4-cli-command-reference)
   - 4.1 Global Flags · 4.2 Task · 4.3 Inbox · 4.4 Now · 4.5 Compass
   - 4.6 Habit · 4.7 Chore · 4.8 Calendar · 4.9 Timer · 4.10 Grocery
   - 4.11 Project · 4.12 Planning · 4.13 Search · 4.14 Profile · 4.15 Memory
   - 4.16 Household · 4.17 Review · 4.18 Plugin · 4.19 Utility
   - 4.20 Task Comments · 4.21 Task Sharing · 4.22 Chore Rotation
   - 4.23 Notifications · 4.24 Stats & Achievements · 4.25 Export
   - 4.26 Account · 4.27 API Keys · 4.28 AI Conversations · 4.29 Extended Pomodoro
5. [TUI Specification](#5-tui-specification)
   - 5.1 Architecture · 5.2 Tab Bar · 5.3 Now · 5.4 Inbox · 5.5 Tasks
   - 5.6 Habits · 5.7 Chores · 5.8 Calendar · 5.9 Timers · 5.10 Grocery
   - 5.11 Task Detail · 5.12 Compass · 5.13 Search · 5.14 Settings · 5.15 Help
   - 5.16 Command Palette · 5.17 Pomodoro Focus · 5.18 Notifications
   - 5.19 Stats & Achievements · 5.20 AI Chat
6. [Global Keybindings (TUI)](#6-global-keybindings-tui)
7. [Search System](#7-search-system)
8. [Help System](#8-help-system)
9. [Output Formatting](#9-output-formatting)
10. [Plugin System](#10-plugin-system)
11. [Shell Completions](#11-shell-completions)
12. [Man Page](#12-man-page)
13. [Error Handling](#13-error-handling)
14. [Integration Testing](#14-integration-testing)
15. [Distribution & Packaging](#15-distribution--packaging)

Appendices: A (API Endpoints) · B (Priority Zones) · C (Date Parsing) ·
D (Duration Parsing) · E (Recurrence Shortcuts) · F (YNAB Plugin) ·
G (Notification Types) · H (Internal File Structure) · I (Request/Response Structures)

---

## 1. Architecture

### 1.1 Binary Structure

Single static binary `mynow` provides both CLI and TUI modes. Default invocation (no subcommand) launches the TUI.

```
mynow              → launches TUI
mynow tui          → launches TUI (explicit)
mynow <command>    → runs CLI command, prints output, exits
```

### 1.2 Internal Layers

```
cmd/mynow/          CLI entry point (Cobra root + subcommands)
internal/
  app/              Application layer — shared by CLI and TUI
    tasks.go        Task CRUD, complete, archive, snooze
    habits.go       Habit complete, skip, streak, chains, schedule, reminders
    chores.go       Chore list, complete, schedule, rotation
    inbox.go        Inbox list, add, process, count
    compass.go      Compass generate, correct, complete, status
    calendar.go     Calendar events CRUD, decline, skip
    timers.go       Countdown + alarm timers
    pomodoro.go     Pomodoro start/smart-start/pause/resume/stop/complete/history/settings
    lists.go        Grocery CRUD, bulk add, check, clear, convert
    projects.go     Project CRUD, move tasks
    planning.go     AI plan, auto-schedule, reschedule
    search.go       Unified search
    profile.go      Whoami, goals, prefs, coaching intensity
    memory.go       Memory CRUD, search, export
    household.go    Household CRUD, members, invites
    comments.go     Task comment CRUD
    sharing.go      Task share, respond, revoke, shared-inbox
    notifications.go Notification list, read, read-all, delete
    stats.go        Productivity stats, pomodoro stats, AI usage
    achievements.go Achievements list, streaks
    export.go       Data export request, list, download, delete
    account.go      Account info, usage, subscription, billing, deletion
    apikeys.go      API key CRUD
    ai.go           AI chat (SSE streaming), conversation CRUD
  api/              HTTP client (one method per API endpoint)
    client.go       Base HTTP, auth injection, retry, SSE reader
    tasks.go        /api/v2/unified-tasks
    habits.go       /api/habits/chains, /api/habits/reminders, /api/v2/scheduling/habits
    chores.go       /api/v2/chores
    compass.go      /api/v2/compass
    calendar.go     /api/v2/calendar
    timers.go       /api/v2/timers
    pomodoro.go     /api/v1/pomodoro
    lists.go        /api/v1/households/.../grocery-list
    projects.go     /api/project
    planning.go     /api/schedules
    search.go       /api/v2/search
    profile.go      /api/v1/customers
    memory.go       /api/v1/customers/memories
    household.go    /api/v1/households
    comments.go     /api/v2/unified-tasks/{id}/comments
    sharing.go      /api/v2/unified-tasks/{id}/share
    notifications.go /api/v2/notifications
    gamification.go /api/v1/gamification
    export.go       /api/v1/customers/exports
    account.go      /api/v1/account-deletion, /api/payments, /api/v1/usage
    apikeys.go      /api/v1/api-keys
    ai.go           /api/ai/chat/stream (SSE), /api/v1/ai/conversations
  auth/             OAuth PKCE + device flow + credential storage
    oauth.go        Browser-based PKCE flow
    device.go       Device authorization flow
    keyring.go      Linux Secret Service (GNOME Keyring / KDE Wallet)
    apikey.go       API key storage (alternative to OAuth)
    tokens.go       Token refresh, in-memory access token cache
  config/           Configuration loading
    config.go       XDG config dirs, env vars, YAML config file
  output/           Output formatting
    formatter.go    Text / JSON / table / quiet modes
    table.go        Column-aligned text tables
    color.go        ANSI color support with --no-color
    markdown.go     Glamour-based markdown rendering
    progress.go     Progress bars, streaming output for AI chat
  tui/              Bubble Tea TUI
    app.go          Root Bubble Tea model, screen router, tab management
    screens/        One file per screen
      now.go
      inbox.go
      next_actions.go
      habits.go
      chores.go
      calendar.go
      compass.go
      timers.go
      pomodoro.go   (Pomodoro focus mode — section 5.17)
      grocery.go
      projects.go
      task_detail.go
      search.go
      settings.go
      help.go
      notifications.go (Notifications overlay — section 5.18)
      stats.go        (Stats & Achievements — section 5.19)
      ai_chat.go      (AI Chat screen — section 5.20)
    components/     Reusable TUI components
      task_list.go    Filterable, sortable task list
      task_row.go     Single task row rendering
      priority_badge.go  Priority zone indicator
      streak_bar.go   Habit streak visualization
      timer_display.go  Countdown/alarm display
      pomodoro_ring.go  Pomodoro progress ring with session pip track
      input.go        Text input field
      confirm.go      Confirmation dialog
      toast.go        Transient notification banner
      statusbar.go    Bottom status bar (with notification badge 🔔 N)
      tabs.go         Tab bar navigation
      modal.go        Modal overlay
      calendar_grid.go  Week/month calendar grid
      comment_list.go Markdown comment list with author/timestamp
      progress_bar.go Horizontal ASCII progress bar
      chart_bar.go    ASCII bar chart for stats
      sse_reader.go   SSE stream reader for AI chat streaming
plugins/            Plugin interface
  plugin.go         Plugin loading, registration, command injection
  ynab/             YNAB budget plugin (see Appendix F)
    plugin.go
test/
  integration/      Docker Compose integration tests
```

### 1.3 Data Flow

```
User Input → Bubble Tea (TUI)
          → internal/app (business logic)
          → internal/api (HTTP calls)
          → MYN Backend (HTTPS)
          → Response → output/formatter (CLI) or TUI model update
```

### 1.4 Concurrency

- API calls are sequential by default (one at a time per command)
- TUI may fire background API calls (e.g., refresh timer, polling compass status)
- Context cancellation on Ctrl+C for all in-flight HTTP requests

---

## 2. Authentication

### 2.1 Auth Methods (in priority order)

1. **API Key** (`mynow login --api-key`) — simplest, recommended for scripting
2. **OAuth 2.0 PKCE** (`mynow login`) — browser-based, for interactive use
3. **Device Authorization** (`mynow login --device`) — for headless/SSH environments

### 2.2 API Key Flow

```
mynow login --api-key
> Enter your MYN API key: myn_xxxx_...
✓ Authenticated as John Doe (john@example.com)
  API key stored in keyring.
```

- Validates key via `GET /api/v1/customers/me` (key in `X-API-KEY` header)
- Stores in Linux Secret Service under service=`mynow`, account=`api-key`
- All subsequent requests use `X-API-KEY: myn_...` header (not `Authorization: Bearer`)

### 2.3 OAuth 2.0 PKCE Flow

```
mynow login
> Registering CLI client with MYN...
> Opening browser for authentication...
> Waiting for callback on http://localhost:19283/callback
✓ Authenticated as John Doe (john@example.com)
  Refresh token stored in keyring.
```

**OAuth Endpoints (all under the MYN backend):**

| Step | Method | Path |
|------|--------|------|
| Discover | GET | `/.well-known/oauth-authorization-server` |
| Register client | POST | `/api/mcp/oauth/register` |
| Authorize | GET | `/api/mcp/oauth/authorize` |
| Exchange code | POST | `/api/mcp/oauth/token` (`grant_type=authorization_code`) |
| Refresh token | POST | `/api/mcp/oauth/token` (`grant_type=refresh_token`) |
| Logout/revoke | POST | `/api/mcp/oauth/logout` |

**PKCE parameters:**
- `code_challenge_method=S256` (only S256 supported)
- `code_challenge = Base64URL(SHA256(code_verifier))`
- `scope=mcp`

**Flow details:**
1. Register a dynamic CLI client: `POST /api/mcp/oauth/register` → receives `client_id` (format: `mcp_<16chars>`)
2. Open browser: `GET /api/mcp/oauth/authorize?response_type=code&client_id=...&redirect_uri=http://localhost:PORT/callback&code_challenge=...&code_challenge_method=S256&state=...`
3. Local HTTP server on ephemeral port receives `?code=...&state=...`
4. Exchange: `POST /api/mcp/oauth/token` with `grant_type=authorization_code&code=...&redirect_uri=...&code_verifier=...`
5. Response: `{ access_token, token_type: "Bearer", expires_in: 3600, refresh_token, scope: "mcp" }`
6. Store refresh token in keyring; access token in memory
7. Auto-refresh on 401 using `POST /api/mcp/oauth/token` with `grant_type=refresh_token`

**Token TTLs:** Access token = 1 hour · Refresh token = 10 years (effectively permanent) · Auth codes = 5 minutes

### 2.4 Device Authorization Flow

> **Note:** Device authorization flow is not yet supported by the MYN backend. This section describes the planned behavior.

```
mynow login --device
> Visit: https://mindyournow.com/device
> Enter code: ABCD-1234
> Waiting for authorization... (polling every 5s)
✓ Authenticated as John Doe (john@example.com)
```

- Requests device code from MYN backend
- Polls token endpoint until authorized or expired
- Same token storage as OAuth PKCE

### 2.5 Credential Storage

| Backend | When Used |
|---------|-----------|
| GNOME Keyring | GNOME desktop detected |
| KDE Wallet | KDE desktop detected |
| `pass` (password-store) | `MYNOW_KEYRING=pass` env var |
| Encrypted file | Fallback: `~/.config/mynow/credentials.enc` (AES-256-GCM, key derived from machine ID) |

### 2.6 Session Commands

```
mynow login                     # OAuth PKCE (default)
mynow login --api-key           # API key auth
mynow login --device            # Device authorization
mynow logout                    # Clear all stored credentials
mynow whoami                    # Show current user info
mynow auth status               # Show auth method, token expiry, etc.
mynow auth refresh              # Force token refresh
```

---

## 3. Configuration

### 3.1 Config File

Location: `$XDG_CONFIG_HOME/mynow/config.yaml` (default: `~/.config/mynow/config.yaml`)

```yaml
# MYN CLI configuration
api:
  url: https://api.mindyournow.com  # Override with MYN_API_URL env var
  timeout: 30s
  retries: 3

auth:
  method: api-key                    # api-key | oauth | device
  keyring: auto                      # auto | gnome | kde | pass | file

display:
  color: auto                        # auto | always | never
  date_format: relative              # relative | iso | short (Mar 9) | long (March 9, 2026)
  time_format: 12h                   # 12h | 24h
  default_output: text               # text | json | table

tui:
  theme: dark                        # dark | light | auto
  refresh_interval: 30s              # Background data refresh
  vim_keys: true                     # j/k navigation
  mouse: false                       # Mouse support
  animations: true                   # Completion animations, transitions

defaults:
  priority: OPPORTUNITY_NOW          # Default priority for new tasks
  task_type: TASK                    # Default type for new items
  calendar_days: 7                   # Default calendar lookahead
  habit_schedule_days: 7             # Default habit schedule lookahead
```

### 3.2 Environment Variables

| Variable | Description | Overrides |
|----------|-------------|-----------|
| `MYN_API_URL` | Backend URL | `api.url` |
| `MYN_API_KEY` | API key (avoids keyring) | stored credential |
| `MYNOW_CONFIG` | Config file path | default path |
| `MYNOW_KEYRING` | Keyring backend | `auth.keyring` |
| `NO_COLOR` | Disable color (standard) | `display.color` |
| `MYNOW_DEBUG` | Enable debug logging | — |

### 3.3 Config Commands

```
mynow config show                   # Print resolved config (redacts secrets)
mynow config set <key> <value>      # Set a config value
mynow config get <key>              # Get a config value
mynow config reset                  # Reset to defaults
mynow config path                   # Print config file path
```

---

## 4. CLI Command Reference

### 4.1 Global Flags

Every command supports these flags:

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | `-j` | Output in JSON format |
| `--quiet` | `-q` | Suppress non-essential output |
| `--no-color` | | Disable color output |
| `--api-url <url>` | | Override API URL for this command |
| `--debug` | | Enable debug logging |
| `--help` | `-h` | Show help for command |

### 4.2 Task Commands

#### `mynow task list`

List tasks with filtering and sorting.

```
mynow task list [flags]

Flags:
  --priority <zone>     Filter: critical, opportunity, horizon, parking
  --type <type>         Filter: task, habit, chore (default: all)
  --project <name|id>   Filter by project
  --completed           Include completed tasks (isCompleted=true)
  --archived            Show archived tasks (GET /api/v2/unified-tasks/archived)
  --today               Only tasks with startDate = today
  --overdue             Only tasks past their startDate
  --household           Include household-shared tasks (includeHousehold=true)
  --sort <field>        Sort by: priority, date, title, created (default: priority)
  --reverse             Reverse sort order
  --page <n>            Page number (default: 0)
  --limit <n>           Page size (default: 50, max: 200)

Examples:
  mynow task list                         # All active tasks
  mynow task list --priority critical     # Critical Now only
  mynow task list --today                 # Today's tasks
  mynow task list --type habit            # Habits only
  mynow task list --project "Q1 Planning" # Tasks in a project
  mynow task list --json | jq '.[] | .title'
```

Output format (text):

```
 CRITICAL NOW
   ● Prepare quarterly report          2h    Mar 9    Q1 Planning
   ● Fix production bug                30m   Mar 9    —

 OPPORTUNITY NOW
   ○ Review pull requests              1h    Mar 9    Engineering
   ○ Update team wiki                  45m   Mar 10   —

 OVER THE HORIZON
   ◌ Research new frameworks           —     Mar 15   R&D
```

#### `mynow task add <title>`

Create a new task. The CLI generates a UUID for each task (the API requires a client-provided `id`).

```
mynow task add <title> [flags]

Flags:
  --priority <zone>     Priority: critical, opportunity, horizon, parking (default: from config)
  --date <date>         Start date: today, tomorrow, monday, 2026-03-15 (default: today)
  --duration <dur>      Duration: 30m, 1h, 2h30m
  --project <name|id>   Assign to project
  --description <text>  Description text
  --type <type>         task, habit, chore (default: task)
  --recurrence <rule>   RRULE for habits/chores: daily, weekly, "FREQ=WEEKLY;BYDAY=MO,WE,FR"

Examples:
  mynow task add "Call Sam"
  mynow task add "Prepare report" --priority critical --duration 2h --date today
  mynow task add "Morning meditation" --type habit --recurrence daily --duration 15m
  mynow task add "Take out trash" --type chore --recurrence "FREQ=WEEKLY;BYDAY=TU,FR"
```

#### `mynow task show <id>`

Show detailed task information.

```
mynow task show <id>

Output:
  Title:       Prepare quarterly report
  ID:          550e8400-...
  Type:        TASK
  Priority:    CRITICAL NOW
  Start Date:  2026-03-09
  Duration:    2h
  Project:     Q1 Planning
  Description: Q1 financials and projections
  Created:     2026-03-01
  Status:      Active
```

#### `mynow task edit <id>`

Edit a task's fields.

```
mynow task edit <id> [flags]

Flags:
  --title <text>        New title
  --priority <zone>     New priority
  --date <date>         New start date
  --duration <dur>      New duration
  --project <name|id>   Move to project (use "" to unassign)
  --description <text>  New description

Examples:
  mynow task edit abc123 --priority critical
  mynow task edit abc123 --date tomorrow --duration 1h
  mynow task edit abc123 --project "Home Renovation"
```

#### `mynow task done <id>`

Mark a task as completed.

```
mynow task done <id>

Output:
  ✓ Completed: "Prepare quarterly report"
```

#### `mynow task archive <id>`

Archive a completed task.

```
mynow task archive <id>

Output:
  ✓ Archived: "Prepare quarterly report"
```

#### `mynow task uncomplete <id>`

Undo a task completion (mark as not done).

```
mynow task uncomplete <id>

Output:
  ↩ Uncompleted: "Prepare quarterly report"
```

#### `mynow task delete <id>`

Soft-delete a task (recoverable via `task restore`). Use `--permanent` for irreversible deletion.

```
mynow task delete <id> [flags]

Flags:
  --force       Skip confirmation prompt
  --permanent   Irreversibly delete (no restore possible)

Output:
  ✓ Deleted: "Prepare quarterly report"
    Restore: mynow task restore <id>
```

#### `mynow task restore <id>`

Restore a soft-deleted task.

```
mynow task restore <id>

Output:
  ✓ Restored: "Prepare quarterly report"
```

#### `mynow task batch`

Update multiple tasks at once.

```
mynow task batch --ids <id1,id2,...> [flags]

Flags:
  --ids <list>          Comma-separated task IDs
  --priority <zone>     Set priority for all
  --project <name|id>   Move all to project
  --date <date>         Set start date for all

Examples:
  mynow task batch --ids abc123,def456 --priority opportunity
  mynow task batch --ids abc123,def456,ghi789 --project "Q1 Planning"
```

#### `mynow task snooze <id>`

Reschedule a task to a later date.

```
mynow task snooze <id> [flags]

Flags:
  --date <date>         Target date (default: tomorrow)
  --days <n>            Snooze by N days

Examples:
  mynow task snooze abc123                 # → tomorrow
  mynow task snooze abc123 --days 3        # → 3 days from now
  mynow task snooze abc123 --date monday   # → next Monday
```

#### `mynow task move <id> <project>`

Move a task to a project.

```
mynow task move abc123 "Q1 Planning"
```

### 4.3 Inbox Commands

The inbox is a special view of tasks with no priority assigned, or tasks explicitly in the inbox zone.

#### `mynow inbox list`

```
mynow inbox list

Output:
  Inbox (3 items)
  1. Call Sam                         added 2h ago
  2. Look into new health insurance   added yesterday
  3. Fix leaky faucet                 added Mar 7
```

#### `mynow inbox add <title>`

Quick-add an item to the inbox (no priority, today's date, type=TASK).

```
mynow inbox add "Call Sam"
mynow inbox add "Buy groceries" --description "Need milk and eggs"
```

#### `mynow inbox process`

Interactive processing — walks through each inbox item and asks for priority assignment.

```
mynow inbox process

> "Call Sam" — added 2h ago
  [c]ritical  [o]pportunity  [h]orizon  [p]arking  [s]kip  [d]elete
  > c
  ✓ "Call Sam" → Critical Now

> "Look into new health insurance" — added yesterday
  [c]ritical  [o]pportunity  [h]orizon  [p]arking  [s]kip  [d]elete
  > h
  ✓ "Look into new health insurance" → Over The Horizon

  Processed 2 of 3 items. 1 remaining.
```

#### `mynow inbox count`

```
mynow inbox count
3
```

### 4.4 Now (Focus) Commands

The "Now" view shows what to focus on right now — Critical Now tasks + today's calendar.

#### `mynow now`

```
mynow now

  🎯 NOW — Monday, March 9

  FOCUS
    ● Prepare quarterly report     2h     Q1 Planning
    ● Fix production bug           30m    —

  UPCOMING
    09:00  Team Standup             30m    Conference Room B
    14:00  1:1 with Manager         30m    Zoom

  HABITS DUE
    ◆ Morning meditation           15m    🔥 45-day streak
    ◆ Read 30 minutes              30m    🔥 12-day streak
```

#### `mynow now focus`

Set or show the current focus task.

```
mynow now focus                    # Show current focus
mynow now focus <id>               # Set focus to task
mynow now focus --clear            # Clear focus
```

### 4.5 Compass (Briefing) Commands

#### `mynow compass`

Show the current compass briefing. Uses `GET /api/v2/compass/current`.

```
mynow compass

Output: (renders briefing summary as markdown via Glamour)
```

#### `mynow compass generate`

Generate a new compass briefing. Uses `POST /api/v2/compass/generate`.

```
mynow compass generate [flags]

Flags:
  --type <type>           Briefing type: daily, evening, weekly, weekly-and-daily, on-demand
                          (default: on-demand, max 3/day)
  --async                 Don't wait for result (default: wait, uses sync generation)

Examples:
  mynow compass generate
  mynow compass generate --type evening
```

#### `mynow compass correct`

Submit course corrections to an active compass session.
Uses `POST /api/v2/compass/corrections/apply`.

```
mynow compass correct [flags]

Flags:
  --summary-id <id>       Compass summary ID (default: current session)
  --task <id>             Related task ID
  --decision <type>       accepted | rejected | modified | completed | archived
  --new-date <date>       Modified start date (for modified decision)
  --reason <text>         User reason for correction

Examples:
  mynow compass correct --task abc123 --decision completed --reason "Finished early"
  mynow compass correct --task abc123 --decision modified --new-date tomorrow
```

#### `mynow compass complete`

End the current compass session.

```
mynow compass complete [flags]

Flags:
  --summary <text>        Session summary
  --decisions <list>      Comma-separated key decisions

Examples:
  mynow compass complete --summary "Productive day, cleared all critical items"
```

#### `mynow compass status`

Show current compass session state. Uses `GET /api/v2/compass/status`.

```
mynow compass status

Output:
  Session: active (started 8:30 AM)
  Briefing ID: 550e8400-...
  Pending corrections: 2
  Last briefing: 8:30 AM today
  Auto-generation: enabled
```

#### `mynow compass history`

Show compass briefing history. Uses `GET /api/v2/compass/history`.

```
mynow compass history [--limit <n>]

Output:
  COMPASS HISTORY

  Mar 9  08:30  Daily     ✓ completed   3 corrections applied
  Mar 8  08:15  Daily     ✓ completed   1 correction applied
  Mar 7  09:00  Daily     — no session
  Mar 6  08:45  Evening   ✓ completed
```

### 4.6 Habit Commands

#### `mynow habit list`

```
mynow habit list [flags]

Flags:
  --due                   Only habits due today
  --schedule [days]       Show upcoming schedule (default: 7 days)

Output:
  HABITS
  ◆ Morning meditation    daily     🔥 45    15m    due today
  ◆ Read 30 minutes       daily     🔥 12    30m    due today
  ◆ Gym workout           MWF       🔥 8     1h     due Wed
  ◆ Weekly review          weekly    🔥 3     30m    due Sun
```

#### `mynow habit done <id>`

Complete a habit for today. Uses `POST /api/v2/unified-tasks/{id}/complete`.

```
mynow habit done abc123

Output:
  ✓ Completed: "Morning meditation"
    Streak: 46 days 🔥
```

#### `mynow habit skip <id>`

Skip a habit without breaking the streak. Uses `POST /api/v2/unified-tasks/{id}/skip`.

```
mynow habit skip <id> [flags]

Flags:
  --reason <text>         Reason for skipping
  --date <date>           Date to skip (default: today)

Examples:
  mynow habit skip abc123 --reason "Sick day"
```

#### `mynow habit streak <id>`

Show detailed streak information. Uses `GET /api/v2/unified-tasks/{id}/streak`.

```
mynow habit streak abc123 [--history]

Output:
  Morning meditation
  Current streak:  45 days 🔥
  Longest streak:  120 days
  Total completions: 892
  Last completed:  today, 7:15 AM

  --history flag shows day-by-day grid:
  Mar:  ✓✓✓✓✓✓✓✓✓  (9/9 so far)
  Feb:  ✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓  (28/28)
  Jan:  ✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓✓  (31/31)
```

#### `mynow habit chains`

List habit chains. Uses `GET /api/habits/chains`. Additional chain management:

```
mynow habit chains                          # List all chains
mynow habit chains create <name>            # Create a chain (POST /api/habits/chains)
mynow habit chains add <chain-id> <habit-id>    # Add habit to chain
mynow habit chains remove <chain-id> <habit-id> # Remove habit from chain
mynow habit chains status <chain-id>        # Chain completion status today
mynow habit chains done <chain-id>          # Batch complete all habits in chain

Output (list):
  Morning Routine (4 habits)
    1. Morning meditation    15m
    2. Journal               10m
    3. Exercise              30m
    4. Healthy breakfast     20m
```

#### `mynow habit schedule`

Trigger AI scheduling for habits and show upcoming schedule.
Uses `POST /api/v2/scheduling/habits/schedule?numberOfDays=<n>`.

```
mynow habit schedule [--days <n>]

Output:
  Mon Mar 9
    ◆ Morning meditation    15m    ✓ done
    ◆ Read 30 minutes       30m    ○ pending
    ◆ Gym workout           1h     ○ pending

  Tue Mar 10
    ◆ Morning meditation    15m    ○ pending
    ◆ Read 30 minutes       30m    ○ pending
```

#### `mynow habit reminders`

Manage habit reminders.

```
mynow habit reminders                          # List all reminders
mynow habit reminders <id>                     # Show reminder for habit
mynow habit reminders <id> --enable --time 07:30   # Set reminder
mynow habit reminders <id> --disable           # Disable reminder
```

### 4.7 Chore Commands

#### `mynow chore list`

```
mynow chore list [flags]

Flags:
  --assigned-to <name>    Filter by assignee
  --due                   Only chores due today

Output:
  CHORES — Taylor Family
  ▪ Take out trash       Tue/Fri    Alex     10m    due tomorrow
  ▪ Vacuum living room   weekly     Jordan   30m    due Sat
  ▪ Clean kitchen        daily      Riley    20m    due today
```

#### `mynow chore done <id>`

```
mynow chore done <id> [--note <text>]

Output:
  ✓ Completed: "Take out trash"
    Next due: Friday, March 13
```

#### `mynow chore schedule`

```
mynow chore schedule [--date <date>] [--week]

Output:
  Mon Mar 9
    ▪ Clean kitchen        Riley     20m    ○ pending
  Tue Mar 10
    ▪ Take out trash       Alex      10m    ○ pending
    ▪ Clean kitchen        Jordan    20m    ○ pending
```

### 4.8 Calendar Commands

#### `mynow calendar`

```
mynow calendar [flags]

Flags:
  --date <date>           Specific date (default: today)
  --days <n>              Number of days to show (default: from config)
  --week                  Show current week

Output:
  Monday, March 9
    09:00 - 09:30  Team Standup         Conference Room B
    14:00 - 14:30  1:1 with Manager     Zoom
    (all day)      Mom's Birthday

  Tuesday, March 10
    10:00 - 11:00  Sprint Planning      Google Meet
```

#### `mynow calendar add`

```
mynow calendar add <title> [flags]

Flags:
  --start <datetime>      Start time (required for non-all-day)
  --end <datetime>        End time
  --date <date>           Date for all-day events
  --all-day               All-day event
  --location <text>       Location
  --attendees <emails>    Comma-separated attendee emails
  --description <text>    Description
  --recurrence <rule>     RRULE for recurring events

Examples:
  mynow calendar add "Lunch with Alex" --start "2026-03-10T12:00" --end "2026-03-10T13:00"
  mynow calendar add "Team offsite" --date 2026-03-15 --all-day
```

#### `mynow calendar delete <id>`

```
mynow calendar delete <id> [--force]
```

#### `mynow calendar decline <id>`

Decline a meeting invitation.

```
mynow calendar decline <id>
```

#### `mynow calendar skip <id>`

Mark a meeting as skipped (MYN-specific, not calendar deletion).

```
mynow calendar skip <id>
```

### 4.9 Timer Commands

#### `mynow timer list`

```
mynow timer list

Output:
  ACTIVE TIMERS
  ⏱  Focus time           COUNTDOWN   14:32 remaining   RUNNING
  ⏱  Deep work block      POMODORO    Session 2/4       WORK PHASE  18:45 remaining
```

#### `mynow timer start`

```
mynow timer start <duration> [flags]

Flags:
  --label <text>          Timer label
  --type <type>           countdown (default), pomodoro

Duration formats: 25m, 1h, 1h30m, 90s

Examples:
  mynow timer start 25m                           # 25-minute countdown
  mynow timer start 25m --label "Focus time"
```

#### `mynow timer pomodoro`

Start a Pomodoro session.

```
mynow timer pomodoro [flags]

Flags:
  --work <dur>            Work phase (default: 25m)
  --break <dur>           Short break (default: 5m)
  --long-break <dur>      Long break (default: 15m)
  --sessions <n>          Sessions before long break (default: 4)
  --auto-start            Auto-start next phase
  --label <text>          Session label

Examples:
  mynow timer pomodoro                            # Standard 25/5/15/4
  mynow timer pomodoro --work 50m --break 10m     # Custom durations
```

#### `mynow timer alarm`

Set an alarm.

```
mynow timer alarm <time> [flags]

Flags:
  --label <text>          Alarm label
  --recurrence <rule>     RRULE for repeating alarms

Examples:
  mynow timer alarm 07:00 --label "Wake up"
  mynow timer alarm 08:30 --label "Standup" --recurrence "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"
```

#### `mynow timer cancel <id>`

```
mynow timer cancel <id>
```

#### `mynow timer snooze <id>`

```
mynow timer snooze <id> [--minutes <n>]    # default: 5
```

### 4.10 Grocery List Commands

#### `mynow grocery`

```
mynow grocery [list]

Output:
  GROCERY LIST — Taylor Family
  ─────────────────────────────
  Produce
    □ Avocados (4)         ripe ones for guacamole
    ☑ Bananas (1 bunch)

  Dairy
    □ Milk (1 gallon)
    □ Eggs (1 dozen)

  Bakery
    □ Bread

  3 unchecked, 1 checked
```

#### `mynow grocery add <item>`

```
mynow grocery add <item> [flags]

Flags:
  --category <cat>        Category (Produce, Dairy, Meat, etc.)
  --quantity <qty>        Quantity description
  --notes <text>          Notes

Examples:
  mynow grocery add "Avocados" --category Produce --quantity 4 --notes "ripe ones"
  mynow grocery add "Milk" --category Dairy --quantity "1 gallon"
```

#### `mynow grocery add-bulk`

Add multiple items at once.

```
mynow grocery add-bulk [flags]

Flags:
  --items <json>          JSON array of items
  --stdin                 Read items from stdin (one per line, format: "name|category|quantity")

Examples:
  mynow grocery add-bulk --items '[{"name":"Milk","category":"Dairy"},{"name":"Bread","category":"Bakery"}]'
  echo -e "Milk|Dairy|1 gallon\nBread|Bakery|" | mynow grocery add-bulk --stdin
```

#### `mynow grocery check <id>`

Toggle an item's checked state.

```
mynow grocery check <id>
```

#### `mynow grocery delete <id>`

Delete a specific grocery item.

```
mynow grocery delete <id> [--force]
Output:
  ✓ Deleted: "Avocados"
```

API: `DELETE /api/v1/households/{hid}/grocery-list/{id}`

#### `mynow grocery clear`

Clear all checked items from the list.

```
mynow grocery clear [--force]
```

API: `DELETE /api/v1/households/{hid}/grocery-list/checked`

#### `mynow grocery convert`

Convert grocery items to MYN tasks (for a shopping trip).

```
mynow grocery convert [flags]

Flags:
  --priority <zone>       Priority for created tasks (default: opportunity)
  --unchecked-only        Only unchecked items (default: true)
```

### 4.11 Project Commands

#### `mynow project list`

```
mynow project list [--archived]

Output:
  PROJECTS
  ● Q1 Planning           8/12 tasks    2 critical
  ● Home Renovation       3/15 tasks    0 critical
  ● Engineering           12/20 tasks   1 critical
```

#### `mynow project show <name|id>`

```
mynow project show "Q1 Planning"

Output:
  Q1 Planning
  ────────────
  Description: First quarter objectives
  Tasks: 8/12 completed (2 critical)

  ● Prepare quarterly report     CRITICAL    Mar 9
  ● Review budget                CRITICAL    Mar 10
  ○ Update roadmap               OPPORTUNITY Mar 12
  ○ Team feedback survey         OPPORTUNITY Mar 15
  ...
```

#### `mynow project create <name>`

```
mynow project create <name> [flags]

Flags:
  --description <text>    Description
  --color <hex>           Color (#3B82F6)
  --parent <name|id>      Parent project for nesting
```

### 4.12 Planning Commands

#### `mynow plan`

Generate an AI plan for a goal.

```
mynow plan <goal> [flags]

Flags:
  --hours <n>             Available hours
  --deadline <date>       Hard deadline
  --priority <zone>       Priority level

Examples:
  mynow plan "Complete Q1 planning" --hours 4 --deadline 2026-03-15
```

#### `mynow schedule`

Auto-schedule today's tasks.

```
mynow schedule [flags]

Flags:
  --date <date>           Date to schedule (default: today)
  --buffer <min>          Buffer between tasks (default: 15)
  --respect-calendar      Keep existing calendar items (default: true)

Output:
  AUTO-SCHEDULED — Monday, March 9
  09:30 - 11:30  Prepare quarterly report    CRITICAL
  12:00 - 12:30  Review pull requests        OPPORTUNITY
  13:00 - 13:45  Update team wiki            OPPORTUNITY

  UNSCHEDULED (not enough time)
  ◌ Research new frameworks    (moved to tomorrow)
```

#### `mynow reschedule`

Reschedule tasks to a different date.

```
mynow reschedule <id...> [flags]

Flags:
  --date <date>           Target date
  --spread <days>         Spread across N days
  --reason <text>         Reason for rescheduling

Examples:
  mynow reschedule abc123 def456 --date friday --reason "Meeting overran"
```

### 4.13 Search Command

Uses `GET /api/v2/search` with query parameters (not POST with body).
Searchable types: `TASK`, `HABIT`, `CHORE`, `GROCERY`. Results are paginated.

```
mynow search <query> [flags]

Flags:
  --type <types>          Comma-separated: TASK,HABIT,CHORE,GROCERY (default: all)
  --priority <zones>      Comma-separated: CRITICAL,OPPORTUNITY_NOW,OVER_THE_HORIZON,PARKING_LOT
  --status <status>       Comma-separated: pending, completed, incomplete
  --from <date>           Filter by startDate >= date (YYYY-MM-DD)
  --to <date>             Filter by startDate <= date (YYYY-MM-DD)
  --include-archived      Include archived items (default: false)
  --limit <n>             Max results (default: 20, max: 100)
  --offset <n>            Pagination offset (default: 0)

Examples:
  mynow search "quarterly report"
  mynow search "budget" --type task,project --priority critical
  mynow search "meeting notes" --from 2026-01-01 --limit 50

Output:
  SEARCH: "quarterly report" (3 results)

  [task]    Prepare quarterly report       CRITICAL    Mar 9     Q1 Planning
            ...Q1 financials and projections...

  [note]    Q4 Quarterly Report Notes      —           Jan 15    —
            ...revenue increased 12% over...

  [project] Q1 Planning                    —           —         8/12 tasks
            ...First quarter objectives...
```

### 4.14 Profile Commands

#### `mynow whoami`

```
mynow whoami

Output:
  John Doe (john@example.com)
  Timezone: America/New_York
  Subscription: Pro (expires Jan 15, 2027)
  Household: Taylor Family (owner)
  Tasks completed: 1,234
  Current streak: 45 days 🔥
```

#### `mynow goals`

```
mynow goals                                   # Show goals (rendered as markdown)
mynow goals edit                              # Open goals in $EDITOR
mynow goals set <markdown>                    # Set goals from string
```

#### `mynow prefs`

```
mynow prefs                                   # List all preferences
mynow prefs get <key>                          # Get specific preference
mynow prefs set <key> <value>                  # Set a preference

Examples:
  mynow prefs set ai.tone friendly
  mynow prefs get ai.tone
```

#### `mynow prefs coaching`

Get or set the AI coaching intensity. Controls how proactively Kaia nudges you.

```
mynow prefs coaching                           # Show current coaching intensity
mynow prefs coaching <level>                   # Set level: off, gentle, proactive

Output:
  Coaching intensity: GENTLE
  Kaia will offer suggestions when you complete tasks or miss habits.

Examples:
  mynow prefs coaching off          # Disable proactive AI coaching
  mynow prefs coaching gentle       # Subtle suggestions only
  mynow prefs coaching proactive    # Active coaching and reminders
```

API: `GET /api/v1/customers/coaching-intensity` / `PUT /api/v1/customers/coaching-intensity`

### 4.15 Memory Commands

> Note: There is no GET-by-ID endpoint for memories — the backend returns all memories as a list via `GET /api/v1/customers/memories`. `memory show <id>` filters client-side.

```
mynow memory list [--limit <n>]                # List all memories (GET /api/v1/customers/memories)
mynow memory show <id>                         # Show specific memory (client-side filter)
mynow memory add <content> [flags]             # Store a memory (POST /api/v1/customers/memories)
  --category <cat>       user_preference, work_context, personal_info, decision, insight, routine
  --tags <tags>          Comma-separated tags
  --importance <level>   low, medium, high, critical
mynow memory update <id> <content> [flags]     # Update a memory (PUT /api/v1/customers/memories/{id})
mynow memory search <query> [flags]            # Search memories (GET /api/v1/customers/memories/search)
  --category <cat>       Filter by category
  --tag <tag>            Filter by tag
mynow memory delete <id>                       # Delete a specific memory
mynow memory delete-all [--force]              # Delete ALL memories (DELETE /api/v1/customers/memories)
mynow memory export                            # Export memories to file (GET /api/v1/customers/memories/export)
  --output <path>        Output file path (default: ./mynow-memories-<date>.json)
```

### 4.16 Household Commands

```
mynow household                                # Show household info + members
mynow household members                        # List members
mynow household invite <email> [--role <role>]  # Invite a member
mynow household leaderboard [--period <p>]     # Gamification leaderboard (WEEKLY default)
mynow household challenges                     # List active household challenges
```

#### `mynow household leaderboard`

Show the household gamification leaderboard.

```
mynow household leaderboard [flags]

Flags:
  --period <p>      WEEKLY | MONTHLY | ALL_TIME (default: WEEKLY)

Output:
  HOUSEHOLD LEADERBOARD — Taylor Family (This Week)

  #1  John Doe     1,450 pts   🔥 12-day streak   8 tasks · 3 habits
  #2  Alex Smith     920 pts   🔥 5-day streak    5 tasks · 2 habits
  #3  Jordan Lee     340 pts   🔥 2-day streak    2 tasks · 1 habit

API: GET /api/v1/gamification/households/{hid}/leaderboard?period=WEEKLY
```

#### `mynow household challenges`

List active household challenges.

```
mynow household challenges

Output:
  HOUSEHOLD CHALLENGES — Taylor Family

  🏆 Habit Streak Week          All members maintain 7+ day streak     3 days left
     Progress: John ✓  Alex ✓  Jordan ✗
  🏆 Team Task Sprint           Complete 30 tasks this week             5 days left
     Progress: 12 / 30 tasks completed

API: GET /api/v1/gamification/households/{hid}/challenges
```

### 4.17 Review Commands

Daily review workflow.

```
mynow review daily                             # Start interactive daily review
mynow review weekly                            # Start interactive weekly review
```

The daily review walks through:
1. Yesterday's incomplete tasks — snooze, complete, or delete each
2. Today's calendar — show upcoming events
3. Today's habits due — status check
4. Inbox processing — prioritize unprocessed items
5. Compass briefing — generate or show latest

### 4.18 Plugin Commands

```
mynow plugin list                              # List installed plugins
mynow plugin enable <name>                     # Enable a plugin
mynow plugin disable <name>                    # Disable a plugin
mynow plugin info <name>                       # Show plugin details
```

### 4.19 Utility Commands

```
mynow version                                  # Print version info
mynow completion bash|zsh|fish                 # Generate shell completions
mynow man                                      # Generate man page to stdout
```

### 4.20 Task Comment Commands

Comments on tasks, habits, and chores support Markdown (max 10,000 characters).
Rate limited to 10 comments per hour per task.

#### `mynow task comment list <task-id>`

```
mynow task comment list <task-id>

Output:
  COMMENTS — Prepare quarterly report (2)

  John Doe  2h ago
  └─ Added the finance data. Spreadsheet is in the shared drive.

  Alice Smith  yesterday
  └─ Projections look good. Can you add a risk analysis section?
```

#### `mynow task comment add <task-id> "<content>"`

```
mynow task comment add <task-id> "<content>" [flags]

Flags:
  --stdin             Read content from stdin (for multiline comments)

Examples:
  mynow task comment add abc123 "Added the finance data."
  echo -e "## Summary\n- Revenue up 12%\n- Risk: TBD" | mynow task comment add abc123 --stdin

Output:
  ✓ Comment added.
```

#### `mynow task comment edit <task-id> <comment-id> "<content>"`

Edit your own comment (authors only; task owners can edit any comment).

```
mynow task comment edit abc123 cmt456 "Updated: also added risk section."
✓ Comment updated.
```

#### `mynow task comment delete <task-id> <comment-id>`

```
mynow task comment delete abc123 cmt456 [--force]
✓ Comment deleted.
```

#### `mynow task comment count <task-id>`

```
mynow task comment count abc123
2
```

---

### 4.21 Task Sharing Commands

Share tasks with household members. Cannot share: chores, habits, recurring tasks, archived tasks.
Share types: `view` (read-only), `edit` (can modify), `delegate` (transfer ownership).

#### `mynow task share <task-id> <member>`

```
mynow task share <task-id> <member-name-or-id> [flags]

Flags:
  --type <type>           view | edit | delegate (default: edit)
  --message <text>        Optional message to the recipient

Examples:
  mynow task share abc123 Alex
  mynow task share abc123 Jordan --type view --message "FYI on this one"
  mynow task share abc123 "Riley Smith" --type delegate

Output:
  ✓ Shared "Prepare quarterly report" with Alex (EDIT access).
    Alex must accept or decline.
```

#### `mynow task share respond <task-id> accept|decline`

```
mynow task share respond abc123 accept
mynow task share respond abc123 decline [--note "Too busy right now"]

Output (accept):
  ✓ Accepted: "Prepare quarterly report" (shared by John Doe, EDIT access).
```

#### `mynow task share revoke <task-id> <member>`

```
mynow task share revoke abc123 Alex
✓ Revoked Alex's access to "Prepare quarterly report".
```

#### `mynow task share list <task-id>`

List all shares for a task (only task owner).

```
mynow task share list abc123

Output:
  SHARES — Prepare quarterly report
  Member    Type    Status     Shared
  Alex      EDIT    ACCEPTED   2h ago
  Jordan    VIEW    PENDING    5m ago
```

#### `mynow shared-inbox`

Show all tasks shared with you (pending + accepted).

```
mynow shared-inbox [flags]

Flags:
  --pending       Only show pending shares awaiting response
  --accepted      Only show accepted (active shared tasks)

Output:
  SHARED WITH ME

  PENDING (1)
  ? Review budget proposal        from Alice Smith   VIEW    5m ago
    → mynow task share respond <id> accept|decline

  ACCEPTED (2)
  ✓ Fix login page bug            from John Doe      EDIT    accepted yesterday
  ✓ Write unit tests for auth     from Sam Lee       EDIT    accepted Mar 7
```

---

### 4.22 Chore Rotation Commands

Manage rotation assignments for household chores.

#### `mynow chore rotation <chore-id>`

Show current rotation status.

```
mynow chore rotation <chore-id>

Output:
  ROTATION — Take out trash

  Current assignee:  Alex    (position 1 of 3)
  Next assignee:     Jordan
  Rotation order:    Alex → Jordan → Riley → (repeat)
  Last rotated:      March 7, 2026
  Total rotations:   15
  Status:            ACTIVE
```

#### `mynow chore rotation advance <chore-id>`

Advance to the next member in rotation (without marking the chore done).

```
mynow chore rotation advance abc123

Output:
  ✓ Rotation advanced.
    Take out trash — now assigned to Jordan (next: Riley)
```

#### `mynow chore rotation reset <chore-id>`

Reset rotation to the first member.

```
mynow chore rotation reset abc123 [--force]

Output:
  ✓ Rotation reset.
    Take out trash — now assigned to Alex (first in rotation)
```

#### `mynow chore rotation order <chore-id> <member1> <member2> ...`

Update the rotation order.

```
mynow chore rotation order abc123 Riley Jordan Alex [flags]

Flags:
  --preserve-position     Keep the current assignee at their relative position

Output:
  ✓ Rotation order updated: Riley → Jordan → Alex → (repeat)
    Currently assigned: Riley
```

---

### 4.23 Notification Commands

Manage in-app notifications (comments, shares, household events, AI messages, timer completions).
Two versions of the API exist (v1 and v2); the CLI uses v2 for full history with read/actioned state.

#### `mynow notifications`

List recent notifications.

```
mynow notifications [flags]

Flags:
  --unread        Show only unread notifications
  --limit <n>     Number to show (default: 20, max: 100)
  --page <n>      Page number (0-indexed, default: 0)

Output:
  NOTIFICATIONS  (2 unread)

  ● John Doe commented on "Prepare quarterly report"    2h ago    [unread]
  ● Task "Review budget" was shared with you            5h ago    [unread]
  ○ "Morning Routine" chain completed 7 days straight   yesterday
  ○ Jordan accepted your share of "Fix login page"      Mar 7
  ○ Riley joined Taylor Family household                Mar 5

  Showing 5 of 12. Use --page to see more.
```

#### `mynow notifications unread`

Show unread count and the 5 most recent unread notifications.

```
mynow notifications unread

Output:
  Unread: 2

  ● John Doe commented on "Prepare quarterly report"    2h ago
  ● Task "Review budget" was shared with you            5h ago
```

#### `mynow notifications read <id>`

Mark a notification as read.

```
mynow notifications read <id>
✓ Marked as read.
```

#### `mynow notifications read-all`

Mark all notifications as read.

```
mynow notifications read-all
✓ Marked 12 notifications as read.
```

#### `mynow notifications delete <id>`

Delete a notification.

```
mynow notifications delete <id> [--force]
✓ Notification deleted.
```

---

### 4.24 Stats & Achievements Commands

#### `mynow stats`

Overall productivity statistics.

```
mynow stats [flags]

Flags:
  --from <date>       Start of period (default: 30 days ago)
  --to <date>         End of period (default: today)

Output:
  PRODUCTIVITY STATS — Last 30 Days

  Tasks completed:      47   (avg 1.6/day)
  Habits completed:     82   (89% completion rate)
  Focus time:         38.5h  (avg 1.3h/day)
  Pomodoro sessions:    62   (52 completed, 84% completion rate)
  Best day:           Mar 2  (8 tasks + 3.5h focus)
  Current streak:      12 days 🔥
```

#### `mynow stats pomodoro`

Pomodoro-specific statistics.

```
mynow stats pomodoro [flags]

Flags:
  --from <date>       Start of period
  --to <date>         End of period
  --task <id>         Stats for a specific linked task
  --group-by <unit>   day | week | month (default: week)

Output:
  POMODORO STATS — Last 30 Days

  Total sessions:       62   (52 completed, 10 cancelled)
  Completion rate:      84%
  Total focus time:   21.7h
  Total break time:    5.4h
  Total interruptions:  8
  Most productive hour: 10 AM
  Average session:     22.3 min

  WEEKLY BREAKDOWN
  Week of Mar 2:    12 sessions   5.0h focus   2 interruptions
  Week of Feb 24:   15 sessions   6.3h focus   3 interruptions
  Week of Feb 17:    9 sessions   3.8h focus   1 interruption
```

#### `mynow stats usage`

AI token usage and daily request limits.

```
mynow stats usage [flags]

Flags:
  --from <date>       Start of period (default: last 30 days)
  --to <date>         End of period

Output:
  AI USAGE — Last 30 Days

  Total cost:          $1.23
  Total tokens:      185,430  (125,200 input / 60,230 output)
  Total operations:      312  (308 successful, 4 failed)
  Success rate:         98.7%
  Avg cost/operation:   $0.004

  Today:  85 / 100 requests  (85% of daily limit remaining)
```

#### `mynow achievements`

List achievements (unlocked + available with progress).

```
mynow achievements [flags]

Flags:
  --unlocked          Only show unlocked achievements
  --available         Only show locked achievements (with progress)

Output:
  ACHIEVEMENTS — 12 unlocked, 1,450 points

  UNLOCKED
  🏆 First Task           Complete your first task              Mar 1     50 pts
  🔥 Week Warrior         7-day streak on any habit             Mar 8    100 pts
  ⚡ Power Hour           5 Pomodoros in one day                Mar 5    150 pts
  🎯 Inbox Zero           Process inbox to 0                    Mar 3     75 pts

  AVAILABLE
  🌟 Century Club         100-day habit streak       45/100             500 pts
  🚀 Productivity Pro     Complete 1,000 tasks       472/1000           250 pts
  🏠 Team Player          Share 10 tasks             3/10               100 pts
```

#### `mynow achievements streaks`

Show all active streaks.

```
mynow achievements streaks

Output:
  STREAKS

  Morning meditation    45 days   🔥  (longest: 120)
  Daily task review     12 days   🔥  (longest: 23)
  Pomodoro sessions      8 days   🔥  (longest: 15)

API: GET /api/v1/gamification/streaks
```

#### `mynow achievements points`

Show total achievement points.

```
mynow achievements points

Output:
  ACHIEVEMENT POINTS
  Total: 1,450 points across 12 achievements
  Rank in household: #1 of 3 members

API: GET /api/v1/gamification/points
```

---

### 4.25 Export Commands

Export your MYN data in JSON, CSV, or iCal format. Exports are async jobs; download when complete.
Export files are retained for 7 days.

#### `mynow export`

Request a data export.

```
mynow export [flags]

Flags:
  --format <fmt>        json | csv | ical (default: json)
  --include <cats>      Comma-separated: tasks,habits,chores,events,memories
                        (default: all categories)

Examples:
  mynow export                              # Interactive: choose format
  mynow export --format json               # JSON, all data
  mynow export --format ical --include events

Output:
  Requesting export (json, all categories)...
  ✓ Export requested. Job ID: exp-abc123
    Check status: mynow export list
    Download when ready: mynow export download exp-abc123
    Files retained for 7 days.
```

#### `mynow export list`

List all export jobs.

```
mynow export list

Output:
  EXPORTS
  ID           Format  Requested   Status       Note
  exp-abc123   json    Mar 9       COMPLETED    ready to download
  exp-def456   csv     Mar 1       COMPLETED    expires Mar 31
  exp-ghi789   ical    Feb 15      DELETED      —
```

#### `mynow export download <id>`

Download a completed export.

```
mynow export download <id> [flags]

Flags:
  --output <path>       Destination path (default: ./mynow-export-<date>.<ext>)

Examples:
  mynow export download exp-abc123
  mynow export download exp-abc123 --output ~/backups/myn.zip

Output:
  Downloading export exp-abc123...
  ✓ Saved to mynow-export-2026-03-09.zip  (2.3 MB)
```

#### `mynow export delete <id>`

Delete an export job.

```
mynow export delete <id> [--force]
✓ Export exp-abc123 deleted.
```

---

### 4.26 Account Commands

Manage your account, subscription, and billing.

#### `mynow account`

Show account summary.

```
mynow account

Output:
  ACCOUNT

  Name:          John Doe
  Email:         john@example.com
  Member since:  January 15, 2025
  Subscription:  Pro (renews January 15, 2027)
  Household:     Taylor Family (owner, 3 members)
  Auth method:   API Key (GNOME Keyring)
```

#### `mynow account usage`

Show usage vs subscription tier limits.

```
mynow account usage

Output:
  USAGE & LIMITS — Pro Tier

  AI requests today:     85 / 100    ████████░░  85%
  Storage:               45 MB / 1 GB  ░░░░░░░░░   5%
  API keys:               2 / 5       ████░░░░░░  40%
```

#### `mynow account subscription`

Show subscription tier and billing status.

```
mynow account subscription

Output:
  SUBSCRIPTION

  Tier:         Pro
  Status:       Active
  Billing:      Annual ($79.99/year)
  Next billing: January 15, 2027
  Auto-renew:   Yes
```

#### `mynow account billing`

Open the Stripe billing portal in a browser (manage payment method, download invoices, etc.).

```
mynow account billing
Opening billing portal in browser...
```

#### `mynow account delete`

Request account deletion. Sends a confirmation email; account is permanently deleted after
a 30-day grace period.

```
mynow account delete [--immediate]

Output:
  WARNING: All your data will be permanently deleted after a 30-day grace period.
  This cannot be undone.

  Type your email to confirm: john@example.com
  ✓ Deletion requested. Confirmation email sent to john@example.com.
    Cancel at any time: mynow account delete cancel
    Check status: mynow account delete status
```

#### `mynow account delete cancel`

Cancel a pending account deletion.

```
mynow account delete cancel
✓ Deletion cancelled. Your account remains active.
```

#### `mynow account delete status`

Check the deletion status.

```
mynow account delete status

Output:
  Deletion requested:   March 9, 2026
  Scheduled for:        April 8, 2026  (30 days remaining)
  Status:               Awaiting email confirmation
```

---

### 4.27 API Key Commands

Create and manage API keys for programmatic / scripted access to MYN.

Keys are prefixed `myn_`. The full key is shown **only at creation** — it cannot be retrieved again. Scopes restrict what the key can access. Rate limits are per-key.

#### `mynow apikey list`

```
mynow apikey list

Output:
  API KEYS (2)
  ID  Name              Scopes                    Created   Last used   Active
  1   Shell scripts     tasks:list,tasks:view     Mar 1     Mar 9       ✓
  2   Automation        tasks:*,habits:list        Feb 15    Mar 8       ✓
```

#### `mynow apikey create <name>`

```
mynow apikey create <name> [flags]

Flags:
  --scopes <list>        Comma-separated scopes (default: tasks:list,tasks:view)
  --expires <date>       Expiration date: YYYY-MM-DD (default: never)
  --description <text>   Description
  --rate-per-min <n>     Per-minute rate limit (default: 60)
  --rate-per-hour <n>    Per-hour rate limit (default: 1000)

Available scopes:
  tasks:list    tasks:view    tasks:create    tasks:update    tasks:delete    tasks:calendar
  habits:list   habits:view   habits:reminders
  schedules:list  schedules:view  schedules:create  schedules:update  schedules:delete
  projects:list   projects:view   projects:create   projects:update   projects:delete
  user:read     admin:full    agent:full

Examples:
  mynow apikey create "Shell scripts" --scopes "tasks:list,tasks:view"
  mynow apikey create "Full access" --scopes read:all
  mynow apikey create "Temp" --expires 2026-04-01

Output:
  ✓ API key created.

  Name:    Shell scripts
  ID:      1
  Scopes:  tasks:list, tasks:view
  Expires: never

  Key (shown once — copy it now):
  myn_abc123def456ghi789jkl012mno345pqr678stu901

  Use with:  export MYN_API_KEY=myn_abc123...
             mynow login --api-key  (to store in keyring)
```

#### `mynow apikey show <id>`

Show key metadata (not the secret value).

```
mynow apikey show <id>

Output:
  API KEY #1

  Name:              Shell scripts
  Scopes:            tasks:list, tasks:view
  Status:            Active
  Created:           March 1, 2026
  Last used:         March 9, 2026
  Expires:           Never
  Rate limits:       60/min, 1000/hour
  Requests (24h):    45
  Unique IPs (24h):  1
```

#### `mynow apikey update <id>`

```
mynow apikey update <id> [flags]

Flags:
  --name <text>          New display name
  --description <text>   New description
  --scopes <list>        Replace all scopes
  --expires <date>       Set expiration (use "never" to remove expiry)
  --enable               Re-enable a disabled key
  --disable              Disable without deleting
  --rate-per-min <n>
  --rate-per-hour <n>
```

#### `mynow apikey revoke <id>`

Permanently revoke a key. All requests using it immediately receive 401.

```
mynow apikey revoke <id> [--force]
✓ API key "Shell scripts" revoked.
```

---

### 4.28 AI Conversation Commands

Chat with Kaia, the MYN AI assistant. Text-only in the CLI (no voice). Responses stream to
stdout by default. Conversations are saved with a 20-message cap per thread; use
`continue` to extend beyond the cap (prior context is summarized).

#### `mynow ai chat`

Start or continue an interactive chat session.

```
mynow ai chat [flags]

Flags:
  --conversation <id>    Continue an existing conversation
  --task <id>            Link session to a specific task (task context injected)
  --no-stream            Wait for complete response before printing

Examples:
  mynow ai chat                              # New conversation
  mynow ai chat --conversation abc123        # Continue existing
  mynow ai chat --task def456                # Task-focused chat

Session:
  > How should I prioritize today?
  < Kaia: Based on your current tasks and calendar...

  Type your message and press Enter. Type 'quit' or press Ctrl+D to end.
  Conversation saved automatically.
```

#### `mynow ai conversations`

List saved conversations.

```
mynow ai conversations [flags]

Flags:
  --archived      Include archived conversations
  --voice         Include voice-only conversations
  --limit <n>     Number to return (default: 20)

Output:
  CONVERSATIONS (8)

  abc123   Task prioritization         3 messages   Mar 9
  def456   Weekly planning session    12 messages   Mar 7
  ghi789   Habit advice                5 messages   Mar 5  [archived]
```

#### `mynow ai conversations show <id>`

Show all messages in a conversation.

```
mynow ai conversations show <id> [flags]

Flags:
  --last <n>    Show only the last N messages (default: all)

Output:
  CONVERSATION — Task prioritization (Mar 9, 2026, 3 messages)

  You:  How should I prioritize today?

  Kaia: Good morning! Looking at your current tasks, I'd suggest starting
        with "Prepare quarterly report" — it's due today and is marked
        Critical. Your calendar is clear until 9:00 AM, giving you a 90-
        minute focus block before your standup...

  You:  What about the production bug?
```

#### `mynow ai conversations search <query>`

Search conversations by title.

```
mynow ai conversations search <query> [flags]

Flags:
  --limit <n>     Number to return (default: 20)

Output:
  CONVERSATION SEARCH — "planning"

  def456   Weekly planning session    12 messages   Mar 7
  jkl012   Q1 planning discussion      4 messages   Feb 28
```

#### `mynow ai conversations count`

Show conversation statistics.

```
mynow ai conversations count

Output:
  Conversations: 23 total  (21 chat, 2 voice)
```

#### `mynow ai conversations archive <id>`

Archive a conversation (soft-archive; shown with `--archived`).

```
mynow ai conversations archive <id>
✓ Conversation archived.
```

#### `mynow ai conversations favorite <id>`

Mark a conversation as a favorite (pinned in the list).

```
mynow ai conversations favorite <id>
✓ Conversation favorited.

mynow ai conversations unfavorite <id>
✓ Removed from favorites.
```

Both use `PATCH /api/v1/ai/conversations/{id}/status` with `{favorited: true/false}`.

#### `mynow ai conversations delete <id>`

Permanently delete a conversation and all its messages.

```
mynow ai conversations delete <id> [--force]
✓ Conversation deleted.
```

#### `mynow ai conversations continue <id>`

Extend a conversation past the 20-message cap. Prior context is automatically summarized.

```
mynow ai conversations continue <id>

Output:
  Continuing "Task prioritization" — context from prior thread summarized.
  > ...
```

---

### 4.29 Extended Pomodoro Commands

The `timer pomodoro` command (section 4.9) provides a basic start. These subcommands expose
the full `/api/v1/pomodoro` API: smart task selection, pause/resume, session history, and settings.

Note: `timer pomodoro` → `POST /api/v1/pomodoro/start`; this is distinct from the general
timer system at `/api/v2/timers`.

#### `mynow timer pomodoro smart`

Context-aware Pomodoro start: detects available time window from calendar and suggests the
best task to work on.

```
mynow timer pomodoro smart [flags]

Flags:
  --available <min>     Minutes available (default: auto-detect from calendar)
  --suggestions <n>     Number of task suggestions to offer (default: 3)
  --task <id>           Skip suggestions, link to specific task directly

Output:
  SMART POMODORO

  Available time: 85 minutes (next meeting: Team Standup at 10:00 AM)

  Suggested tasks:
  1. Prepare quarterly report    2h total    HIGH  "Most critical, fits window"
  2. Fix production bug          30m total   MED   "Overdue by 1 day"
  3. Review pull requests        1h total    MED   "2 PRs pending"

  Select [1-3] or Enter to start without linking: 1

  ⏱  Starting 25-minute Pomodoro — Prepare quarterly report
     25:00 ━━━━━━━━━━░░░░░░░░░░  (p=pause, s=stop, q=quit display)
```

#### `mynow timer pomodoro current`

Show the currently active Pomodoro session.

```
mynow timer pomodoro current

Output:
  ACTIVE POMODORO

  Task:        Prepare quarterly report
  Type:        WORK — session 2 of 4
  Remaining:   18:45
  Started:     9:30 AM
  Interruptions: 0
```

#### `mynow timer pomodoro pause`

Pause the running Pomodoro.

```
mynow timer pomodoro pause
✓ Paused at 18:45 remaining.
```

#### `mynow timer pomodoro resume`

Resume a paused Pomodoro.

```
mynow timer pomodoro resume
✓ Pomodoro resumed.
```

#### `mynow timer pomodoro stop`

Cancel the current Pomodoro session.

```
mynow timer pomodoro stop [--force]
✓ Pomodoro stopped and cancelled.
```

#### `mynow timer pomodoro complete`

Mark the current Pomodoro as complete before the timer finishes.

```
mynow timer pomodoro complete [--note <text>]
✓ Pomodoro marked complete.
  Next: 5-minute short break (session 2 of 4).
```

#### `mynow timer pomodoro history`

View past Pomodoro sessions.

```
mynow timer pomodoro history [flags]

Flags:
  --from <date>       Start date (default: 7 days ago)
  --to <date>         End date (default: today)
  --status <s>        completed | cancelled (default: all)
  --task <id>         Filter by linked task

Output:
  POMODORO HISTORY

  Mar 9   work    25m  Prepare quarterly report  ✓ completed
  Mar 9   work    25m  Prepare quarterly report  ✓ completed
  Mar 9   break    5m  —                          ✓ completed
  Mar 8   work    25m  Fix production bug         ✗ cancelled (interruption)
  Mar 8   work    25m  Fix production bug         ✓ completed
```

#### `mynow timer pomodoro settings`

View or update Pomodoro settings. With no flags, prints current settings.

```
mynow timer pomodoro settings [flags]

Flags:
  --work <min>              Work phase duration 1-60 (default: 25)
  --short-break <min>       Short break 1-30 (default: 5)
  --long-break <min>        Long break 1-60 (default: 15)
  --sessions <n>            Sessions before long break 1-10 (default: 4)
  --auto-start-breaks       Auto-start break after work phase
  --no-auto-start-breaks
  --auto-start-work         Auto-start work after break
  --no-auto-start-work
  --sound <name>            Completion sound: default|bell|chime|gong|ding|alarm|urgent|none
  --notifications           Enable completion notifications
  --no-notifications

Output (no flags):
  POMODORO SETTINGS

  Work duration:            25 min
  Short break:               5 min
  Long break:               15 min
  Sessions until long break: 4
  Auto-start breaks:         off
  Auto-start work:           off
  Completion sound:          default
  Notifications:             on
```

---

## 5. TUI Specification

### 5.1 Screen Architecture

The TUI uses a tab-based layout with a persistent status bar.

```
┌─────────────────────────────────────────────────────────────────┐
│ ● mynow   Now  Inbox  Tasks  Habits  Chores  Cal  ⏱  🛒  ⚙   │  ← Tab bar
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│                     (Active Screen)                             │  ← Main content
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│ John Doe  ●  3 inbox  2 habits due  ⏱ 14:32  🔔 2  ?=help /=search│  ← Status bar
└─────────────────────────────────────────────────────────────────┘
```

### 5.2 Tab Bar

Tabs and their keyboard shortcuts:

| Tab | Key | Description |
|-----|-----|-------------|
| Now | `1` | Current focus view |
| Inbox | `2` | Unprocessed items |
| Tasks | `3` | All tasks by priority zone |
| Habits | `4` | Habits with streaks |
| Chores | `5` | Household chores |
| Cal | `6` | Calendar view |
| Timers | `7` | Active timers |
| Grocery | `8` | Grocery list |
| Settings | `9` | App settings |

Navigate tabs: `1-9` (direct jump), `Tab`/`Shift+Tab` (cycle), `[`/`]` (prev/next).

### 5.3 Screen: Now (default)

The landing screen. Shows what to focus on right now.

```
┌─────────────────────────────────────────────────────────────────┐
│  🎯 NOW — Monday, March 9, 2026                                │
│                                                                 │
│  CRITICAL NOW                                                   │
│  ► ● Prepare quarterly report          2h     Q1 Planning      │  ← focused item
│    ● Fix production bug                30m    —                 │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│                                                                 │
│  UPCOMING TODAY                                                 │
│    09:00  Team Standup                  30m    Conference Rm B  │
│    14:00  1:1 with Manager              30m    Zoom             │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│                                                                 │
│  HABITS DUE                                                     │
│    ◆ Morning meditation                 15m    🔥 45            │
│    ◆ Read 30 minutes                    30m    🔥 12            │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│                                                                 │
│  COMPASS                                                        │
│  "Productive morning ahead. Focus on the quarterly report..."   │
│  (press g to generate new briefing)                             │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  j/k or ↑/↓    Navigate items
  Enter          Open task detail
  d              Mark done
  s              Snooze
  g              Generate compass briefing
  f              Set as focus task
  n              Quick-add task
```

### 5.4 Screen: Inbox

```
┌─────────────────────────────────────────────────────────────────┐
│  📥 INBOX (3 items)                                             │
│                                                                 │
│  ► Call Sam                                     added 2h ago    │
│    Look into new health insurance               added yesterday │
│    Fix leaky faucet                             added Mar 7     │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  Press p to process inbox, n to add new item                    │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  j/k          Navigate
  Enter        Open detail
  n            Quick-add to inbox
  p            Start processing (interactive prioritization)
  d            Delete item
  c            Set as Critical Now
  o            Set as Opportunity Now
  h            Set as Over the Horizon
  x            Set as Parking Lot

Processing mode (p):
  Shows one item at a time with priority selection:
  c/o/h/x      Assign priority
  s            Skip
  d            Delete
  q            Exit processing
```

### 5.5 Screen: Tasks

```
┌─────────────────────────────────────────────────────────────────┐
│  📋 TASKS                                    [filter: all ▼]    │
│                                                                 │
│  CRITICAL NOW (2)                                               │
│  ► ● Prepare quarterly report     2h    Mar 9   Q1 Planning    │
│    ● Fix production bug           30m   Mar 9   —              │
│                                                                 │
│  OPPORTUNITY NOW (3)                                            │
│    ○ Review pull requests         1h    Mar 9   Engineering    │
│    ○ Update team wiki             45m   Mar 10  —              │
│    ○ Schedule dentist             15m   Mar 10  —              │
│                                                                 │
│  OVER THE HORIZON (2)                                           │
│    ◌ Research new frameworks      —     Mar 15  R&D            │
│    ◌ Plan summer vacation         —     Mar 20  —              │
│                                                                 │
│  PARKING LOT (1)                                                │
│    · Learn Rust                   —     —       —              │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  n=new  d=done  e=edit  s=snooze  /=search  f=filter            │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  j/k          Navigate
  Enter        Open task detail
  n            Create new task
  d            Mark done
  e            Edit task (opens edit modal)
  s            Snooze task
  a            Archive task
  m            Move to project
  /            Search/filter tasks
  f            Toggle filter panel
  c/o/h/x      Quick-set priority (critical/opportunity/horizon/parking)
  Space        Toggle selection (multi-select)
  D            Bulk done (selected)
```

Filter panel (f):

```
  Filter by:
    Priority:  [x] Critical  [x] Opportunity  [x] Horizon  [ ] Parking
    Type:      [x] Tasks     [x] Habits       [x] Chores
    Project:   [All ▼]
    Status:    [x] Active    [ ] Completed     [ ] Archived
    Date:      [All ▼]  Today | This week | Overdue | Custom
```

### 5.6 Screen: Habits

```
┌─────────────────────────────────────────────────────────────────┐
│  🔄 HABITS                                                      │
│                                                                 │
│  TODAY (3 due)                                                  │
│  ► ◆ Morning meditation     15m   🔥 45   ✓ done               │
│    ◆ Read 30 minutes        30m   🔥 12   ○ pending            │
│    ◆ Gym workout            1h    🔥 8    ○ pending            │
│                                                                 │
│  CHAINS                                                         │
│    Morning Routine (4 habits, 3/4 done today)                   │
│    ████████████░░░░  75%                                        │
│                                                                 │
│  7-DAY VIEW                                                     │
│               Mon Tue Wed Thu Fri Sat Sun                       │
│  Meditation    ✓   ✓   ✓   ✓   ✓   ✓   ✓                      │
│  Reading       ✓   ✓   ✓   ✓   ○   -   -                      │
│  Gym           ✓   -   ✓   -   ○   -   -                       │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  d=done  k=skip  r=reminders  c=chains                          │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  j/k          Navigate
  d            Mark done
  K            Skip (preserve streak)
  Enter        Show streak detail
  r            Manage reminders
  c            View chains
  S            Show full schedule
```

### 5.7 Screen: Chores

```
┌─────────────────────────────────────────────────────────────────┐
│  🏠 CHORES — Taylor Family                                     │
│                                                                 │
│  TODAY                                                          │
│  ► ▪ Clean kitchen            Riley    20m    ○ pending         │
│                                                                 │
│  THIS WEEK                                                      │
│    ▪ Take out trash (Tue)     Alex     10m    ○ pending         │
│    ▪ Vacuum living room (Sat) Jordan   30m    ○ pending         │
│    ▪ Laundry (Sun)            Alex     1h     ○ pending         │
│                                                                 │
│  ASSIGNMENTS                                                    │
│    Alex:   3 chores/week  |  Jordan: 2 chores/week             │
│    Riley:  2 chores/week                                        │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  d=done  s=schedule  a=assign                                   │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  j/k          Navigate
  d            Mark done
  Enter        Show chore detail
  s            View full schedule
  a            Change assignment
```

### 5.8 Screen: Calendar

```
┌─────────────────────────────────────────────────────────────────┐
│  📅 CALENDAR — March 2026                    [day|week|agenda]  │
│                                                                 │
│  ◄ Mon 9   Tue 10   Wed 11   Thu 12   Fri 13   Sat 14   Sun ► │
│  ─────────────────────────────────────────────────────────────  │
│  Monday, March 9                                                │
│                                                                 │
│  ► 09:00 - 09:30  Team Standup          Conference Room B      │
│    14:00 - 14:30  1:1 with Manager      Zoom                   │
│    (all day)      Mom's Birthday                                │
│                                                                 │
│  Tuesday, March 10                                              │
│    10:00 - 11:00  Sprint Planning       Google Meet             │
│    15:00 - 16:00  Design Review         Figma                   │
│                                                                 │
│  Wednesday, March 11                                            │
│    09:00 - 09:30  Team Standup          Conference Room B      │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  n=new event  D=decline  S=skip  ←/→=prev/next week            │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  j/k          Navigate events
  Enter        Event detail
  n            Create new event
  D            Decline meeting
  S            Skip meeting
  h/l or ←/→   Previous/next week
  t            Jump to today
  v            Cycle view mode: day → week → agenda
```

### 5.9 Screen: Timers

```
┌─────────────────────────────────────────────────────────────────┐
│  ⏱ TIMERS                                                       │
│                                                                 │
│  ┌───────────────────────────┐                                  │
│  │                           │                                  │
│  │        14 : 32            │   Focus time                     │
│  │        ━━━━━━━━░░░░       │   COUNTDOWN · RUNNING            │
│  │                           │   25m total                      │
│  └───────────────────────────┘                                  │
│                                                                 │
│  ┌───────────────────────────┐                                  │
│  │        18 : 45            │   Deep work block                │
│  │        ━━━━━━━━━━░░░      │   POMODORO · Session 2/4 · WORK │
│  │                           │   25m work / 5m break            │
│  └───────────────────────────┘                                  │
│                                                                 │
│  ALARMS                                                         │
│    07:00  Morning wake-up      daily                            │
│    08:30  Standup reminder     weekdays                         │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  n=new timer  p=pomodoro  a=alarm  Space=pause/resume  x=cancel │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  j/k          Navigate timers
  n            New countdown timer
  p            New pomodoro session
  a            New alarm
  Space        Pause/resume selected timer
  x            Cancel selected timer
  z            Snooze ringing timer
```

### 5.10 Screen: Grocery

```
┌─────────────────────────────────────────────────────────────────┐
│  🛒 GROCERY LIST — Taylor Family                                │
│                                                                 │
│  Produce                                                        │
│  ► □ Avocados (4)                ripe ones for guacamole        │
│    ☑ Bananas (1 bunch)                                          │
│                                                                 │
│  Dairy                                                          │
│    □ Milk (1 gallon)                                            │
│    □ Eggs (1 dozen)                                             │
│                                                                 │
│  Bakery                                                         │
│    □ Bread                                                      │
│                                                                 │
│  Meat                                                           │
│    □ Chicken breast (2 lbs)                                     │
│                                                                 │
│                                              4 unchecked, 1 ✓  │
│  ─────────────────────────────────────────────────────────────  │
│  n=add  Space=check  d=delete  C=clear checked                  │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  j/k          Navigate
  n            Add new item
  Space        Toggle check/uncheck
  d            Delete item
  C            Clear all checked items
  /            Search items
```

### 5.11 Screen: Task Detail (overlay)

Opens as a full-screen overlay when pressing Enter on any task/habit/chore.

```
┌─────────────────────────────────────────────────────────────────┐
│  Prepare quarterly report                              [TASK]   │
│  ═══════════════════════════════════════════════════════════════ │
│                                                                 │
│  Priority:     ● CRITICAL NOW                                   │
│  Start Date:   March 9, 2026                                    │
│  Duration:     2 hours                                          │
│  Project:      Q1 Planning                                      │
│  Status:       Active                                           │
│  Created:      March 1, 2026                                    │
│                                                                 │
│  ─── Description ───────────────────────────────────────────── │
│                                                                 │
│  Q1 financials and projections. Need to gather data from        │
│  the finance team and prepare slides for the board meeting.     │
│                                                                 │
│  ─── Actions ───────────────────────────────────────────────── │
│                                                                 │
│  d=done  e=edit  s=snooze  a=archive  m=move  p=set priority   │
│  Esc/q=back                                                     │
└─────────────────────────────────────────────────────────────────┘
```

For habits, the detail also shows:
- Current streak + longest streak
- Completion history grid
- Chain membership
- Reminder settings

For tasks with comments or shares, the detail adds:

```
  ─── Comments (2) ─────────────────────────────────────────── │
                                                                 │
  John Doe  2h ago                                              │
  └─ Added the finance data. Spreadsheet in shared drive.       │
                                                                 │
  Alice Smith  yesterday                                         │
  └─ Projections look good. Add a risk section?                 │
                                                                 │
  ─── Shared with ─────────────────────────────────────────── │
  Alex     EDIT     ACCEPTED   shared 2h ago                    │
  Jordan   VIEW     PENDING    shared 5m ago                    │
                                                                 │
  C=add comment  T=start Pomodoro for this task                  │
```

Additional keybindings in task detail:
```
  C          Open/add comment (opens text input; Markdown supported)
  Shift+S    Share task with a household member
  T          Start Pomodoro linked to this task
```

### 5.12 Screen: Compass

Accessible via the Now screen or directly via `Shift+G`.

```
┌─────────────────────────────────────────────────────────────────┐
│  🧭 COMPASS BRIEFING — March 9, 2026                           │
│  Session started 8:30 AM                                        │
│                                                                 │
│  ─── Summary ──────────────────────────────────────────────── │
│                                                                 │
│  Good morning! You have a productive day ahead with 2 critical  │
│  tasks and 3 meetings. Your morning is clear until the standup  │
│  at 9:00 — use this window for the quarterly report.            │
│                                                                 │
│  ─── Critical Now ─────────────────────────────────────────── │
│  ● Prepare quarterly report          2h                         │
│  ● Fix production bug                30m                        │
│                                                                 │
│  ─── Opportunity Now ──────────────────────────────────────── │
│  ○ Review pull requests              1h                         │
│                                                                 │
│  ─── Suggestions ──────────────────────────────────────────── │
│  • Block 9:30-11:30 for the quarterly report before meetings    │
│  • Batch the PR reviews after lunch                             │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  g=regenerate  c=submit correction  C=complete session  q=back  │
└─────────────────────────────────────────────────────────────────┘
```

### 5.13 Screen: Search (overlay)

Triggered by `/` from any screen.

```
┌─────────────────────────────────────────────────────────────────┐
│  🔍 Search: quarterly report█                                   │
│  ─────────────────────────────────────────────────────────────  │
│  Types: [all ▼]     Priority: [all ▼]     Status: [all ▼]     │
│  ─────────────────────────────────────────────────────────────  │
│                                                                 │
│  ► [task]    Prepare quarterly report     CRITICAL  Mar 9       │
│              Q1 financials and projections                       │
│                                                                 │
│    [note]    Q4 Quarterly Report Notes    —         Jan 15      │
│              revenue increased 12% over...                       │
│                                                                 │
│    [project] Q1 Planning                  —         8/12 tasks  │
│              First quarter objectives                            │
│                                                                 │
│  3 results                                                      │
│  ─────────────────────────────────────────────────────────────  │
│  Enter=open  Tab=cycle filters  Esc=close                       │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  Type to search (debounced, fires after 300ms idle)
  j/k or ↑/↓   Navigate results
  Enter         Open selected result
  Tab           Cycle through filter dropdowns
  Esc           Close search
```

### 5.14 Screen: Settings

```
┌─────────────────────────────────────────────────────────────────┐
│  ⚙ SETTINGS                                                     │
│                                                                 │
│  Account                                                        │
│    User:       John Doe (john@example.com)                      │
│    Auth:       API Key (stored in GNOME Keyring)                │
│    Household:  Taylor Family (owner)                            │
│                                                                 │
│  Display                                                        │
│  ► Theme:       dark ▼                                          │
│    Date format: relative ▼                                      │
│    Time format: 12h ▼                                           │
│    Animations:  on ▼                                            │
│                                                                 │
│  API                                                            │
│    Backend:     https://api.mindyournow.com                     │
│    Timeout:     30s                                             │
│                                                                 │
│  Defaults                                                       │
│    Priority:    Opportunity Now ▼                                │
│    Calendar days: 7                                             │
│                                                                 │
│  Actions                                                        │
│    [Logout]  [Clear cache]  [Reset config]                      │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  Enter=change  q=back                                           │
└─────────────────────────────────────────────────────────────────┘
```

### 5.15 Screen: Help (overlay)

Triggered by `?` from any screen.

```
┌─────────────────────────────────────────────────────────────────┐
│  HELP — mynow TUI                                    [Esc=close]│
│  ═══════════════════════════════════════════════════════════════ │
│                                                                 │
│  NAVIGATION                                                     │
│    1-9          Jump to tab (Now, Inbox, Tasks, ...)            │
│    Tab/S-Tab    Cycle tabs                                      │
│    [ / ]        Previous / next tab                             │
│    j / k        Move down / up (or ↑/↓)                        │
│    g / G        First / last item                               │
│    Enter        Open detail view                                │
│    Esc / q      Go back / close overlay                         │
│                                                                 │
│  ACTIONS                                                        │
│    n            New item (context-dependent)                    │
│    d            Mark done                                       │
│    e            Edit                                            │
│    s            Snooze                                          │
│    a            Archive                                         │
│    Space        Toggle check / select                           │
│    c/o/h/x      Set priority (critical/opp/horizon/parking)    │
│                                                                 │
│  SEARCH & FILTER                                                │
│    /            Open search overlay                             │
│    f            Toggle filter panel                             │
│                                                                 │
│  TIMERS                                                         │
│    n            New countdown                                   │
│    p            New pomodoro                                    │
│    Space        Pause / resume                                  │
│                                                                 │
│  OTHER                                                          │
│    ?            Toggle this help                                │
│    :            Command palette                                 │
│    Ctrl+C       Quit                                            │
│                                                                 │
│  Full documentation: mynow --help                               │
│  Man page: man mynow                                            │
└─────────────────────────────────────────────────────────────────┘
```

### 5.16 Command Palette

Triggered by `:` — like vim command mode. Supports Tab completion and fuzzy matching.
History: last 50 commands stored, recalled with `↑`/`↓`.

```
┌─────────────────────────────────────────────────────────────────┐
│  : █                                                            │
│  ─────────────────────────────────────────────────────────────  │
│  add task ...          Add a new task                           │
│  add habit ...         Add a new habit                          │
│  add chore ...         Add a new chore                          │
│  inbox add ...         Add to inbox                             │
│  search ...            Search all items                         │
│  goto now              Switch to Now screen                      │
│  goto inbox            Switch to Inbox screen                    │
│  goto tasks            Switch to Tasks screen                    │
│  goto habits           Switch to Habits screen                   │
│  goto chores           Switch to Chores screen                   │
│  goto calendar         Switch to Calendar screen                 │
│  goto timers           Switch to Timers screen                   │
│  goto grocery          Switch to Grocery screen                  │
│  goto settings         Switch to Settings screen                 │
│  goto stats            Open Stats & Achievements screen          │
│  goto ai               Open AI Chat screen                       │
│  goto pomodoro         Open Pomodoro focus screen                │
│  compass generate      Generate new compass briefing             │
│  compass correct       Submit compass correction                 │
│  timer 25m             Start 25-min countdown                    │
│  pomodoro              Start Pomodoro (standard 25/5)            │
│  pomodoro smart        Smart Pomodoro (AI task suggestion)       │
│  pomodoro stop         Stop current Pomodoro                     │
│  alarm 07:00           Set alarm for 07:00                       │
│  notifications         Open notifications overlay                │
│  achievements          Open achievements screen                  │
│  ai chat               Open AI chat screen                       │
│  export                Request data export                       │
│  logout                Log out of MYN                            │
│  quit / q              Quit the TUI                              │
└─────────────────────────────────────────────────────────────────┘
```

Fuzzy matching: `: pom` matches `pomodoro`, `: gs` matches `goto settings`, etc.
Argument completion: after `add task `, completes `--priority`, `--date`, etc.

### 5.17 Screen: Pomodoro Focus Mode

Accessed from the Timers screen (`P`), from the command palette (`: pomodoro`), or via
`mynow tui --screen pomodoro`. Full-screen immersive focus mode with session management.

```
┌─────────────────────────────────────────────────────────────────┐
│  🍅 POMODORO — Prepare quarterly report          [Session 2/4]  │
│                                                                 │
│                                                                 │
│                     ┌───────────────────┐                       │
│                     │                   │                       │
│                     │     18 : 45       │                       │
│                     │                   │                       │
│                     │  ████████████░░░  │                       │
│                     │    WORK PHASE     │                       │
│                     └───────────────────┘                       │
│                                                                 │
│                      Space=pause  s=stop                        │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│                                                                 │
│  SESSION PROGRESS                                               │
│  [WORK][BREAK][WORK ►][BREAK][LONG BREAK]                      │
│   ✓     ✓      now    pending  pending                          │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│                                                                 │
│  Today's sessions:  6 completed  ·  150 min focus  ·  0 interrupt│
│                                                                 │
│  SUGGESTIONS (next task)                                        │
│    1. Fix production bug      30m  CRITICAL  "overdue 1 day"   │
│    2. Review pull requests     1h  OPPORTUNITY                  │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  Space=pause/resume  s=stop  c=complete early  i=interrupt      │
│  h=session history  !=settings  q=back                          │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  Space        Pause / resume
  s            Stop (cancel current session)
  c            Complete session early
  i            Record an interruption
  n            Add note to current session
  h            View session history
  !            Open Pomodoro settings
  q / Esc      Return to Timers screen (session keeps running)
```

**Break screen** (automatically shown between work sessions):

```
┌─────────────────────────────────────────────────────────────────┐
│  ☕ BREAK — Short break (session 2 of 4)                        │
│                                                                 │
│                     ┌───────────────────┐                       │
│                     │     04 : 23       │                       │
│                     │  ████████░░░░░░░  │                       │
│                     │   SHORT BREAK     │                       │
│                     └───────────────────┘                       │
│                                                                 │
│  Next up: Work session 3 — Prepare quarterly report             │
│                                                                 │
│  Space=skip break  q=stop session                               │
└─────────────────────────────────────────────────────────────────┘
```

### 5.18 Overlay: Notifications

Triggered from any screen via `N` (Shift+N). Shows the most recent 10 notifications.
Unread count shown in the status bar: `🔔 2` when there are unread notifications.

```
┌─────────────────────────────────────────────────────────────────┐
│  🔔 NOTIFICATIONS (2 unread)                         [N=close]  │
│  ─────────────────────────────────────────────────────────────  │
│                                                                 │
│  ► ● John Doe commented on "Prepare quarterly report"  2h ago  │
│    ● "Review budget" was shared with you               5h ago  │
│    ○ Morning Routine chain: 7 days completed          yesterday │
│    ○ Jordan accepted "Fix login page" share            Mar 7   │
│    ○ Riley joined Taylor Family                        Mar 5   │
│                                                                 │
│  2 unread · 5 shown of 12                                       │
│  ─────────────────────────────────────────────────────────────  │
│  Enter=open related  r=mark read  R=mark all read  d=delete     │
│  j/k=navigate  N or Esc=close                                   │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  j/k          Navigate notifications
  Enter        Navigate to the related item (task, conversation, etc.)
  r            Mark selected notification as read
  R            Mark all notifications as read
  d            Delete selected notification
  N or Esc     Close overlay
```

### 5.19 Screen: Stats & Achievements

Accessible via command palette (`: goto stats`) or `mynow tui --screen stats`.

```
┌─────────────────────────────────────────────────────────────────┐
│  📊 STATS & ACHIEVEMENTS                     [Tab=switch view]  │
│  ─────────────────────────────────────────────────────────────  │
│  [Overview] [Pomodoro] [Habits] [Achievements] [Usage]          │
│                                                                 │
│  OVERVIEW — Last 30 Days                                        │
│                                                                 │
│  Tasks completed:     47   ████████████░░░░  avg 1.6/day       │
│  Habits completed:    82   ██████████████░░  89% rate           │
│  Focus time:        38.5h  █████████░░░░░░░  avg 1.3h/day      │
│  Pomodoros done:      52   ████████████░░░░  84% completion     │
│                                                                 │
│  Current streak:  12 days 🔥                                   │
│  Best day:        Mar 2  (8 tasks + 3.5h focus)                │
│                                                                 │
│  ACHIEVEMENTS (12 unlocked, 1,450 pts)                          │
│  🏆 First Task  🔥 Week Warrior  ⚡ Power Hour  🎯 Inbox Zero  │
│                                                                 │
│  NEXT ACHIEVEMENT                                               │
│  🌟 Century Club — 100-day habit streak                        │
│  Progress: ████████████░░░░░░░░░░  45/100 days                 │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  Tab=switch view  a=all achievements  q=back                    │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  Tab          Cycle between Overview / Pomodoro / Habits / Achievements / Usage
  a            Show full achievements list (unlocked + available)
  j/k          Navigate achievement list
  Enter        Show achievement detail
  q / Esc      Go back
```

**Pomodoro sub-view:**

```
│  POMODORO STATS — Last 30 Days                                  │
│                                                                 │
│  Sessions:   62 total  ·  52 completed (84%)  ·  10 cancelled  │
│  Focus:    21.7h  ·  Break: 5.4h  ·  Interruptions: 8         │
│  Best hour: 10 AM  ·  Avg session: 22.3 min                    │
│                                                                 │
│  WEEKLY BAR CHART                                               │
│  Week of Mar 2:   ██████████████████████  12  sessions  5.0h  │
│  Week of Feb 24:  ████████████████████████████  15  6.3h      │
│  Week of Feb 17:  █████████████████  9  3.8h                  │
```

### 5.20 Screen: AI Chat

Accessible via command palette (`: ai chat`), `mynow tui --screen ai`, or `N` on an
AI-related notification. Provides streaming text chat with Kaia inside the TUI.

```
┌─────────────────────────────────────────────────────────────────┐
│  🤖 KAIA — AI Chat             [Conv: Task prioritization ▼]   │
│  ─────────────────────────────────────────────────────────────  │
│                                                                 │
│  You: How should I prioritize today?                           │
│                                                                 │
│  Kaia: Good morning! Looking at your schedule, I'd suggest     │
│  starting with "Prepare quarterly report" — it's your most     │
│  critical item and your calendar is clear until 9 AM, giving  │
│  you a solid 90-minute window. After your standup, tackle the  │
│  production bug since it's also critical and 30 minutes.       │
│                                                                 │
│  You: What about the production bug?                           │
│                                                                 │
│  Kaia: The production bug is also marked Critical and is       │
│  overdue by 1 day. I'd prioritize it right after the report... │
│                                                                 │
│  ─────────────────────────────────────────────────────────────  │
│  > type your message here...                                   █│
│  ─────────────────────────────────────────────────────────────  │
│  Enter=send  Shift+Enter=newline  Ctrl+L=clear  n=new conv     │
│  h=history  q=back  (3 messages, 17 remaining in thread)       │
└─────────────────────────────────────────────────────────────────┘

Keybindings:
  Enter           Send message
  Shift+Enter     Insert newline (multiline message)
  Ctrl+L          Clear current input
  n               Start a new conversation
  h               Open conversation history picker
  Esc / q         Return to previous screen (conversation auto-saved)

Conversation picker (h):
  Shows list of recent conversations with title and message count.
  Enter to switch, d to delete, a to archive.
```

---

## 6. Global Keybindings (TUI)

### 6.1 Navigation

| Key | Action |
|-----|--------|
| `1`-`9` | Jump to tab by number |
| `Tab` | Next tab |
| `Shift+Tab` | Previous tab |
| `[` / `]` | Previous / next tab |
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `g` | Jump to first item |
| `G` | Jump to last item |
| `Enter` | Open detail / confirm |
| `Esc` | Back / close overlay |
| `q` | Back (from detail), quit (from root) |
| `Ctrl+C` | Quit immediately |

### 6.2 Actions (context-dependent)

| Key | Action |
|-----|--------|
| `n` | New item (task, event, timer, grocery item — depends on screen) |
| `d` | Mark done / complete |
| `e` | Edit selected item |
| `s` | Snooze task |
| `a` | Archive task |
| `Space` | Toggle check / multi-select |
| `D` | Delete (with confirmation) |
| `m` | Move to project |

### 6.3 Priority Quick-Set

| Key | Priority |
|-----|----------|
| `c` | Critical Now |
| `o` | Opportunity Now |
| `h` | Over The Horizon |
| `x` | Parking Lot |

### 6.4 Global Overlays

| Key | Overlay |
|-----|---------|
| `/` | Search |
| `?` | Help |
| `:` | Command palette |
| `N` (Shift+N) | Notifications overlay |

### 6.5 Screen-Specific Keys

Listed in each screen's section above. The status bar always shows context-relevant key hints.

---

## 7. Search System

### 7.1 TUI Search (`/`)

- Opens a full-screen overlay with a search input at the top
- Debounced: fires API call after 300ms of idle typing
- Results grouped by entity type (task, habit, chore, event, project, note, memory)
- Filter chips below the search bar: type, priority, status
- `Enter` opens the selected result in its detail view
- `Esc` closes search and returns to the previous screen
- Search history: last 20 searches stored locally, accessible via `↑` in the search input

### 7.2 CLI Search (`mynow search`)

- One-shot search with results printed to stdout
- Supports `--json` for machine-readable output
- All filter flags available (see Section 4.13)

### 7.3 Screen-Local Filtering

Each list screen (Tasks, Habits, Chores, Grocery) supports `f` to toggle a filter panel. This is client-side filtering of already-loaded data — not an API call.

---

## 8. Help System

### 8.1 CLI Help

Every command and subcommand has `--help` / `-h` with:
- Synopsis (usage pattern)
- Description
- Available flags with defaults
- Examples

```
mynow task add --help

Add a new task to Mind Your Now.

Usage:
  mynow task add <title> [flags]

Flags:
  --priority <zone>     Priority zone: critical, opportunity, horizon, parking
                        (default: from config, usually "opportunity")
  --date <date>         Start date: today, tomorrow, monday, YYYY-MM-DD
                        (default: today)
  --duration <dur>      Duration: 30m, 1h, 2h30m
  --project <name|id>   Assign to a project
  --description <text>  Task description
  --type <type>         Task type: task, habit, chore (default: task)
  --recurrence <rule>   Recurrence rule for habits/chores

Global Flags:
  -j, --json        Output in JSON format
  -q, --quiet       Suppress non-essential output
      --no-color    Disable color output
  -h, --help        Show this help

Examples:
  mynow task add "Call Sam"
  mynow task add "Prepare report" --priority critical --duration 2h
  mynow task add "Meditate" --type habit --recurrence daily --duration 15m
```

### 8.2 TUI Help (`?`)

Full-screen overlay showing all keybindings organized by category. Content is context-aware — shows screen-specific bindings first, then global bindings.

### 8.3 Man Page

Generated from Cobra command tree using `cobra-doc` or custom generator.

```
MYNOW(1)                    User Commands                    MYNOW(1)

NAME
       mynow - Mind Your Now CLI & TUI

SYNOPSIS
       mynow [command] [flags]
       mynow                    (launches TUI)

DESCRIPTION
       A fast, scriptable, Linux-native terminal client for Mind Your
       Now.  Provides both a command-line interface for scripting and
       automation, and an interactive terminal user interface (TUI)
       for daily use.

COMMANDS
       task        Manage tasks (add, list, edit, done, snooze, archive, comment, share)
       inbox       Manage inbox items (add, list, process, count)
       now         Current focus view
       compass     AI-powered daily briefing
       habit       Habit tracking and streaks
       chore       Household chore management (including rotation)
       calendar    Calendar events
       timer       Timers, alarms, and pomodoro
       grocery     Grocery list management
       project     Project management
       plan        AI-powered planning
       schedule    Auto-schedule tasks
       search      Search across all entities
       review      Daily/weekly review workflows
       memory      Memory store and recall
       household   Household and member management
       shared-inbox Tasks shared with you by household members
       notifications In-app notification management
       stats       Productivity statistics and AI usage
       achievements Gamification achievements and streaks
       export      Data export (JSON/CSV/iCal)
       account     Account info, subscription, billing, deletion
       apikey      API key management for programmatic access
       ai          AI conversation management (Kaia chat)
       login       Authenticate with MYN
       logout      Clear credentials
       whoami      Show current user
       config      Manage configuration
       plugin      Manage plugins
       tui         Launch interactive TUI
       version     Show version information
       completion  Generate shell completions
       man         Generate man page

AUTHENTICATION
       mynow login              Browser-based OAuth 2.0 PKCE
       mynow login --api-key    API key authentication
       mynow login --device     Device authorization (headless)

ENVIRONMENT
       MYN_API_URL     Backend API URL
       MYN_API_KEY     API key (overrides stored credential)
       MYNOW_CONFIG    Config file path
       NO_COLOR        Disable color output

FILES
       ~/.config/mynow/config.yaml    Configuration
       ~/.config/mynow/plugins/       Plugin directory

SEE ALSO
       https://mindyournow.com
       https://github.com/mindyournow/myn-cli

AUTHORS
       Mind Your Now Contributors
```

### 8.4 Shell Completions

Generated by Cobra for bash, zsh, and fish:

```
mynow completion bash > /etc/bash_completion.d/mynow
mynow completion zsh > "${fpath[1]}/_mynow"
mynow completion fish > ~/.config/fish/completions/mynow.fish
```

Completions include:
- All commands and subcommands
- Flag names and values (where enumerable)
- Priority zone names
- Task type names
- Date shortcuts (today, tomorrow, monday, etc.)
- Project names (from cache)

---

## 9. Output Formatting

### 9.1 Text Output (default)

- Color-coded priority zones (red=critical, yellow=opportunity, blue=horizon, gray=parking)
- Unicode symbols for task states: `●` (critical), `○` (opportunity), `◌` (horizon), `·` (parking), `✓` (done), `◆` (habit)
- Relative dates ("2h ago", "yesterday", "Mar 7")
- Column-aligned tables for list views
- Streak fire emoji for habits: `🔥 45`

### 9.2 JSON Output (`--json`)

Every command outputs a JSON object/array when `--json` is passed. Structure matches the API response shape with client-side additions:

```json
{
  "tasks": [
    {
      "id": "550e8400-...",
      "title": "Prepare quarterly report",
      "taskType": "TASK",
      "priority": "CRITICAL",
      "startDate": "2026-03-09",
      "duration": "2h",
      "projectName": "Q1 Planning",
      "isCompleted": false
    }
  ],
  "count": 1
}
```

### 9.3 Quiet Output (`--quiet`)

Only essential output. For mutations: just the result status. For queries: just the data (no headers, decorations, hints).

### 9.4 No Color (`--no-color` or `NO_COLOR=1`)

Strips ANSI escape codes. Useful for piping, logging, and accessibility.

### 9.5 Markdown Rendering

Glamour renders markdown content in:
- Task/note descriptions
- Compass briefing summaries
- Goals display
- Help content
- Memory content

---

## 10. Plugin System

### 10.1 Plugin Interface

Plugins are Go shared objects (`.so`) or standalone binaries in `~/.config/mynow/plugins/`.

#### Shared Object Plugin

```go
// Plugin must export these symbols
func Name() string                    // Plugin name
func Version() string                 // Plugin version
func Description() string             // One-line description
func Commands() []*cobra.Command      // CLI commands to register
func TUIScreens() []tui.Screen       // TUI screens to register (optional)
func Init(api *api.Client) error      // Initialize with API client
```

#### Binary Plugin

Standalone binary named `mynow-<name>`. Invoked as `mynow <name> [args]`.

### 10.2 Plugin Directory

```
~/.config/mynow/plugins/
  openclaw.so           # Shared object plugin
  mynow-example         # Binary plugin
  plugins.yaml          # Plugin state (enabled/disabled)
```

### 10.3 Plugin Commands

```
mynow plugin list                  # List all plugins
mynow plugin enable <name>        # Enable
mynow plugin disable <name>       # Disable
mynow plugin info <name>          # Show details
```

### 10.4 Core Client Independence

The core `mynow` binary compiles and runs without any plugins. Plugins are optional extensions — proprietary ones (like OpenClaw) are distributed separately.

---

## 11. Shell Completions

### 11.1 Supported Shells

- **bash**: via `complete` command
- **zsh**: via `compdef` / `_mynow` function
- **fish**: via `complete` command

### 11.2 Installation

```bash
# bash
mynow completion bash | sudo tee /etc/bash_completion.d/mynow

# zsh
mynow completion zsh > "${fpath[1]}/_mynow"
compinit

# fish
mynow completion fish > ~/.config/fish/completions/mynow.fish
```

### 11.3 Dynamic Completions

- Project names fetched from local cache (refreshed on `mynow project list`)
- Priority zones: `critical`, `opportunity`, `horizon`, `parking`
- Task types: `task`, `habit`, `chore`
- Date shortcuts: `today`, `tomorrow`, `monday`, `tuesday`, etc.

---

## 12. Man Page

### 12.1 Generation

Man page generated from Cobra command tree at build time:

```bash
mynow man > mynow.1
```

Installed to `/usr/share/man/man1/mynow.1.gz` via package managers.

### 12.2 Sections

- NAME, SYNOPSIS, DESCRIPTION
- COMMANDS (all subcommands with brief descriptions)
- AUTHENTICATION (all auth methods)
- CONFIGURATION (config file format, env vars)
- ENVIRONMENT (all env vars)
- FILES (config paths, plugin dirs, credential locations)
- EXIT CODES
- EXAMPLES
- SEE ALSO
- AUTHORS

---

## 13. Error Handling

### 13.1 Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Usage error (bad flags, missing args) |
| 3 | Authentication error (not logged in, expired token) |
| 4 | Network error (cannot reach backend) |
| 5 | API error (4xx/5xx from backend) |
| 6 | Rate limited (429) |

### 13.2 Error Output

Errors go to stderr. JSON mode wraps errors:

```json
{
  "error": "Authentication required",
  "code": 3,
  "hint": "Run 'mynow login' to authenticate"
}
```

### 13.3 Retry Behavior

- Network errors: retry up to 3 times with exponential backoff (1s, 2s, 4s)
- 429 (rate limited): wait for `Retry-After` header, then retry once
- 401 (unauthorized): attempt token refresh once, then fail
- 5xx: retry up to 2 times with 2s delay

### 13.4 Offline Behavior

When the backend is unreachable:
- CLI: error message with exit code 4
- TUI: shows cached data (if available) with "offline" indicator in status bar, retries in background

---

## 14. Integration Testing

### 14.1 Test Stack

Docker Compose stack that mirrors production:
- PostgreSQL 16 (Alpine)
- Redis 7 (Alpine)
- MYN Spring Boot backend (built from source)

### 14.2 Test Data Bootstrap

Uses the demo account API to create reproducible test data:

```
POST /api/v1/admin/demo/recreate-account
Header: X-Demo-API-Key: <key>
→ Returns JWT token + account with full sample data
```

Sample data includes: tasks in all priority zones, habits with streaks, chores with assignments, calendar events, conversation history, compass history.

### 14.3 Test Modes

1. **Full Docker** (default): spins up PostgreSQL + Redis + backend
2. **External backend**: point `MYN_TEST_BACKEND_URL` at a running instance
3. **Mock**: HTTP record/replay for unit tests (no Docker needed)

### 14.4 Test Categories

```
test/integration/
  auth_test.go              # Login flows, token refresh, logout
  tasks_test.go             # CRUD, complete, archive, search
  habits_test.go            # Streaks, skip, chains, schedule
  chores_test.go            # List, complete, schedule
  calendar_test.go          # Events CRUD, decline, skip
  compass_test.go           # Generate, correct, complete
  timers_test.go            # Countdown, pomodoro, alarm, snooze
  grocery_test.go           # CRUD, bulk add, check, convert
  projects_test.go          # CRUD, move tasks
  planning_test.go          # Plan, auto-schedule, reschedule
  search_test.go            # Unified search with filters
  profile_test.go           # Whoami, goals, preferences
  memory_test.go            # Store, recall, search, delete
  household_test.go         # Members, invites
  review_test.go            # Daily review workflow
  cli_output_test.go        # JSON output, quiet mode, no-color
  demo_account_test.go      # Bootstrap test data
  setup.go                  # Docker Compose lifecycle
  docker-compose.yml        # Test stack definition
```

### 14.5 CI Integration

- **Unit tests**: every PR (GitHub Actions)
- **Integration tests**: weekly scheduled + manual trigger
- Backend source cloned via deploy token (private repo)

---

## 15. Distribution & Packaging

### 15.1 Build

GoReleaser builds static binaries with version info baked in via ldflags:

```
-X main.version={{.Version}}
-X main.commit={{.ShortCommit}}
-X main.date={{.Date}}
```

CGO disabled for full static linking.

### 15.2 Platforms

| OS | Arch | Binary |
|----|------|--------|
| Linux | amd64 | `mynow_linux_amd64` |
| Linux | arm64 | `mynow_linux_arm64` |

### 15.3 Package Formats

Built by GoReleaser + nfpm:

| Format | Target |
|--------|--------|
| `.tar.gz` | Manual install |
| `.deb` | Debian, Ubuntu, Pop!_OS |
| `.rpm` | Fedora, RHEL, openSUSE |
| `.apk` | Alpine |
| Arch | Arch Linux, Manjaro |

### 15.4 Release Process

1. Tag: `git tag v0.1.0`
2. Push: `git push origin v0.1.0`
3. GitHub Actions runs GoReleaser
4. Binaries + packages uploaded to GitHub Releases
5. Checksums published

### 15.5 Reproducible Builds

- All dependencies vendored (`go mod vendor`)
- Build flags ensure deterministic output
- Checksums published alongside binaries
- Source fully open — buildable with standard Go toolchain

---

## Appendix A: API Endpoint Mapping

Complete mapping of CLI commands to MYN API endpoints.

| Command | Method | Endpoint |
|---------|--------|----------|
| `task list` | GET | `/api/v2/unified-tasks` (params: `type`, `date`, `includeHousehold`, `isCompleted`, `page`, `size`, `sort`) |
| `task list --archived` | GET | `/api/v2/unified-tasks/archived` |
| `task show <id>` | GET | `/api/v2/unified-tasks/{id}` |
| `task add` | POST | `/api/v2/unified-tasks` (client-generated UUID `id` required) |
| `task edit <id>` | PATCH | `/api/v2/unified-tasks/{id}` |
| `task done <id>` | POST | `/api/v2/unified-tasks/{id}/complete` |
| `task uncomplete <id>` | POST | `/api/v2/unified-tasks/{id}/uncomplete` |
| `task archive <id>` | POST | `/api/v2/unified-tasks/{id}/archive` |
| `task delete <id>` | DELETE | `/api/v2/unified-tasks/{id}` (soft-delete; recoverable) |
| `task delete --permanent <id>` | DELETE | `/api/v2/unified-tasks/{id}/permanent` (irreversible) |
| `task restore <id>` | POST | `/api/v2/unified-tasks/{id}/restore` |
| `task move <id>` | PUT | `/api/project/{projectId}/moveTaskToProject/{taskId}` |
| `task snooze <id>` | PATCH | `/api/v2/unified-tasks/{id}` (update `startDate`) |
| `task batch` | PATCH | `/api/v2/unified-tasks/batch` (body: `{ids: [], updates: {}}`) |
| `inbox list` | GET | `/api/v2/unified-tasks` (filter: no priority) |
| `inbox add` | POST | `/api/v2/unified-tasks` |
| `inbox count` | GET | `/api/v2/unified-tasks` (count: no priority) |
| `habit list` | GET | `/api/v2/unified-tasks?type=HABIT` |
| `habit done <id>` | POST | `/api/v2/unified-tasks/{id}/complete` |
| `habit skip <id>` | POST | `/api/v2/unified-tasks/{id}/skip` |
| `habit streak <id>` | GET | `/api/v2/unified-tasks/{id}/streak` |
| `habit chains` | GET | `/api/habits/chains` |
| `habit chains create` | POST | `/api/habits/chains` |
| `habit chains add <cid> <hid>` | POST | `/api/habits/chains/{chainId}/habits` |
| `habit chains remove <cid> <hid>` | DELETE | `/api/habits/chains/{chainId}/habits/{habitId}` |
| `habit chains status <cid>` | GET | `/api/habits/chains/{chainId}/status` |
| `habit chains done <cid>` | POST | `/api/habits/chains/{chainId}/batch-complete` |
| `habit schedule` | POST | `/api/v2/scheduling/habits/schedule?numberOfDays=<n>` |
| `habit schedule status` | GET | `/api/v2/scheduling/habits/status` |
| `habit reminders` | GET | `/api/habits/reminders` |
| `habit reminders smart-time <id>` | POST | `/api/habits/reminders/{habitId}/calculate-smart-time` |
| `chore list` | GET | `/api/v2/chores` |
| `chore done <id>` | POST | `/api/v2/chores/{id}/complete` |
| `chore schedule` | GET | `/api/v2/chores/schedule` |
| `calendar` | GET | `/api/v2/calendar/events` |
| `calendar add` | POST | `/api/v2/calendar/standalone-events` |
| `calendar delete <id>` | DELETE | `/api/v2/calendar/events/{id}` |
| `calendar decline <id>` | POST | `/api/v2/calendar/meetings/{id}/decline` |
| `calendar skip <id>` | POST | `/api/v2/calendar/meetings/{id}/skip` |
| `compass` | GET | `/api/v2/compass/current` |
| `compass generate` | POST | `/api/v2/compass/generate` (body: `{type: DAILY\|EVENING\|WEEKLY\|ON_DEMAND, sync: bool}`) |
| `compass correct` | POST | `/api/v2/compass/corrections/apply` |
| `compass correct --undo` | POST | `/api/v2/compass/corrections/undo` |
| `compass complete` | POST | `/api/v2/compass/complete` |
| `compass status` | GET | `/api/v2/compass/status` |
| `compass history` | GET | `/api/v2/compass/history` |
| `timer list` | GET | `/api/v2/timers` |
| `timer start` | POST | `/api/v2/timers/countdown` |
| `timer alarm` | POST | `/api/v2/timers/alarm` |
| `timer cancel <id>` | POST | `/api/v2/timers/{id}/cancel` |
| `timer snooze <id>` | POST | `/api/v2/timers/{id}/snooze` |
| `timer pomodoro` | POST | `/api/v1/pomodoro/start` |
| `timer pomodoro smart` | POST | `/api/v1/pomodoro/smart-start` |
| `timer pomodoro current` | GET | `/api/v1/pomodoro/current` |
| `timer pomodoro pause` | POST | `/api/v1/pomodoro/pause` |
| `timer pomodoro resume` | POST | `/api/v1/pomodoro/resume` |
| `timer pomodoro stop` | POST | `/api/v1/pomodoro/stop` |
| `timer pomodoro complete` | POST | `/api/v1/pomodoro/complete` |
| `timer pomodoro history` | GET | `/api/v1/pomodoro/sessions` |
| `timer pomodoro settings` | GET | `/api/v1/pomodoro/settings` |
| `timer pomodoro settings --work` | PUT | `/api/v1/pomodoro/settings` |
| `stats pomodoro` | GET | `/api/v1/pomodoro/stats` |
| `grocery` | GET | `/api/v1/households/{hid}/grocery-list` |
| `grocery add` | POST | `/api/v1/households/{hid}/grocery-list/items` |
| `grocery add-bulk` | POST | `/api/v1/households/{hid}/grocery-list/items/bulk` |
| `grocery check <id>` | PATCH | `/api/v1/households/{hid}/grocery-list/items/{id}` |
| `grocery delete <id>` | DELETE | `/api/v1/households/{hid}/grocery-list/{id}` |
| `grocery clear` | DELETE | `/api/v1/households/{hid}/grocery-list/checked` |
| `grocery convert` | POST | `/api/v1/households/{hid}/grocery-list/convert-to-tasks` |
| `project list` | GET | `/api/project` |
| `project show <id>` | GET | `/api/project/{id}` |
| `project create` | POST | `/api/project/create` |
| `plan` | POST | `/api/schedules/plan` |
| `schedule` | POST | `/api/schedules/auto` |
| `reschedule` | POST | `/api/schedules/reschedule` |
| `search` | GET | `/api/v2/search` (query params: `q`, `types[]`, `statuses[]`, `priorities[]`, `startDate`, `endDate`, `includeArchived`, `limit`, `offset`) |
| `whoami` | GET | `/api/v1/customers/me` |
| `goals` | GET | `/api/v1/customers/goals` |
| `goals set` | PUT | `/api/v1/customers/goals` |
| `prefs` | GET | `/api/v1/customers/preferences` |
| `prefs set` | PUT | `/api/v1/customers/preferences` |
| `prefs coaching` | GET | `/api/v1/customers/coaching-intensity` |
| `prefs coaching <level>` | PUT | `/api/v1/customers/coaching-intensity` (body: `{intensity: OFF\|GENTLE\|PROACTIVE}`) |
| `memory list` | GET | `/api/v1/customers/memories` (returns all; no paginated GET by ID) |
| `memory add` | POST | `/api/v1/customers/memories` |
| `memory update <id>` | PUT | `/api/v1/customers/memories/{id}` |
| `memory search` | GET | `/api/v1/customers/memories/search` |
| `memory delete <id>` | DELETE | `/api/v1/customers/memories/{id}` |
| `memory delete-all` | DELETE | `/api/v1/customers/memories` |
| `memory export` | GET | `/api/v1/customers/memories/export` |
| `household` | GET | `/api/v1/customers/me` (extract households) |
| `household members` | GET | `/api/v1/households/{hid}/members` |
| `household invite` | POST | `/api/v1/households/{hid}/invites` |
| `task comment list <id>` | GET | `/api/v2/unified-tasks/{id}/comments` |
| `task comment add <id>` | POST | `/api/v2/unified-tasks/{id}/comments` |
| `task comment edit <id> <cid>` | PUT | `/api/v2/unified-tasks/{id}/comments/{cid}` |
| `task comment delete <id> <cid>` | DELETE | `/api/v2/unified-tasks/{id}/comments/{cid}` |
| `task comment count <id>` | GET | `/api/v2/unified-tasks/{id}/comments/count` |
| `task share <id>` | POST | `/api/v2/unified-tasks/{id}/share` |
| `task share respond <id>` | POST | `/api/v2/unified-tasks/{id}/share/respond` |
| `task share revoke <id>` | DELETE | `/api/v2/unified-tasks/{id}/share/{memberId}` |
| `task share list <id>` | GET | `/api/v2/unified-tasks/{id}/shares` |
| `shared-inbox` | GET | `/api/v2/unified-tasks/shared-with-me` |
| `chore rotation <id>` | GET | `/api/v2/unified-tasks/{id}/rotation/status` |
| `chore rotation advance <id>` | POST | `/api/v2/unified-tasks/{id}/rotation/advance` |
| `chore rotation reset <id>` | POST | `/api/v2/unified-tasks/{id}/rotation/reset` |
| `chore rotation order <id>` | PUT | `/api/v2/unified-tasks/{id}/rotation/order` |
| `notifications` | GET | `/api/v2/notifications` |
| `notifications unread` | GET | `/api/v2/notifications/unread` |
| `notifications read <id>` | POST | `/api/v2/notifications/mark-read` |
| `notifications read-all` | POST | `/api/v2/notifications/mark-read` (markAll=true) |
| `notifications delete <id>` | DELETE | `/api/v2/notifications/{id}` |
| `stats` | GET | multiple: `/api/v1/gamification/streaks` + `/api/v1/pomodoro/stats` + `/api/v1/usage/today` |
| `stats pomodoro` | GET | `/api/v1/pomodoro/stats` |
| `stats usage` | GET | `/api/v1/token-usage/my-usage` (params: `startDate`, `endDate`) |
| `achievements` | GET | `/api/v1/gamification/achievements` (unlocked) |
| `achievements --available` | GET | `/api/v1/gamification/achievements/available` (locked, with progress) |
| `achievements streaks` | GET | `/api/v1/gamification/streaks` |
| `achievements points` | GET | `/api/v1/gamification/points` |
| `achievements challenges` | GET | `/api/v1/gamification/households/{hid}/challenges` |
| `household leaderboard` | GET | `/api/v1/gamification/households/{hid}/leaderboard` (param: `period=WEEKLY`) |
| `export` | POST | `/api/v1/customers/request-export` |
| `export list` | GET | `/api/v1/customers/exports` |
| `export download <id>` | GET | `/api/v1/customers/exports/{id}/download` |
| `export delete <id>` | DELETE | `/api/v1/customers/exports/{id}` |
| `account` | GET | `/api/v1/customers` |
| `account usage` | GET | `/api/v1/usage/today` |
| `account subscription` | GET | `/api/v1/usage/limits` |
| `account billing` | POST | `/api/payments/create-customer-portal-session` |
| `account delete` | POST | `/api/v1/account-deletion/request` |
| `account delete cancel` | POST | `/api/v1/account-deletion/cancel` |
| `account delete status` | GET | `/api/v1/account-deletion/status` |
| `apikey list` | GET | `/api/v1/api-keys` |
| `apikey create` | POST | `/api/v1/api-keys` |
| `apikey show <id>` | GET | `/api/v1/api-keys/{id}` |
| `apikey update <id>` | PATCH | `/api/v1/api-keys/{id}` |
| `apikey revoke <id>` | DELETE | `/api/v1/api-keys/{id}` |
| `ai chat` | POST | `/api/ai/chat/stream` (SSE streaming) |
| `ai conversations` | GET | `/api/v1/ai/conversations` |
| `ai conversations show <id>` | GET | `/api/v1/ai/conversations/{id}/messages` |
| `ai conversations search <q>` | GET | `/api/v1/ai/conversations/search` |
| `ai conversations count` | GET | `/api/v1/ai/conversations/stats` (returns `totalConversations`, `webConversations`, `voiceConversations`) |
| `ai conversations archive <id>` | PATCH | `/api/v1/ai/conversations/{id}/status` (body: `{isArchived: true}`) |
| `ai conversations favorite <id>` | PATCH | `/api/v1/ai/conversations/{id}/status` (body: `{favorited: true}`) |
| `ai conversations delete <id>` | DELETE | `/api/v1/ai/conversations/{id}` |
| `ai conversations continue <id>` | POST | `/api/v1/ai/conversations/{id}/continue` |

---

## Appendix B: Priority Zones

MYN uses a 4-zone priority system:

| Zone | API Value | CLI Flag | TUI Key | Symbol | Color |
|------|-----------|----------|---------|--------|-------|
| Critical Now | `CRITICAL` | `--priority critical` | `c` | `●` | Red |
| Opportunity Now | `OPPORTUNITY_NOW` | `--priority opportunity` | `o` | `○` | Yellow |
| Over The Horizon | `OVER_THE_HORIZON` | `--priority horizon` | `h` | `◌` | Blue |
| Parking Lot | `PARKING_LOT` | `--priority parking` | `x` | `·` | Gray |
| Inbox | `null` | (no `--priority`) | — | `?` | White |

**Inbox zone**: Tasks with `priority: null` are considered "inbox" items (unprocessed). The API has no
explicit `INBOX` priority constant — a null priority means the item hasn't been assigned yet.
The `inbox list` command filters for `priority == null` client-side after fetching all tasks.

---

## Appendix C: Date Parsing

The CLI accepts multiple date formats:

| Input | Resolved To |
|-------|-------------|
| `today` | Current date |
| `tomorrow` | Current date + 1 |
| `yesterday` | Current date - 1 |
| `monday`..`sunday` | Next occurrence of that day |
| `next week` | Monday of next week |
| `+3d` | 3 days from now |
| `+1w` | 1 week from now |
| `2026-03-15` | Exact ISO date |
| `Mar 15` | March 15 of current year |
| `3/15` | March 15 of current year |

---

## Appendix D: Duration Parsing

| Input | Seconds |
|-------|---------|
| `30s` | 30 |
| `5m` | 300 |
| `25m` | 1500 |
| `1h` | 3600 |
| `1h30m` | 5400 |
| `2h` | 7200 |

---

## Appendix E: Recurrence Rule Shortcuts

| Shortcut | RRULE |
|----------|-------|
| `daily` | `FREQ=DAILY` |
| `weekly` | `FREQ=WEEKLY` |
| `monthly` | `FREQ=MONTHLY` |
| `weekdays` | `FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR` |
| `MWF` | `FREQ=WEEKLY;BYDAY=MO,WE,FR` |
| `TTh` | `FREQ=WEEKLY;BYDAY=TU,TH` |

Full RRULE strings also accepted directly.

---

## Appendix F: YNAB Plugin Specification

The YNAB integration is delivered as an optional plugin (`mynow-ynab` or `ynab.so`).
The MYN backend provides the full YNAB OAuth + budget API at `/api/v1/ynab/`.

### F.1 Authentication

YNAB OAuth 2.0 is handled server-side by the MYN backend:

```
mynow ynab connect          # Open browser → MYN redirects to YNAB OAuth
mynow ynab disconnect        # Revoke YNAB connection
mynow ynab status            # Show connection status + default budget
```

### F.2 Budget Commands

```
mynow ynab budget            # Budget overview (Ready to Assign, Age of Money)
mynow ynab budget accounts   # All accounts grouped by type (checking/savings/credit/loans)
mynow ynab budget categories # All category groups with balances
mynow ynab budget months     # Monthly budget history
```

**Budget overview output:**
```
  YNAB BUDGET — Family Budget

  Ready to Assign:  $234.50
  Age of Money:     18 days

  ACCOUNTS
  Checking (Chase)        $3,421.22    Savings (Ally)    $12,450.00
  Visa Credit Card       -$1,203.45    Student Loan     -$8,500.00

  Net Worth:  $6,167.77
```

### F.3 Transaction Commands

```
mynow ynab transactions [--since <date>]
mynow ynab transactions add [flags]       # Create a transaction
mynow ynab transactions bulk [flags]      # Bulk create from stdin (CSV)
mynow ynab transactions bill <task-id>    # Record bill payment linked to MYN task
```

**Create transaction flags:**
```
  --account <id|name>     Account to post to
  --payee <name>          Payee name
  --category <name>       Budget category
  --amount <dollars>      Amount in dollars (positive=inflow, negative=outflow)
  --date <date>           Transaction date (default: today)
  --memo <text>           Memo
```

### F.4 Analytics Commands

```
mynow ynab analytics spending [--months 3]   # Category spending breakdown
mynow ynab analytics payees [--months 3]     # Top payees by frequency/amount
mynow ynab analytics trends [--months 6]     # Monthly income vs spending
mynow ynab analytics net-worth               # Net worth summary
mynow ynab analytics debt                    # Debt tracking and payoff timeline
```

### F.5 TUI Integration

When the YNAB plugin is enabled, a `YNAB` tab appears in the tab bar (accessible via `0`).

```
┌─────────────────────────────────────────────────────────────────┐
│ ● mynow  Now  Inbox  Tasks  Habits  Chores  Cal  ⏱  🛒  ⚙  💰 │
├─────────────────────────────────────────────────────────────────┤
│  💰 YNAB — Family Budget                                        │
│                                                                 │
│  Ready to Assign:  $234.50    Age of Money:  18 days           │
│                                                                 │
│  [Overview] [Accounts] [Categories] [Transactions] [Analytics] │
│                                                                 │
│  SPENDING THIS MONTH (top categories)                           │
│  Groceries          $342.50  ████████████░░  +$42 over budget  │
│  Restaurants        $156.00  ████████░░░░░░  on track          │
│  Transportation      $89.45  ████░░░░░░░░░░  $60 remaining     │
└─────────────────────────────────────────────────────────────────┘
```

### F.6 API Endpoint Mapping (YNAB Plugin)

| Command | Method | Endpoint |
|---------|--------|----------|
| `ynab connect` | GET | `/api/v1/ynab/authorize` (browser redirect) |
| `ynab disconnect` | POST | `/api/v1/ynab/disconnect` |
| `ynab status` | GET | `/api/v1/ynab/status` |
| `ynab budget` | GET | `/api/v1/ynab/budget/overview` |
| `ynab budget accounts` | GET | `/api/v1/ynab/budget/accounts` |
| `ynab budget categories` | GET | `/api/v1/ynab/budget/categories` |
| `ynab budget months` | GET | `/api/v1/ynab/budget/months` |
| `ynab transactions` | GET | `/api/v1/ynab/transactions` |
| `ynab transactions add` | POST | `/api/v1/ynab/transactions` |
| `ynab transactions bulk` | POST | `/api/v1/ynab/transactions/bulk` |
| `ynab transactions bill <id>` | POST | `/api/v1/ynab/transactions/bill-payment/{taskId}` |
| `ynab analytics spending` | GET | `/api/v1/ynab/analytics/spending` |
| `ynab analytics payees` | GET | `/api/v1/ynab/analytics/payees` |
| `ynab analytics trends` | GET | `/api/v1/ynab/analytics/trends` |
| `ynab analytics net-worth` | GET | `/api/v1/ynab/analytics/net-worth` |
| `ynab analytics debt` | GET | `/api/v1/ynab/analytics/debt` |

---

## Appendix G: Notification Types

Reference for all notification event types that can appear in `mynow notifications`.

| Event Type | Trigger | Related Item |
|------------|---------|--------------|
| `COMMENT_ADDED` | New comment on a task you own or share | task |
| `TASK_SHARED` | A task has been shared with you (pending) | task |
| `TASK_SHARE_ACCEPTED` | Your share invitation was accepted | task |
| `TASK_SHARE_DECLINED` | Your share invitation was declined | task |
| `ACHIEVEMENT_UNLOCKED` | A gamification achievement was unlocked | achievement |
| `HOUSEHOLD_INVITE` | You received a household invitation | household |
| `HOUSEHOLD_MEMBER_JOINED` | Someone joined your household | household |
| `CHORE_ROTATION` | A chore rotation has advanced to you | chore |
| `AI_MESSAGE` | Kaia sent a proactive message | conversation |
| `TIMER_COMPLETE` | A countdown timer completed | timer |
| `ALARM` | An alarm timer fired | timer |
| `TASK_REMINDER` | A task-linked reminder fired | task |

---

## Appendix H: Internal File Structure (complete)

Updated internal package structure including all new features:

```
cmd/mynow/
  main.go                   Version vars, root Cobra command, global flags

internal/
  app/
    tasks.go                Task CRUD, complete, archive, snooze, move
    habits.go               Habit complete, skip, streak, chains, schedule, reminders
    chores.go               Chore list, complete, schedule, rotation
    inbox.go                Inbox list, add, process, count
    compass.go              Compass generate, correct, complete, status
    calendar.go             Calendar events CRUD, decline, skip
    timers.go               Countdown, alarm
    pomodoro.go             Pomodoro start, smart-start, pause, resume, stop, complete,
                            current, history, settings, stats, suggestions
    lists.go                Grocery CRUD, bulk add, check, clear, convert
    projects.go             Project CRUD, move tasks
    planning.go             AI plan, auto-schedule, reschedule
    search.go               Unified search
    profile.go              Whoami, goals, prefs, coaching intensity
    memory.go               Memory CRUD, search, export
    household.go            Household CRUD, members, invites
    comments.go             Task comment CRUD
    sharing.go              Task share, respond, revoke, shared-inbox
    notifications.go        Notification list, read, read-all, delete
    stats.go                Productivity stats, AI usage stats
    achievements.go         Achievements, streaks, leaderboard
    export.go               Data export request, list, download, delete
    account.go              Account info, usage, subscription, billing, deletion
    apikeys.go              API key CRUD
    ai.go                   AI chat (streaming), conversation CRUD
    ynab.go                 YNAB connect, budget, transactions, analytics (plugin)

  api/
    client.go               Base HTTP, auth injection, retry, SSE streaming
    tasks.go                /api/v2/unified-tasks
    habits.go               /api/habits/chains, /api/habits/reminders, /api/v2/scheduling/habits
    chores.go               /api/v2/chores
    compass.go              /api/v2/compass
    calendar.go             /api/v2/calendar
    timers.go               /api/v2/timers
    pomodoro.go             /api/v1/pomodoro
    lists.go                /api/v1/households/.../grocery-list
    projects.go             /api/project
    planning.go             /api/schedules
    search.go               /api/v2/search
    profile.go              /api/v1/customers
    memory.go               /api/v1/customers/memories
    household.go            /api/v1/households
    comments.go             /api/v2/unified-tasks/{id}/comments
    sharing.go              /api/v2/unified-tasks/{id}/share
    notifications.go        /api/v2/notifications
    gamification.go         /api/v1/gamification
    export.go               /api/v1/customers/exports
    account.go              /api/v1/account-deletion, /api/payments, /api/v1/usage
    apikeys.go              /api/v1/api-keys
    ai.go                   /api/ai/chat/stream, /api/v1/ai/conversations
    ynab.go                 /api/v1/ynab (plugin)

  auth/
    oauth.go                Browser PKCE flow
    device.go               Device authorization flow
    keyring.go              Linux Secret Service (GNOME Keyring / KDE Wallet)
    apikey.go               API key credential storage
    tokens.go               Token refresh, in-memory access token cache

  config/
    config.go               XDG config, env vars, YAML config file

  output/
    formatter.go            Text / JSON / table / quiet modes
    table.go                Column-aligned text tables
    color.go                ANSI color support with --no-color
    markdown.go             Glamour markdown rendering
    progress.go             Progress bars, streaming output for AI chat

  tui/
    app.go                  Root Bubble Tea model, screen router, tab management
    screens/
      now.go                Now (focus) screen
      inbox.go              Inbox screen + processing mode
      next_actions.go       Tasks screen with filter panel
      habits.go             Habits screen with 7-day grid
      chores.go             Chores screen with rotation info
      calendar.go           Calendar screen (day/week/agenda)
      compass.go            Compass briefing screen
      timers.go             Timers screen (countdown + alarms)
      pomodoro.go           Pomodoro focus mode screen (5.17)
      grocery.go            Grocery list screen
      projects.go           Project list and detail
      task_detail.go        Task/habit/chore detail overlay
      search.go             Search overlay
      settings.go           Settings screen
      help.go               Help overlay
      notifications.go      Notifications overlay (5.18)
      stats.go              Stats & achievements screen (5.19)
      ai_chat.go            AI chat screen (5.20)
    components/
      task_list.go          Filterable, sortable task list
      task_row.go           Single task row rendering
      priority_badge.go     Priority zone indicator
      streak_bar.go         Habit streak visualization
      timer_display.go      Countdown/pomodoro circular display
      pomodoro_ring.go      Pomodoro progress ring
      input.go              Text input field
      confirm.go            Confirmation dialog
      toast.go              Transient notification banner
      statusbar.go          Bottom status bar (with notification badge)
      tabs.go               Tab bar navigation
      modal.go              Modal overlay
      calendar_grid.go      Week/month calendar grid
      comment_list.go       Task comment list with Markdown rendering
      progress_bar.go       Horizontal progress bar (achievements, usage)
      chart_bar.go          ASCII bar chart (stats screen)
      sse_reader.go         SSE stream reader for AI chat streaming

plugins/
  plugin.go               Plugin interface, loading, command injection
  ynab/                   YNAB plugin (reference implementation)
    plugin.go

test/
  integration/
    setup.go              Docker Compose lifecycle, health wait
    docker-compose.yml    PostgreSQL 16 + Redis 7 + MYN Spring Boot
    demo_account_test.go  Bootstrap test data
    auth_test.go
    tasks_test.go
    habits_test.go
    chores_test.go
    calendar_test.go
    compass_test.go
    timers_test.go
    pomodoro_test.go      Pomodoro start, smart-start, settings, stats
    grocery_test.go
    projects_test.go
    planning_test.go
    search_test.go
    profile_test.go
    memory_test.go
    household_test.go
    review_test.go
    comments_test.go      Task comment CRUD
    sharing_test.go       Task share flows, shared-inbox
    notifications_test.go
    stats_test.go
    achievements_test.go
    export_test.go
    account_test.go
    apikeys_test.go
    ai_test.go            AI chat + conversation history
    cli_output_test.go    JSON, quiet, no-color output modes
```

---

## Appendix I: Key Request / Response Structures

Reference for implementors. Field names are exact as used by the backend.

### I.1 UnifiedTask (create / update)

**POST `/api/v2/unified-tasks`** — client must generate the UUID `id` field.

```json
{
  "id": "<client-generated-uuid>",
  "title": "Prepare quarterly report",
  "taskType": "TASK",
  "priority": "CRITICAL",
  "startDate": "2026-03-09",
  "dueDate": null,
  "duration": 120,
  "description": "Q1 financials and projections",
  "recurrenceRule": null,
  "isAutoScheduled": false,
  "householdId": null
}
```

**PATCH `/api/v2/unified-tasks/{id}`** — all fields optional:

```json
{
  "title": "Updated title",
  "priority": "OPPORTUNITY_NOW",
  "startDate": "2026-03-10",
  "duration": 60
}
```

**Response DTO fields** (key ones):

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Prepare quarterly report",
  "description": "Q1 financials...",
  "taskType": "TASK",
  "priority": "CRITICAL",
  "startDate": "2026-03-09",
  "dueDate": null,
  "duration": 120,
  "isCompleted": false,
  "isArchived": false,
  "recurrenceRule": null,
  "createdDate": "2026-03-01T08:00:00",
  "lastUpdated": "2026-03-09T09:15:00",
  "streakCount": null,
  "commentCount": 2,
  "schedules": [],
  "calendarEvents": []
}
```

### I.2 Unified Task Query Parameters

`GET /api/v2/unified-tasks`

| Param | Type | Description |
|-------|------|-------------|
| `type` | string | `TASK` \| `HABIT` \| `CHORE` \| `RECURRING_TASK` |
| `isCompleted` | boolean | Filter by completion state |
| `includeHousehold` | boolean | Include household shared tasks |
| `page` | int | 0-indexed page (default: 0) |
| `size` | int | Page size (default: 50, max: 200) |
| `sort` | string | e.g. `createdDate,desc` or `title,asc` |

**Inbox filter**: `priority == null` — filter client-side after fetching all tasks, or fetch with no priority filter and check for null `priority` in the response.

### I.3 OAuth Token Request/Response

**POST `/api/mcp/oauth/token`** — authorization code exchange:

```
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code
&code=<authorization_code>
&redirect_uri=http://localhost:PORT/callback
&code_verifier=<pkce_verifier>
&client_id=<mcp_xxxxxxxxxxxxxxxx>
```

**Response:**

```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "rt_...",
  "scope": "mcp"
}
```

**Refresh:**

```
grant_type=refresh_token
&refresh_token=<refresh_token>
```

### I.4 AI Chat Request

**POST `/api/ai/chat/stream`** (SSE streaming):

```json
{
  "currentMessage": "How should I prioritize today?",
  "conversationId": "optional-existing-conversation-id",
  "isMobile": false,
  "isVoice": false
}
```

Response: `text/event-stream` — each `data:` line contains a text chunk. End signaled by `[DONE]`.

### I.5 ConversationDTO

```json
{
  "conversationId": "uuid",
  "customerId": "uuid",
  "title": "Task prioritization",
  "createdAt": "2026-03-09T08:30:00",
  "updatedAt": "2026-03-09T09:15:00",
  "messageCount": 3,
  "isVoice": false,
  "lastMessageAt": "2026-03-09T09:15:00",
  "lastMessagePreview": "Based on your schedule, I'd suggest...",
  "topics": ["prioritization", "planning"],
  "isArchived": false,
  "isDeleted": false,
  "favorited": false
}
```

### I.6 Batch Task Update

**PATCH `/api/v2/unified-tasks/batch`**:

```json
{
  "ids": ["uuid-1", "uuid-2", "uuid-3"],
  "updates": {
    "priority": "OPPORTUNITY_NOW",
    "startDate": "2026-03-10"
  }
}
```

Response: array of updated `UnifiedTaskDTO` objects.

### I.7 Pomodoro Smart-Start Response

**POST `/api/v1/pomodoro/smart-start`** — task suggestions:

```json
{
  "availableMinutes": 85,
  "suggestions": [
    {
      "task": { "id": "...", "title": "Prepare quarterly report", "priority": "CRITICAL" },
      "reason": "Most critical, fits window",
      "estimatedMinutes": 120,
      "confidence": "HIGH"
    }
  ],
  "session": {
    "id": "...",
    "status": "RUNNING",
    "phase": "WORK",
    "sessionNumber": 1,
    "remainingSeconds": 1500
  }
}
```

### I.8 Compass Generate Request/Response

**POST `/api/v2/compass/generate`**:

```json
{
  "type": "ON_DEMAND",
  "sync": true
}
```

Types: `DAILY` | `EVENING` | `WEEKLY` | `WEEKLY_AND_DAILY` | `ON_DEMAND` (max 3/day)

Response (sync, HTTP 200): Compass summary object.
Response (async, HTTP 202): `{ "jobId": "...", "message": "Briefing queued" }`.

### I.9 Chore Rotation Status Response

**GET `/api/v2/unified-tasks/{id}/rotation/status`**:

```json
{
  "choreId": "uuid",
  "choreTitle": "Take out trash",
  "currentAssigneeId": "member-uuid",
  "currentAssigneeName": "Alex",
  "nextAssigneeName": "Jordan",
  "rotationOrder": ["Alex", "Jordan", "Riley"],
  "currentPosition": 0,
  "lastRotated": "2026-03-07T00:00:00",
  "totalRotations": 15,
  "status": "ACTIVE",
  "rotationType": "ROUND_ROBIN"
}
```
