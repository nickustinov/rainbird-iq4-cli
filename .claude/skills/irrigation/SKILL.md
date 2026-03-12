---
name: irrigation
description: Check irrigation schedule status, detect overlaps, suggest and apply seasonal adjustments. Use when the user asks about watering, irrigation, schedule, or seasonal adjustments.
allowed-tools: Bash, Read, Edit
argument-hint: [status|adjust|apply]
---

# Irrigation schedule manager

You manage the irrigation schedule for a property in Corroios, Portugal with 5 Rain Bird controllers sharing a single water pump. **No two controllers can water simultaneously.**

## Setup

1. Read [LAWN.md](../../../LAWN.md) for the full zone inventory, current schedule, design rationale, and review calendar.
2. The CLI tool is `./iq4-cli` in the project root. Build it first if missing: `go build -o iq4-cli .`
3. Check LAWN.md for any excluded controllers (e.g. greenhouse managed separately).

## Commands

### `/irrigation` or `/irrigation status`

Full status check:

1. Read LAWN.md for expected state
2. Pull live data from all 4 managed controllers using the CLI:
   ```bash
   iq4-cli controllers
   iq4-cli programs
   iq4-cli runtimes <controller-id>    # for each controller
   iq4-cli start-times <controller-id> # for each controller
   ```
3. Compare live state to LAWN.md – flag any drift
4. Build the weekly timeline and check for overlaps (accounting for adjusted runtimes, not base)
5. Report:
   - Controller online/offline status
   - Any schedule drift from LAWN.md
   - Overlap analysis with exact time windows
   - Tightest gap between controllers

### `/irrigation adjust`

Suggest seasonal adjustments for the current month:

1. Read the review calendar in LAWN.md
2. Determine the current month and look up expected seasonal percentages
3. Pull current seasonal adjust values from all programs
4. Compare and suggest changes
5. If runtimes or watering frequency should change too, suggest those
6. Recalculate total adjusted runtimes and verify no overlaps with the new values
7. Present changes in a clear table, wait for user confirmation before applying

Use these guidelines for Corroios, Portugal:

| Month | ET level | Seasonal adjust range | Frequency |
|-------|----------|----------------------|-----------|
| Nov–Feb | Very low | 0–25% or off | 0–1x/week |
| Mar | Low | 30–45% | 2x/week |
| Apr | Moderate | 45–60% | 2x/week |
| May | Moderate–high | 60–80% | 2–3x/week |
| Jun–Aug | Peak | 80–100% | 3x/week, consider 2 runs/day for spray heads |
| Sep | Moderate–high | 60–80% | 2–3x/week |
| Oct | Moderate | 40–55% | 2x/week |

When increasing frequency beyond 2x/week, spread days evenly and re-check overlaps. If adding a second daily run (summer cycle-and-soak), the evening run must not overlap with any other controller.

### `/irrigation apply`

Apply previously suggested changes:

1. Apply each change using the CLI tool – never use raw curl/API calls
2. Verify each change by re-reading from the API
3. Update LAWN.md to reflect the new state
4. Show the updated weekly timeline

## CLI reference

```bash
# Read
iq4-cli programs                          # all programs, all controllers
iq4-cli programs <controller-id>          # programs for one controller
iq4-cli program <program-id>              # full detail with start times + runtimes
iq4-cli runtimes <controller-id>          # station runtimes
iq4-cli start-times <controller-id>       # start times grouped by program
iq4-cli controllers                       # all controllers with online status

# Modify
iq4-cli set-adjust <program-id> <percent> # seasonal adjust (e.g. 45)
iq4-cli set-days <program-id> <days>      # water days (e.g. MoTh, MoWeFr, 0100100)
iq4-cli set-runtime <step-id> <duration>  # base runtime (e.g. 10m, 30m)
iq4-cli add-start <program-id> <time>     # add start time (e.g. 04:30)
iq4-cli del-start <program-id> <id>       # delete start time
iq4-cli add-step <program-id> <station-id>  # assign station to program
iq4-cli del-step <step-id>                  # remove station from program
```

## Critical rules

1. **Never let controllers overlap.** Always calculate adjusted runtimes (base × seasonal %) and verify gaps between sequential controllers on shared days.
2. **All watering before sunrise.** Sunrise varies: ~06:50 March, ~06:15 June, ~07:30 December.
3. **Always use the CLI tool**, never raw API/curl calls.
4. **Update LAWN.md** after any changes so it stays the source of truth.
5. **Respect exclusions** – check LAWN.md for controllers marked as excluded.
6. **Base runtimes reflect sprinkler type** – don't change them unless the user asks. Seasonal adjust is the main lever.
7. **Sandy soil** – drains fast, better to water more frequently at lower volumes than infrequent deep soaks.
