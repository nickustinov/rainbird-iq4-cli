# iq4-cli – LLM build guide

Rain Bird IQ4 command-line tool. Go, stdlib only, JSON output.

## Required reading

- **[docs/IQ4-API.md](docs/IQ4-API.md)** – IQ4 cloud API reference: auth flow, endpoints, data models.
- **[LAWN.md](LAWN.md)** – user's irrigation schedule (gitignored, created by `/setup`).
- **[LAWN.md.example](LAWN.md.example)** – template for LAWN.md.

## Data model

```
Company → Sites → Controllers (Satellites)
                    ├── Stations (physical valve zones, 4 or 8 per controller)
                    └── Programs (A, B, C – irrigation schedules)
                        ├── Start times (when to run, multiple per program)
                        ├── Program steps (station → runtime mapping)
                        └── Seasonal adjust (% scaling of runtimes)
```

## Project rules

1. **Stdlib only** – no external dependencies.
2. **JSON output** – all commands write JSON to stdout, status/errors to stderr.
3. **Flat package** – all `.go` files in the root package.
4. **Token storage** – JWT stored at `~/.iq4/token` with 0600 permissions.
5. **No sensitive data in code** – no usernames, IPs, MACs, or personal IDs in committed files.

## Skills

- **`/setup`** – initial setup wizard: auth, scan controllers, create LAWN.md
- **`/irrigation`** – status check, overlap analysis, drift detection
- **`/irrigation adjust`** – suggest seasonal adjustments for current month
- **`/irrigation apply`** – apply suggested changes

## Style

1. **European-style titles** – never use American Title Case.
2. **En dashes, not em dashes** – use – not —.

## Git

1. **No autonomous commits** – never commit without explicit user instruction.
2. **No co-authorship** – never add Co-Authored-By lines.
