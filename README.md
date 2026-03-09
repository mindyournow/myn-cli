# mynow — Mind Your Now CLI & TUI

A fast, scriptable, Linux-native terminal client for [Mind Your Now](https://mindyournow.com).

## Features

- **Single binary** — one install gives you both CLI and interactive TUI
- **Scriptable** — every MYN operation available as a CLI command with JSON output
- **Linux-first** — works in terminals, SSH sessions, and remote environments
- **Package-friendly** — `.deb`, `.rpm`, `.apk`, Arch packages via GoReleaser

## Install

### From GitHub Releases

```bash
# Download latest release
curl -sL https://github.com/mindyournow/myn-cli/releases/latest/download/mynow_linux_amd64.tar.gz | tar xz
sudo mv mynow /usr/local/bin/
```

### From Source

```bash
git clone https://github.com/mindyournow/myn-cli.git
cd myn-cli
make build
./bin/mynow version
```

## Usage

```bash
# Authenticate
mynow login

# Manage inbox
mynow inbox add "Call Sam"
mynow inbox list

# Tasks
mynow task done abc123
mynow task snooze abc123

# Current focus
mynow now list
mynow now focus

# Daily review
mynow review daily

# Interactive TUI
mynow tui
# or just:
mynow

# JSON output for scripting
mynow now list --json | jq .
```

## Configuration

Config is stored in `~/.config/mynow/`. Credentials use the Linux Secret Service API (GNOME Keyring / KDE Wallet).

Set `MYN_API_URL` to point at a custom backend:

```bash
export MYN_API_URL=https://api.myn.localhost
```

## Plugins

```bash
mynow plugin list
mynow plugin enable openclaw
```

Plugins live in `~/.config/mynow/plugins/`. The core client is fully functional without any plugins.

## Development

```bash
make build        # Build binary
make test         # Unit tests
make lint         # Linting
make clean        # Clean build artifacts
```

### Integration Tests

Integration tests run against a real MYN backend via Docker Compose:

```bash
# Requires MYN backend source at ~/Projects/myn/api (or set MYN_BACKEND_PATH)
make test-integration
```

## License

[MIT](LICENSE)
