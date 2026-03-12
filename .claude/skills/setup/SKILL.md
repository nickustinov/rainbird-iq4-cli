---
name: setup
description: Initial setup wizard for iq4-cli. Guides user through authentication, scans controllers and stations, and creates LAWN.md with their irrigation schedule.
allowed-tools: Bash, Read, Edit, Write, AskUserQuestion
argument-hint: ""
---

# iq4-cli setup wizard

Guide the user through setting up their irrigation schedule management.

**IMPORTANT: This is a read-only setup process. Do NOT modify any controller settings, programs, runtimes, start times, or station assignments. Only read data and create LAWN.md. Any changes to the irrigation schedule must be explicitly requested and confirmed by the user in a separate step.**

## Step 1 – build the CLI

Check if `iq4-cli` binary exists in the project root. If not, build it:

```bash
cd <project-root> && go build -o iq4-cli .
```

## Step 2 – authentication

Tell the user to run the login command themselves in a separate terminal. **Never ask for their password.** Display this:

> Open another terminal and run:
>
> ```bash
> ./iq4-cli login YOUR_USERNAME YOUR_PASSWORD
> ```
>
> This authenticates with Rain Bird and stores a token locally at `~/.iq4/token`. Your password is not stored. The token expires after ~2 hours – re-run login when needed.
>
> Let me know when you're logged in.

Wait for the user to confirm before proceeding. Then verify the token works by running `./iq4-cli sites` – if it fails, tell the user to try logging in again.

## Step 3 – scan controllers

Run `iq4-cli controllers` to get all controllers and their online status. Present them to the user in a table. Ask if any should be excluded from schedule management (like a greenhouse with its own schedule).

## Step 4 – scan stations

For each managed controller, run `iq4-cli stations <controller-id>`. The output includes `areaLevel2Name` (landscape type) and `areaLevel3Name` (sprinkler type) if configured in the Rain Bird app.

Present each controller's stations in a table. For stations that have landscape/sprinkler info, note it. For stations missing this info, ask the user to describe:
- What it waters (lawn, garden beds, ornamental, trees, etc.)
- Sprinkler type (spray heads, rotors, drip lines, bubblers, etc.)
- Any special conditions (shade, slope, swamp risk, etc.)

## Step 5 – gather property info

Ask the user for:
- **Location** (city, country – needed for sunrise times and ET estimates)
- **Soil type** (sandy, clay, loam, mixed)
- **Single pump constraint?** (can multiple controllers run simultaneously or not?)

## Step 6 – scan current schedule

For each managed controller, pull:
- `iq4-cli programs <controller-id>`
- `iq4-cli runtimes <controller-id>`
- `iq4-cli start-times <controller-id>`

Present the current schedule to the user. Note any issues:
- Empty programs (stations assigned but no start times, or vice versa)
- Ghost steps (0 runtime assignments)
- Potential overlaps
- Missing station assignments

## Step 7 – station-specific notes

Present a summary of all stations across all controllers and ask:

> Do any zones have special conditions I should know about? For example:
> - Shade (partial, heavy) – affects water needs and swamp risk
> - Slope or drainage issues
> - Sensitive plants or new plantings
> - Areas that tend to flood or stay soggy
> - Anything else that affects how much or how often to water
>
> Just tell me which stations and what the issue is, or say "none" to skip.

Record any notes the user provides – these go in a Notes column on the station tables in LAWN.md.

## Step 8 – create LAWN.md

Using LAWN.md.example as the template, create LAWN.md with:
- Property details
- All controller and station info (including sprinkler types and user notes)
- Current schedule (programs, days, start times, runtimes)
- Weekly timeline showing all watering windows
- Design rationale based on sprinkler types and soil
- Seasonal review calendar based on location climate

## Step 9 – explain next steps

Tell the user:

> Your irrigation schedule is now documented in LAWN.md. Here's how to use the tools:
>
> - **`/irrigation`** – check current status, detect drift from LAWN.md, overlap analysis
> - **`/irrigation adjust`** – get seasonal adjustment suggestions for the current month
> - **`/irrigation apply`** – apply suggested changes via the CLI
>
> LAWN.md is in .gitignore (it contains your personal controller/station IDs).
> Review the schedule and run `/irrigation adjust` when the season changes.
