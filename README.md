# iq4-cli

Command-line tool for the Rain Bird IQ4 cloud API. Manage irrigation schedules, programs, and controllers from the terminal.

Built for use with LLMs – all output is JSON, all commands are composable. Includes [Claude Code](https://github.com/anthropics/claude-code) skills for automated schedule management.

## Claude Code skills

This project includes skills that let Claude manage your irrigation schedule through natural language:

### `/setup`

Interactive setup wizard. Guides you through authentication, scans your controllers and stations, detects sprinkler types and landscape settings, and creates a `LAWN.md` file documenting your full irrigation schedule.

### `/irrigation`

Full status check – pulls live data from all controllers, compares against your documented schedule, checks for overlaps, and flags any drift.

### `/irrigation adjust`

Suggests seasonal adjustment changes based on your location's climate and the current month. Accounts for sprinkler types, soil, and pump constraints.

### `/irrigation apply`

Applies the suggested changes via the CLI and updates LAWN.md.

## Why

Rain Bird's 2.0 app moved schedule management to the IQ4 cloud. The local controller API no longer returns schedule data on firmware 4.98+. This tool talks to the IQ4 cloud API directly, giving you full read/write access to your irrigation programs.

## Install

Requires Go 1.21+.

```bash
git clone https://github.com/nickustinov/rainbird-iq4-cli.git
cd rainbird-iq4-cli
go build -o iq4-cli .
```

## Quick start

```bash
# Log in with your Rain Bird username and password
./iq4-cli login <username> <password>

# List all controllers
./iq4-cli controllers

# Show programs for a controller
./iq4-cli programs <controller-id>

# Full program detail (with start times and runtimes)
./iq4-cli program <program-id>
```

The CLI handles the full OIDC auth flow automatically. Your JWT token is stored at `~/.iq4/token` with 0600 permissions. Tokens expire after ~2 hours – just re-run `login` when needed.

## Commands

### Read

```bash
iq4-cli sites                             # list all sites
iq4-cli controllers                       # list all controllers with connection status
iq4-cli stations <controller-id>          # list stations (with sprinkler/landscape type)
iq4-cli programs                          # list all programs across all controllers
iq4-cli programs <controller-id>          # list programs for a specific controller
iq4-cli program <program-id>              # full program detail
iq4-cli start-times                       # list all start times
iq4-cli start-times <controller-id>       # list start times for a controller
iq4-cli runtimes <controller-id>          # list station runtimes
```

### Write

```bash
iq4-cli set-adjust <program-id> <percent>            # set seasonal adjust (e.g. 45)
iq4-cli set-days <program-id> <days>                 # set water days (e.g. MoWeFr, 0101010)
iq4-cli set-runtime <step-id> <duration>             # set base runtime (e.g. 10m, 1h30m)
iq4-cli add-start <program-id> <time>                # add start time (e.g. 06:00)
iq4-cli del-start <program-id> <start-time-id>       # delete a start time
iq4-cli add-step <program-id> <station-id>           # assign station to a program
iq4-cli del-step <step-id>                           # remove station from a program
```

### Auth

```bash
iq4-cli login <username> <password>       # authenticate and store token
iq4-cli logout                            # clear stored token
```

## Output format

All commands output JSON to stdout. Status messages and errors go to stderr. This makes the CLI easy to compose with `jq`, scripts, or LLMs:

```bash
# Get all program names
./iq4-cli programs | jq '.[].name'

# Check if any controller is offline
./iq4-cli controllers | jq '.[] | select(.isConnected == false) | .name'
```

## Data model

```
Company → Sites → Controllers (Satellites)
                    ├── Stations (physical valve zones)
                    └── Programs (irrigation schedules)
                        ├── Start times (when to run)
                        ├── Program steps (station → runtime)
                        └── Seasonal adjust (% scaling)
```

## API documentation

See [docs/IQ4-API.md](docs/IQ4-API.md) for the full IQ4 cloud API reference, including all verified endpoints and write operation formats.

## Design

- **Go stdlib only** – zero external dependencies
- **JSON everywhere** – designed for LLM and script consumption
- **Flat package** – all code in root, easy to read and modify
- **Secure token storage** – 0600 permissions, never logged

## License

MIT
