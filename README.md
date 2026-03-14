# cycling-ride-collector

Go CLI application that collects cycling rides from Strava and exports one consolidated workout per day.

## Overview

The app fetches rides for a date range, applies a daily selection strategy, and writes normalized output for training analysis.

Main capabilities:
- Fetch Strava activities by period
- Map Strava payloads to internal workout model
- Resolve multiple rides per day with policy-based selection (`longest` or `merge`)
- Persist daily results to CSV or Excel (`.xlsx`)
- Handle OAuth token refresh and local token persistence
- Send the generated Excel report by email through the Gmail API

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
- A Strava API app (client id/secret)
- A token JSON file (`TOKEN_FILE_PATH`) with at least a Strava refresh token
- For email sending: a Google Cloud OAuth app, Gmail API enabled, and a Google refresh token obtained outside the app

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
| `STRAVA_OAUTH_BASE_URL` | yes | Strava OAuth base URL (example: `https://www.strava.com/oauth`) |
| `STRAVA_CLIENT_ID` | yes | Strava client id (used for refresh) |
| `STRAVA_CLIENT_SECRET` | yes | Strava client secret (used for refresh) |
| `TOKEN_FILE_PATH` | yes | Path to local token JSON file |
| `GOOGLE_CLIENT_ID` | email only | Google OAuth client id used to exchange the stored refresh token |
| `GOOGLE_CLIENT_SECRET` | email only | Google OAuth client secret used to exchange the stored refresh token |
| `GOOGLE_OAUTH_TOKEN_URL` | no | Google OAuth token URL. Defaults to `https://oauth2.googleapis.com/token` |
| `EMAIL_FROM` | email only | Sender email address used in the Gmail message |
| `EMAIL_TO` | email only | Comma-separated recipient list |
| `EMAIL_SUBJECT` | no | Email subject. Defaults to `Cycling Workout Report` |
| `EXCEL_TEMPLATE_PATH` | no* | Excel template path used by Excel exporter (example: `template.xlsx`) |
| `EXCEL_TEMPLATE_SHEETNAME` | no* | Template sheet name (example: `Registos`) |
| `EXCEL_TEMPLATE_STARTCELL` | no* | Start cell where workout rows are written (example: `B8`) |

Notes:
- `OutputFilePath` is defined by CLI flag `--output-file` (or auto-generated if omitted).
- Minimal workout duration and policy are CLI parameters, not environment variables.
- `*` Excel template env vars are required when the Excel exporter is selected in code.

### Token File Format

`TOKEN_FILE_PATH` supports both the legacy Strava-only shape and the canonical combined shape below:

```json
{
  "strava": {
    "access_token": "...",
    "refresh_token": "...",
    "expires_at": "2030-01-01T00:00:00Z"
  },
  "google": {
    "refresh_token": "...",
    "access_token": "",
    "expires_at": "0001-01-01T00:00:00Z"
  }
}
```

One-time Google setup happens outside the app:

1. Create a Google Cloud project and enable the Gmail API.
2. Create OAuth2 credentials of type Desktop app and record `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET`.
3. Use `google_client.http` to complete consent once, exchange the authorization code, and capture a refresh token.
4. Store that refresh token under the `google` section in `tokens.json`.

The app never opens a browser or runs an interactive Google consent flow. It only exchanges the stored refresh token for short-lived access tokens and keeps those access tokens in memory.

### Google HTTP Client Helper

The repository includes `google_client.http` to help with manual Gmail OAuth testing:
- Step 1 opens the consent URL.
- Step 2 exchanges the authorization code for access and refresh tokens.
- Step 3 refreshes an access token from the stored refresh token.
- Step 4 sends a test email through Gmail.

The send-email example in `google_client.http` uses the Gmail media upload endpoint with a readable MIME message body:

```http
POST https://gmail.googleapis.com/upload/gmail/v1/users/me/messages/send?uploadType=media
Authorization: Bearer {{accessToken}}
Content-Type: message/rfc822

From: your-address@gmail.com
To: your-address@gmail.com
Subject: Test email from google_client.http
MIME-Version: 1.0
Content-Type: text/plain; charset=utf-8

This is a test email sent from the HTTP client file.
```

That keeps the request editable in plain text instead of requiring a base64-encoded `raw` payload in the HTTP client file.

## Usage

### CLI Flags

| Flag | Required | Description |
|---|---|---|
| `--start-date` | yes | Start date in `MM/DD/YYYY` |
| `--end-date` | yes | End date in `MM/DD/YYYY` |
| `--daily-workout-policy` | no | `longest` (default) or `merge` |
| `--min-duration` | no | Minimum workout duration in minutes (floored to `30`) |
| `--output-file` | no | Output path or basename (see exporter behavior below) |

### Example Run

```bash
export STRAVA_API_BASE_URL="https://www.strava.com/api/v3"
export STRAVA_OAUTH_BASE_URL="https://www.strava.com/oauth"
export STRAVA_CLIENT_ID="<client-id>"
export STRAVA_CLIENT_SECRET="<client-secret>"
export TOKEN_FILE_PATH="/absolute/path/to/token.json"
export GOOGLE_CLIENT_ID="<google-client-id>"
export GOOGLE_CLIENT_SECRET="<google-client-secret>"
export EMAIL_FROM="me@gmail.com"
export EMAIL_TO="coach@example.com,me@gmail.com"

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

Current output-file behavior by exporter:
- CSV exporter: uses `--output-file` exactly as provided.
- Excel exporter: appends `.xlsx` to `cfg.OutputFilePath` in `app.go`.

Example:
- `--output-file workouts_summary_2026-01-01_to_2026-01-07` produces `workouts_summary_2026-01-01_to_2026-01-07.xlsx` with Excel exporter.
- When email sending is configured, that generated `.xlsx` file is attached to the Gmail message after the report is saved.

Excel template notes:
- `EXCEL_TEMPLATE_PATH` must point to a real `.xlsx` file (legacy `.xls` is not supported by `excelize`).
- `EXCEL_TEMPLATE_SHEETNAME` and `EXCEL_TEMPLATE_STARTCELL` must match the target template.

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
