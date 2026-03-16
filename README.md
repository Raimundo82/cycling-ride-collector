# cycling-ride-collector

Go CLI application that collects cycling rides from Strava and exports one consolidated workout per day.

## Overview

The app fetches rides for a date range, applies a daily selection strategy, and writes normalized output for training analysis.

Main capabilities:
- Fetch Strava activities by period
- Map Strava payloads to internal workout model
- Resolve multiple rides per day with policy-based selection (`longest` or `merge`)
- Persist daily results to CSV or Excel (`.xlsx`)
- Send workout reports via email with file attachments
- Handle OAuth token refresh from environment configuration (Strava and Google)

## Architecture

The codebase follows a clean architecture split:

```text
internal/
├── domain/          # Core business entities and rules
├── application/     # Use cases and contracts
├── infrastructure/  # Strava HTTP providers, auth, CSV, file repositories
└── config/          # Environment-based configuration
```

Architecture and auth diagrams (Mermaid) are in `docs/diagrams`:
- `docs/diagrams/c4-context.md`
- `docs/diagrams/c4-container.md`
- `docs/diagrams/domain-model.md`
- `docs/diagrams/auth/c4-component.md`
- `docs/diagrams/auth/class-diagram.md`
- `docs/diagrams/auth/sequence-diagram.md`

## Prerequisites

- Go `1.25.3+`
- A Strava API app (client id/secret) and refresh token
- A Google OAuth app (client id/secret) and refresh token (for email notifications)

## Installation

```bash
git clone https://github.com/Raimundo82/cycling-ride-collector.git
cd cycling-ride-collector
go mod download
go build -o cycling-ride-collector .
```

## Configuration

The app uses environment variables loaded by `internal/config/config.go`.

| Variable | Required | Description |
|---|---|---|
| `STRAVA_API_BASE_URL` | yes | Strava API base URL (example: `https://www.strava.com/api/v3`) |
| `STRAVA_OAUTH_BASE_URL` | yes | Strava OAuth base URL (example: `https://www.strava.com/oauth/token`) |
| `STRAVA_CLIENT_ID` | yes | Strava client id (used for refresh) |
| `STRAVA_CLIENT_SECRET` | yes | Strava client secret (used for refresh) |
| `STRAVA_REFRESH_TOKEN` | yes | Strava refresh token used to obtain access tokens |
| `GOOGLE_CLIENT_ID` | yes* | Google OAuth client id (used for email notifications) |
| `GOOGLE_CLIENT_SECRET` | yes* | Google OAuth client secret (used for email notifications) |
| `GOOGLE_OAUTH_TOKEN_URL` | yes* | Google OAuth token URL (example: `https://oauth2.googleapis.com/token`) |
| `EMAIL_FROM` | yes* | Sender email address for workout reports |
| `EMAIL_TO` | yes* | Recipient email address(es) for workout reports (comma-separated) |
| `EMAIL_SUBJECT` | yes* | Subject line for workout report emails |
| `EXCEL_TEMPLATE_PATH` | no** | Excel template path used by Excel exporter (example: `template.xlsx`) |
| `EXCEL_TEMPLATE_SHEETNAME` | no** | Template sheet name (example: `Registos`) |
| `EXCEL_TEMPLATE_STARTCELL` | no** | Start cell where workout rows are written (example: `B8`) |

Notes:
- `OutputFilePath` is defined by CLI flag `--output-file` (or auto-generated if omitted).
- `*` Google OAuth and email configuration are required when running with the `--cron` flag (for scheduled email notifications).
- `**` Excel template env vars are required when the Excel exporter is selected in code.

## Usage

### CLI Flags

| Flag | Required | Description |
|---|---|---|
| `--start-date` | yes | Start date in `MM/DD/YYYY` |
| `--end-date` | yes | End date in `MM/DD/YYYY` |
| `--daily-workout-policy` | no | `longest` (default) or `merge` |
| `--min-duration` | no | Minimum workout duration in minutes (floored to `30`) |
| `--output-file` | no | Output path or basename (see exporter behavior below) |
| `--cron` | no | Enable cron mode (schedules email notifications; requires Google OAuth and email config) |

