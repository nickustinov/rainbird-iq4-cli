---
name: setup
description: Initial setup wizard for iq4-cli. Guides user through authentication, scans controllers and stations, and creates LAWN.md with their irrigation schedule.
allowed-tools: Bash, Read, Edit, Write, AskUserQuestion
argument-hint: ""
---

# iq4-cli setup wizard

Guide the user through setting up their irrigation schedule management.

## Step 1 – build the CLI

Check if `iq4-cli` binary exists in the project root. If not, build it:

```bash
cd <project-root> && go build -o iq4-cli .
```

## Step 2 – authentication

Ask the user for their Rain Bird username and password, then run:

```bash
iq4-cli login <username> <password>
```

The CLI handles the full OIDC auth flow. The token is stored at `~/.iq4/token` (0600 permissions, password is not stored). Tokens expire after ~2 hours – re-run login when needed.

## Step 3 – scan controllers

Run `iq4-cli controllers` to get all controllers and their online status. Present them to the user in a table. Ask if any should be excluded from schedule management (like a greenhouse with its own schedule).

## Step 4 – scan stations

For each managed controller, run:
- `iq4-cli stations <controller-id>` for basic info
- Also fetch full station details from the API to check if landscape/sprinkler types are set (`areaLevel2Name`, `areaLevel3Name` fields from `GetStationListForSatellite`)

Present each controller's stations. For stations that have landscape/sprinkler info set in IQ4, note it. For stations missing this info, ask the user to describe:
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

## Step 7 – create LAWN.md

Using LAWN.md.example as the template, create LAWN.md with:
- Property details
- All controller and station info (including sprinkler types)
- Current schedule (programs, days, start times, runtimes)
- Weekly timeline showing all watering windows
- Design rationale based on sprinkler types and soil
- Seasonal review calendar based on location climate

## Step 8 – explain next steps

Tell the user:

> Your irrigation schedule is now documented in LAWN.md. Here's how to use the tools:
>
> - **`/irrigation`** – check current status, detect drift from LAWN.md, overlap analysis
> - **`/irrigation adjust`** – get seasonal adjustment suggestions for the current month
> - **`/irrigation apply`** – apply suggested changes via the CLI
>
> LAWN.md is in .gitignore (it contains your personal controller/station IDs).
> Review the schedule and run `/irrigation adjust` when the season changes.
