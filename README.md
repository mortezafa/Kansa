# Kansa

A lightweight macOS daemon that tracks how much time you spend in Japanese
immersion tools — automatically, in the background.

## What it does

Kansa watches the frontmost application every 500ms and tracks accumulated
foreground time for a set of immersion-related programs:

- Anki
- mpv
- ttsu
- asbplayer
- VLC

Sessions are persisted to a local SQLite database (`kansa.db`), keyed by
program and date, so time spent accumulates across the day even if you
switch apps and come back.

## Usage

Build:

```bash
go build -o kansa .
```

Start the daemon:

```bash
./kansa
```

Runs in the background, writing logs to `kansa.log` and its PID to
`kansa.pid` in the working directory.

View accumulated time:

```bash
./kansa report
```

Stop the daemon:

```bash
kill $(cat kansa.pid)
```

## How it works

- Detects the active app via AppleScript (`osascript` + System Events).
- Tracks per-program state transitions (Running / Paused) in memory.
- On each pause, upserts that day's accumulated duration into SQLite.
- `report` reads directly from the database and prints totals per
  program per day.

## Requirements

- macOS (uses AppleScript to detect the frontmost app)
- Go 1.23+

## Status

Actively developed, single-user tool. Currently macOS-only.