### Example Run

```bash
export STRAVA_API_BASE_URL="https://www.strava.com/api/v3"
export STRAVA_OAUTH_BASE_URL="https://www.strava.com/oauth/token"
export STRAVA_CLIENT_ID="<strava-client-id>"
export STRAVA_CLIENT_SECRET="<strava-client-secret>"
export STRAVA_REFRESH_TOKEN="<strava-refresh-token>"
export GOOGLE_CLIENT_ID="<google-client-id>"
export GOOGLE_CLIENT_SECRET="<google-client-secret>"
export GOOGLE_OAUTH_TOKEN_URL="https://oauth2.googleapis.com/token"
export EMAIL_FROM="your-email@gmail.com"
export EMAIL_TO="recipient@example.com"
export EMAIL_SUBJECT="Cycling Workout Report"

./cycling-ride-collector \
  --start-date 01/01/2026 \
  --end-date 01/07/2026 \
  --daily-workout-policy merge \
  --min-duration 30 \
  --output-file workouts_summary_2026-01-01_to_2026-01-07
```

## Output Exporters

Two repository implementations exist:
- CSV: `internal/infrastructure/activity/repository/csv`
- Excel: `internal/infrastructure/activity/repository/excel`

Current behavior in this branch:
- Exporter selection is done in code in `app.go` (inside `NewApp`), not via CLI/env flag.
- The current default is Excel.

To switch exporter, edit the repository passed to `usecase.NewSaveWorkoutPeriod` in `app.go`:
- CSV: `activity_csv.NewCSVWorkoutPeriodSaver(cfg.OutputFilePath)`
- Excel: `activity_excel.NewExcelWorkoutPeriodSaverWithOptions(...)`

Current output-file behavior:
- When `--output-file` is provided: uses the exact path specified.
- When `--output-file` is omitted: auto-generates a filename with `.xlsx` suffix (e.g., `workouts_summary_2026-01-01_to_2026-01-07.xlsx`).
- When email notifications are enabled, the generated `.xlsx` file is automatically attached to the Gmail message after the report is saved.

Excel template notes:
- `EXCEL_TEMPLATE_PATH` must point to a real `.xlsx` file (legacy `.xls` is not supported by `excelize`).
- `EXCEL_TEMPLATE_SHEETNAME` and `EXCEL_TEMPLATE_STARTCELL` must match the target template.

## Email Notifications

When running with the `--cron` flag and the appropriate Google OAuth and email configuration set, the app will:

1. Generate a consolidated workout report for the specified date range
2. Save the report as an Excel file
3. Send the report via Gmail to the configured recipient(s)

Required configuration:
- Google OAuth credentials (`GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_OAUTH_TOKEN_URL`)
- Email settings (`EMAIL_FROM`, `EMAIL_TO`, `EMAIL_SUBJECT`)

## Daily Workout Policies

### `longest`
Selects the longest workout of the day above `--min-duration`.

### `merge`
Merges all workouts above `--min-duration`:
- Sum: distance, duration, elevation
- Weighted by duration: average power, average heart rate, cadence, normalized power
- Keep max: maximum heart rate

If no workout satisfies the threshold for a day, a rest-day row is emitted.

## Development

### Make Targets

```bash
make format
make check-format
make lint
make test
make precommit
make prepush
```

### Test Commands

```bash
# all tests
make test

# equivalent
go test ./... -v
```

If your environment has cache permission restrictions, run tests with a local cache path:

```bash
GOCACHE=$(pwd)/.gocache go test ./... -v
```

## Contributing

1. Install hooks: `lefthook install`
2. Implement changes
3. Run checks: `make precommit`
4. Open a pull request

## License

TBD
