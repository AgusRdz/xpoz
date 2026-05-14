# xpoz

[![CI](https://github.com/AgusRdz/xpoz/actions/workflows/ci.yml/badge.svg)](https://github.com/AgusRdz/xpoz/actions/workflows/ci.yml)
[![Release](https://github.com/AgusRdz/xpoz/actions/workflows/release.yml/badge.svg)](https://github.com/AgusRdz/xpoz/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Expose local ports to the internet using **your own domain** and **Cloudflare Tunnel** — no central server, no xpoz account, no third-party tunnel service.

```
xpoz 3000
✓ Tunnel active
  https://a3f9c.yourdomain.com
  forwarding to  localhost:3000
  Press Ctrl+C to stop
```

---

## How it works

xpoz sets up a single wildcard CNAME (`*.yourdomain.com → your-tunnel.cfargotunnel.com`) and a single wildcard ingress rule in cloudflared that points to a local reverse proxy. Every tunnel you open is a new route in that proxy — no cloudflared restart needed, SSL is handled by Cloudflare automatically.

```
Your browser
    │  HTTPS (SSL by Cloudflare)
    ▼
Cloudflare Edge  ──  *.yourdomain.com  ──  Cloudflare Tunnel
                                                    │
                                          xpoz proxy :2080
                                                    │
                                           localhost:3000
```

---

## Prerequisites

- A domain managed in **Cloudflare DNS**
- **cloudflared** installed and authenticated ([download](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/downloads/))
- A Cloudflare Tunnel created (`cloudflared tunnel create <name>`)
- A Cloudflare API token with **Zone → DNS → Edit** permission

---

## Installation

**Linux / macOS**

```sh
curl -sSL https://raw.githubusercontent.com/AgusRdz/xpoz/main/install.sh | sh
```

**Windows (PowerShell)**

```powershell
irm https://raw.githubusercontent.com/AgusRdz/xpoz/main/install.ps1 | iex
```

**From source**

```sh
git clone https://github.com/AgusRdz/xpoz.git
cd xpoz
make install   # builds and copies to ~/.local/bin
```

---

## Setup

Run once after installing:

```sh
xpoz setup
```

The wizard will ask for:

| Field | Where to find it |
|---|---|
| Domain | The domain you manage in Cloudflare (e.g. `yourdomain.com`) |
| API Token | [dash.cloudflare.com/profile/api-tokens](https://dash.cloudflare.com/profile/api-tokens) — needs **Zone → DNS → Edit** |
| Tunnel ID | Output of `cloudflared tunnel list` |

Setup will:
1. Validate your API token and tunnel ID
2. Create `*.yourdomain.com` wildcard CNAME in Cloudflare DNS
3. Add the wildcard ingress rule to your cloudflared config
4. Save `~/.xpoz/config.toml`

> **After setup**, restart cloudflared so it picks up the new ingress rule:
> ```sh
> # Linux (systemd)
> sudo systemctl restart cloudflared
> # macOS (launchd)
> sudo launchctl stop com.cloudflare.cloudflared && sudo launchctl start com.cloudflare.cloudflared
> ```

---

## Usage

### Ephemeral tunnel — random subdomain

```sh
xpoz 3000
# → https://a3f9c.yourdomain.com
```

A new random subdomain each time. Press `Ctrl+C` to stop.

### Named tunnel — persistent URL

```sh
xpoz 3000 --name my-app
# → https://my-app.yourdomain.com  (always the same)
```

The subdomain is saved. Next time you run `xpoz 3000 --name my-app` it reuses the same URL.

### Reset subdomain

```sh
xpoz 3000 --name my-app --reset
# → https://x7k2m.yourdomain.com  (new random subdomain)
```

### Background mode

```sh
xpoz 3000 --name my-app --bg
# Tunnel running in background (PID 12345)
# https://my-app.yourdomain.com
# Stop with: xpoz stop my-app
```

The tunnel keeps running after the terminal closes.

---

## Commands

### `xpoz list`

List all saved tunnels:

```
NAME       URL                              PORT   STATUS    LAST USED
my-app     https://my-app.yourdomain.com    3000   running   2m ago
api        https://a3f9c.yourdomain.com     8080   stopped   3h ago
```

### `xpoz start [name]`

Start a saved background tunnel, or all of them:

```sh
xpoz start my-app   # start one
xpoz start          # start all saved tunnels
```

### `xpoz stop [name]`

Stop a running background tunnel, or all of them:

```sh
xpoz stop my-app   # stop one
xpoz stop          # stop all
```

### `xpoz delete <name>`

Remove a saved tunnel from the database:

```sh
xpoz delete my-app
```

The URL is freed but the wildcard CNAME remains (it is shared by all tunnels).

### `xpoz doctor`

Run a health check on your setup:

```
xpoz doctor

  ✓  Config file is valid
     ~/.xpoz/config.toml

  ✓  cloudflared binary is installed
     /usr/local/bin/cloudflared

  ✓  cloudflared config exists
     ~/.cloudflared/config.yml

  ✓  Proxy port is available
     :2080 is free

  ✓  Store database is accessible
     ~/.xpoz/tunnels.db

  ✓  Wildcard CNAME exists in DNS
     *.yourdomain.com → abc123.cfargotunnel.com

✓ All checks passed — xpoz is ready.
```

### `xpoz service install / uninstall`

Register xpoz as a system service so saved tunnels start automatically on boot:

```sh
xpoz service install    # installs and starts
xpoz service uninstall  # stops and removes
```

Uses systemd on Linux, launchd on macOS, and Windows Service on Windows.

### `xpoz update`

Update xpoz to the latest release:

```sh
xpoz update
```

### `xpoz auto-update [on|off]`

Enable or disable the background update check (enabled by default):

```sh
xpoz auto-update off
xpoz auto-update on
```

### `xpoz config`

Show the current configuration:

```sh
xpoz config        # print config file path and contents
xpoz config edit   # open in $EDITOR
```

### `xpoz uninstall`

Remove xpoz, its config, and tunnel database:

```sh
xpoz uninstall              # removes everything
xpoz uninstall --keep-data  # keeps ~/.xpoz/tunnels.db
```

---

## Configuration

Config file: `~/.xpoz/config.toml` (created by `xpoz setup`, permissions `0600`)

```toml
[domain]
name = "yourdomain.com"

[cloudflare]
token   = "your-cloudflare-api-token"
zone_id = "your-zone-id"
tunnel  = "your-tunnel-id"

[proxy]
port = 2080        # internal proxy port (change if 2080 is in use)

[subdomains]
style  = "hex"     # hex | slug | word
length = 5         # hex length (e.g. 5 → "a3f9c")
```

**Subdomain styles:**

| Style | Example | Combinations |
|---|---|---|
| `hex` (default) | `a3f9c` | 1,048,576+ |
| `slug` | `bright-river` | ~4,800 |
| `word` | `storm` | ~98 |

---

## Development

Requirements: Docker (for the Go toolchain container)

```sh
git clone https://github.com/AgusRdz/xpoz.git
cd xpoz
make hooks   # install pre-commit lint hook
make test    # run all tests
make build   # build binary to bin/xpoz
make lint    # run golangci-lint
```

**All Makefile targets:**

```
make build          Build via Docker (recommended)
make local-build    Build without Docker (requires Go in PATH)
make test           Run all tests
make test-unit      Run unit tests only (fast)
make lint           Run golangci-lint
make fmt            Format all Go source files
make hooks          Install git pre-commit hook
make install        Build and install to ~/.local/bin
make snapshot       Test goreleaser pipeline locally (no publish)
make keygen         Generate Ed25519 signing key pair (one-time)
make deps           Download and tidy all Go modules
make clean          Remove build artifacts
```

### Releasing

```sh
make release-patch   # v1.0.0 → v1.0.1
make release-minor   # v1.0.0 → v1.1.0
make release-major   # v1.0.0 → v2.0.0
```

This generates a `CHANGELOG.md` via git-cliff, creates a commit and tag, and pushes — the release workflow publishes the binaries automatically.

---

## Data stored locally

| Path | Contents | Permissions |
|---|---|---|
| `~/.xpoz/config.toml` | Domain, API token, tunnel ID | `0600` |
| `~/.xpoz/tunnels.db` | Saved tunnel names and URLs | `0600` |
| `~/.xpoz/pids/<name>.pid` | PID of background tunnel process | — |
| `~/.xpoz/update_check` | Cached latest version tag | — |

On Windows, `~/.xpoz/` is `%APPDATA%\xpoz\`.

---

## License

MIT — see [LICENSE](LICENSE)
