# go-strava-weekly

A Go application for collecting Strava workout data and exporting daily summaries to CSV.

## Overview

This application fetches workout data from Strava and exports daily workout summaries to CSV, helping you track and analyze your cycling performance. It supports multiple strategies to select a "daily workout" when more than one ride exists on the same day.

## Features

- 📊 Fetch workout data from Strava API
- 📈 Merge multiple workouts from the same day with intelligent averaging
- 🧭 Configurable daily workout policy (`longest` or `merge`)
- ⚡ Support for power metrics (average, normalized)
- ❤️ Heart rate tracking (average, maximum)
- 🚴 Cadence and elevation tracking
- 🦵 Leg sensations mapping from Strava private notes
- 🗂️ Export to CSV
- 🐳 Docker support for easy deployment
- ⚙️ Configurable minimum workout duration filter

## Workout Types

The application supports the following workout types:
- **Estrada** - Road cycling
- **Rolo** - Indoor trainer/roller
- **Prova** - Race
- **Descanso** - Rest day

## Architecture

The project follows Clean Architecture principles with clear separation of concerns:

```
internal/
├── domain/          # Core business entities (Workout)
├── application/     # Use cases and business logic
│   ├── contracts/   # Interface definitions
│   └── usecase/     # Business logic implementation
├── infrastructure/  # External integrations (CSV, API clients)
└── config/          # Configuration management
```

### Architecture Diagrams

Detailed architecture diagrams are available as Mermaid diagrams in the [`docs/diagrams`](/docs/diagrams) directory:

- **[C4 Context Diagram](/docs/diagrams/c4-context.md)** - System context and external interactions
- **[C4 Container Diagram](/docs/diagrams/c4-container.md)** - High-level application architecture
- **[Domain Model Diagram](/docs/diagrams/domain-model.md)** - Domain entities and relationships

These diagrams render automatically in GitHub's markdown viewer. See the [diagrams README](/docs/diagrams/README.md) for more viewing options.

## Prerequisites

- Go 1.25.3 or higher
- Docker (optional, for containerized deployment)

## Installation

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/Raimundo82/go-strava-weekly.git
cd go-strava-weekly
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o strava-weekly
```

### Docker

Build the Docker image:
```bash
docker build -t strava-weekly .
```

## Configuration

The application can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `MINIMAL_WORKOUT_DURATION` | Minimum workout duration in minutes to include | `30` |
| `STRAVA_API_BASE_URL` | Strava API base URL | `https://www.strava.com/api/v3` |
| `STRAVA_ACCESS_TOKEN` | Current Strava access token | - |
| `OUTPUT_FILE_PATH` | Output CSV path | auto-generated when `--output-file` is omitted |
| `DAILY_WORKOUT_POLICY` | Daily workout selection policy (`longest` or `merge`) | `longest` |

### Strava Authentication

The application uses Bearer token authentication for Strava API requests. To set up:

1. Create a Strava API application at https://www.strava.com/settings/api
2. Obtain an access token through the OAuth2 flow
3. Set `STRAVA_ACCESS_TOKEN` environment variable

All API requests will include the access token in the `Authorization: Bearer <token>` header.

**Note**: Strava access tokens expire after 6 hours. You will need to handle token refresh manually or implement your own refresh logic as needed.

### Example

```bash
export MINIMAL_WORKOUT_DURATION=45
export DAILY_WORKOUT_POLICY=merge
./strava-weekly --start-date 01/01/2026 --end-date 01/07/2026
```

## Usage

### Running Locally

```bash
./strava-weekly \
  --start-date 01/01/2026 \
  --end-date 01/07/2026 \
  --access-token "$STRAVA_ACCESS_TOKEN" \
  --daily-workout-policy merge \
  --min-duration 30 \
  --output-file workouts_summary_2026-01-01_to_2026-01-07.csv
```

### Running with Docker

```bash
docker run --rm \
  -e STRAVA_API_BASE_URL=https://www.strava.com/api/v3 \
  -e STRAVA_ACCESS_TOKEN="$STRAVA_ACCESS_TOKEN" \
  strava-weekly \
  --start-date 01/01/2026 \
  --end-date 01/07/2026
```

## Development

### Available Make Commands

```bash
# Format code
make format

# Check code formatting (used in CI)
make check-format

# Run linter
make lint

# Run tests
make test

# Run all pre-commit checks
make precommit

# Run all pre-push checks
make prepush
```

### Code Quality Tools

The project uses several tools to maintain code quality:

- **gofumpt** - Stricter gofmt formatter
- **gci** - Go imports formatter
- **golangci-lint** - Comprehensive linter
- **goconvey** - Testing framework
- **lefthook** - Git hooks manager

### Running Tests

```bash
# Run all tests
go test ./... -v

# Or use make
make test
```

### Formatting Code

```bash
# Auto-format code
make format

# Check formatting without changes
make check-format
```

### Linting

```bash
make lint
```

## How It Works

1. **Fetch Workouts**: The application fetches workout data from Strava for each day in a user-provided date range.

2. **Filter Workouts**: Only workouts meeting the minimum duration threshold are processed

3. **Select Daily Workout**: If multiple workouts exist on the same day, the selected policy is applied:
   - **`longest`**: pick the longest workout above the minimum duration.
   - **`merge`**: merge all workouts above the minimum duration:
     - distance, duration, and elevation are summed
     - power, heart rate, and cadence are weighted by duration
     - maximum heart rate is preserved

4. **Save Data**: The resulting workout for each day is saved as a CSV row.

## Data Model

Each workout includes the following metrics:

- **ID**: Unique workout identifier
- **Workout Type**: Estrada, Rolo, Prova, or Descanso
- **Start Time**: When the workout began
- **Distance**: Total distance in kilometers
- **Duration**: Total duration in minutes
- **Elevation**: Total elevation gain in meters
- **Power Metrics**:
  - Average Power (watts)
  - Normalized Power (watts)
- **Heart Rate**:
  - Average Heart Rate (bpm)
  - Maximum Heart Rate (bpm)
- **Cadence**: Average cadence in RPM
- **Leg Sensations**: Optional perceived legs condition (from Strava private note mapping)

## Project Status

🚧 **In Development** - Core Strava ingestion and CSV export are implemented and evolving.

## Contributing

This project uses git hooks managed by lefthook. After cloning:

1. Install lefthook: `lefthook install`
2. Make your changes
3. Run checks: `make precommit`
4. Submit a pull request

## License

[Add your license information here]

## Author

Raimundo82
