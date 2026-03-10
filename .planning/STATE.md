# Agent State: CLI-1

## Specialist Feedback

- **[2026-03-09T22:56Z] review-agent → CHANGES-REQUESTED** — `.planning/feedback/001-review-agent-changes-requested.md`

## Current Status

**All 16 blockers have been addressed.**

### Fixes Applied:

| Blocker | Fix |
|---------|-----|
| B1 (Zero tests) | Added 6 test files: config_test.go, keyring_test.go, oauth_test.go, auth_test.go, client_test.go, output_test.go, app_test.go |
| B2 (math/rand) | OAuth PKCE now uses `crypto/rand` for code verifier and state generation |
| B3 (Missing client_id) | Token exchange now includes `client_id` in form data for public clients |
| B4 (Redirect URI mismatch) | Both registration and callback now use consistent `http://localhost:{port}/callback` |
| B5 (Server never shut down) | `startCallbackServer()` returns server reference; `Shutdown()` called in defer |
| B6 (Weak key derivation) | Now uses PBKDF2 with 100k iterations + 32-byte salt instead of single SHA-256 |
| B7 (Dir permissions) | Credentials directory created with 0700 instead of 0755 |
| B8 (URL parsing) | `url.Parse` error now checked and returned |
| B9 (Retry POST) | Non-idempotent requests (POST) are not retried to avoid empty body issue |
| B10 (Token save error) | `SaveRefreshToken` error now returned instead of ignored |
| B11 (Clear dead code) | Fixed to properly collect and return errors in `Clear()` |
| B12 (Global flags) | Wired via `PersistentPreRunE` that sets formatter on each command |
| B13 (Help text) | Changed required args from `[arg]` to `<arg>` (e.g., `done <id>`) |
| B14 (SilenceUsage) | Added `rootCmd.SilenceUsage = true` |
| B15 (No context.Context) | All App methods now accept `context.Context` as first parameter |
| B16 (Error-swallowing) | `config.Load()` and `app.New()` now return `(*T, error)` tuples |

### Files Modified:
- `cmd/mynow/main.go`
- `internal/config/config.go`
- `internal/app/app.go`
- `internal/api/client.go`
- `internal/output/output.go`

### Files Created:
- `internal/auth/oauth.go` + `oauth_test.go`
- `internal/auth/keyring.go` + `keyring_test.go`
- `internal/auth/auth_test.go`
- `internal/config/config_test.go`
- `internal/api/client_test.go`
- `internal/output/output_test.go`
- `internal/app/app_test.go`

## Test Results

All tests pass:
```
ok  	github.com/mindyournow/myn-cli/internal/api	(cached)
ok  	github.com/mindyournow/myn-cli/internal/app	(cached)
ok  	github.com/mindyournow/myn-cli/internal/auth	(cached)
ok  	github.com/mindyournow/myn-cli/internal/config	(cached)
ok  	github.com/mindyournow/myn-cli/internal/output	(cached)
```

## Remaining Work

None - all blockers addressed. Ready for re-review.
