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
- Load non-sensitive configuration from `config.json`
- Override sensitive OAuth secrets and refresh tokens from environment variables

## Architecture

The codebase follows a clean architecture split:

```text
internal/
├── domain/          # Core business entities and rules
├── application/     # Use cases and contracts
├── infrastructure/  # Strava HTTP providers, auth, CSV, file repositories
└── config/          # Config loading, validation, and env secret overlay
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

The app uses a split configuration model:

- `config.json`: non-sensitive values
- `.env` or process environment variables: sensitive values only

Runtime lookup order:

1. `config.json` in the current working directory
2. environment variables for sensitive fields only

Sensitive env vars are mandatory. If they are empty, the application will fail validation at startup.

### `config.json`

Example:

```json
{
  "outputFilePath": "",
  "strava": {
    "clientId": "131590",
    "apiBaseUrl": "https://www.strava.com/api/v3",
    "oauthBaseUrl": "https://www.strava.com/oauth/token"
  },
  "googleOAuth": {
    "clientId": "your-google-client-id.apps.googleusercontent.com",
    "oauthBaseUrl": "https://oauth2.googleapis.com/token"
  },
  "email": {
    "from": "your-email@gmail.com",
    "to": "recipient@example.com",
    "subject": "Cycling Workout Report"
  },
  "excelTemplate": {
    "templatePath": "template.xlsx",
    "sheetName": "Registos",
    "startCell": "B8"
  }
}
```

### Environment Variables

| Variable | Required | Description |
|---|---|---|
| `STRAVA_CLIENT_SECRET` | yes | Strava client secret used to refresh access tokens |
| `STRAVA_REFRESH_TOKEN` | yes | Strava refresh token used to obtain access tokens |
| `GOOGLE_CLIENT_SECRET` | yes* | Google OAuth client secret used for email notifications |
| `GOOGLE_REFRESH_TOKEN` | yes* | Google OAuth refresh token used for Gmail access |

Notes:
- `OutputFilePath` can be set in `config.json` or overridden by CLI flag `--output-file`.
- `*` Google OAuth secrets are required when email delivery is used.
- Non-sensitive values such as client ids, API URLs, email addresses, and Excel template settings come from `config.json`.
- `.env` is loaded automatically at startup when present.

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
export STRAVA_CLIENT_SECRET="<strava-client-secret>"
export STRAVA_REFRESH_TOKEN="<strava-refresh-token>"
export GOOGLE_CLIENT_SECRET="<google-client-secret>"
export GOOGLE_REFRESH_TOKEN="<google-refresh-token>"

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
- `googleOAuth.clientId` and `googleOAuth.oauthBaseUrl` in `config.json`
- `email.from`, `email.to`, and `email.subject` in `config.json`
- `GOOGLE_CLIENT_SECRET` and `GOOGLE_REFRESH_TOKEN` in environment variables

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

### Build Instructions

Before building for any platform:

1. prepare `config.json` with the non-sensitive values you want to use
2. keep secrets out of `config.json`
3. provide secrets via `.env` or environment variables at runtime

#### Linux

Build on Linux:

```bash
go build -o cycling-ride-collector .
```

Run on Linux with secrets in the environment:

```bash
export STRAVA_CLIENT_SECRET="<strava-client-secret>"
export STRAVA_REFRESH_TOKEN="<strava-refresh-token>"
export GOOGLE_CLIENT_SECRET="<google-client-secret>"
export GOOGLE_REFRESH_TOKEN="<google-refresh-token>"

./cycling-ride-collector --start-date 01/01/2026 --end-date 01/07/2026
```

Cross-build a Windows executable from Linux:

```bash
GOOS=windows GOARCH=amd64 go build -o cycling-ride-collector.exe .
```

#### Windows

Build on Windows PowerShell:

```powershell
go build -o cycling-ride-collector.exe .
```

Run on Windows PowerShell with secrets in the environment:

```powershell
$env:STRAVA_CLIENT_SECRET="<strava-client-secret>"
$env:STRAVA_REFRESH_TOKEN="<strava-refresh-token>"
$env:GOOGLE_CLIENT_SECRET="<google-client-secret>"
$env:GOOGLE_REFRESH_TOKEN="<google-refresh-token>"

.\cycling-ride-collector.exe --start-date 01/01/2026 --end-date 01/07/2026
```

Cross-build a Linux binary from Windows PowerShell:

```powershell
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o cycling-ride-collector .
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
```

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
