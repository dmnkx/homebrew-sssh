[English](README.md) | [한국어](README.ko.md)

# 🔑 sssh

> Pick a `~/.ssh/config` Host alias in the terminal, then SSH — with list, add, edit, remove, and show.

`sssh` is a small CLI around OpenSSH. With no arguments it opens a TUI to choose a host. Connecting always runs the system `ssh` binary (`syscall.Exec`); `sssh` does not re-implement SSH.

```sh
sssh
sssh prod
sssh add prod --host 10.0.0.1 --user ubuntu
```

## What is sssh?

`sssh` lets you jump to hosts by the **Host aliases** already in `~/.ssh/config`, and manage those blocks from the command line.

Wildcard hosts (`Host *`, `Host *.internal`, …) are skipped for listing and connect. If one block has several aliases (`Host prod staging`), each alias is selectable on its own.

The command that actually runs is:

```text
ssh <alias>
```

Only the alias is passed. User, Port, IdentityFile, ProxyJump, and the rest come from **OpenSSH reading your config**. `sssh` does not rebuild those flags.

## Key Features

- **TUI picker** — Run `sssh` with no args, filter hosts, press Enter to connect
- **Direct connect** — `sssh <alias>` or `sssh connect <alias>`
- **Config CRUD** — `list`, `add`, `edit`, `rm`, `show` against `~/.ssh/config`
- **OpenSSH, not a fork** — Exec the real `ssh`; your existing config keeps working
- **Safe writes** — Backup to `<path>.bak` before overwrite; `0700` dirs, `0600` files
- **Homebrew + GitHub Releases** — Install without a local Go toolchain

## Install

**Homebrew**

```sh
brew tap dmnkx/sssh
brew install sssh
```

The tap uses [`Formula/sssh.rb`](Formula/sssh.rb) from this repo (`github.com/dmnkx/homebrew-sssh`) on `main`. If you tapped an older clone that has no formula:

```sh
brew untap dmnkx/sssh
brew tap dmnkx/sssh
brew install sssh
```

If you already tapped and a newer version is out:

```sh
brew update
brew install sssh
```

**From source** (Go 1.24+)

```sh
git clone https://github.com/dmnkx/homebrew-sssh.git
cd homebrew-sssh
go build -o sssh .
./sssh --help
```

Requires `ssh` on `PATH`.

## Quick Start

```sh
# Interactive host picker
sssh

# Connect by alias
sssh prod

# Preview the ssh command without connecting
sssh --print-cmd prod

# List selectable hosts
sssh list

# Add / edit / remove / dump a host block
sssh add prod --host 10.0.0.1 --user ubuntu --port 22 --key ~/.ssh/id_ed25519
sssh edit prod --user ubuntu
sssh show prod
sssh rm prod
```

Use another config file:

```sh
sssh --config /tmp/myconfig list
```

A missing config file is treated as empty (not an error). Writes still backup an existing file first.

## Commands

| Command | What it does |
|---------|----------------|
| `sssh` | TUI picker if any selectable Host exists; error if none |
| `sssh <alias>` | Connect immediately |
| `sssh connect <alias>` | Same as above |
| `sssh list` | One line per selectable Host |
| `sssh add <alias>` | Add a Host block (fails if the alias exists — use `edit`) |
| `sssh edit <alias>` | Change only the flags you pass; `--user ""` clears a field |
| `sssh rm <alias>` | Drop that alias; remove the whole block if it was the only name |
| `sssh show <alias>` | Print the Host in SSH config form |

### `sssh add` flags

| Flag | Required | SSH keyword |
|------|----------|-------------|
| `--host` | yes | HostName (IP or DNS) |
| `--user` | no | User |
| `--port` | no | Port |
| `--key` | no | IdentityFile |
| `--jump` | no | ProxyJump |

Aliases must not contain `*`, `?`, or `!`.

### Global flags

| Flag | Description |
|------|-------------|
| `--config <path>` | SSH config file (default `~/.ssh/config`) |
| `--print-cmd` | Print the `ssh` command instead of exec |

## TUI

| Key | Action |
|-----|--------|
| ↑ / ↓ | Move cursor |
| Type | Filter by alias, HostName, or User (case-insensitive) |
| Backspace / Delete | Delete one filter character |
| Enter | Connect to the current row (no-op if the filter matches nothing) |
| Esc, Ctrl+C | Quit without connecting |
| `q` | Quit only when the filter is empty; otherwise types `q` |

## Layout

```text
main.go                 entrypoint
internal/cli/           cobra (list, add, edit, rm, show, connect)
internal/sshcfg/        parse and save ~/.ssh/config
internal/connect/       find ssh, print command, Exec
internal/tui/           host picker
```

Tests live in `*_test.go` next to each package.

```sh
go test ./...
```

## CI / Releases

Workflow: [`.github/workflows/ci.yml`](.github/workflows/ci.yml)

1. Push to `main`/`master` and pull requests → `go test`, `go vet`, `go build`, GoReleaser config check
2. Push a `v*` tag → after tests, GoReleaser publishes GitHub Release binaries and updates `Formula/sssh.rb` on `main` with asset URLs and SHA256

## License

[MIT](LICENSE)
