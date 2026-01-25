# go-strava-weekly

A Go application for syncing Strava workout data to Google Sheets on a weekly basis.

## Overview

This application fetches workout data from Strava and exports it to Google Sheets, helping you track and analyze your weekly cycling performance. The application merges multiple workouts from the same day and calculates weighted averages for key metrics.

## Features

- 📊 Fetch workout data from Strava API
- 📈 Merge multiple workouts from the same day with intelligent averaging
- ⚡ Support for power metrics (average, normalized)
- ❤️ Heart rate tracking (average, maximum)
- 🚴 Cadence and elevation tracking
- 🔄 Sync data to Google Sheets
- 🐳 Docker support for easy deployment
- ⚙️ Configurable minimum workout duration filter

## Workout Types

The application supports three types of workouts:
- **Estrada** - Road cycling
- **Rolo** - Indoor trainer/roller
- **Mixed** - Combined workouts

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

Detailed architecture diagrams are available as PlantUML diagrams in the [`docs/diagrams`](/docs/diagrams) directory:

- **[C4 Context Diagram](/docs/diagrams/c4-context.puml)** - System context and external interactions
- **[C4 Container Diagram](/docs/diagrams/c4-container.puml)** - High-level application architecture
- **[Domain Model Diagram](/docs/diagrams/domain-model.puml)** - Domain entities and relationships

See the [diagrams README](/docs/diagrams/README.md) for instructions on viewing these diagrams online.

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

### Example

```bash
export MINIMAL_WORKOUT_DURATION=45
./strava-weekly
```

## Usage

### Running Locally

```bash
./strava-weekly
```

### Running with Docker

```bash
docker run -e MINIMAL_WORKOUT_DURATION=30 strava-weekly
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

1. **Fetch Workouts**: The application fetches workout data from Strava for a specified date range (typically the last 7 days)

2. **Filter Workouts**: Only workouts meeting the minimum duration threshold are processed

3. **Merge Workouts**: If multiple workouts exist on the same day, they are intelligently merged:
   - Distance, duration, and elevation are summed
   - Power, heart rate, and cadence are weighted by duration
   - Maximum heart rate is preserved

4. **Save Data**: Merged workout data is saved to the configured destination (Google Sheets)

## Data Model

Each workout includes the following metrics:

- **ID**: Unique workout identifier
- **Workout Type**: Road (Estrada), Trainer (Rolo), or Mixed
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

## Project Status

🚧 **In Development** - The Strava API integration and Google Sheets export features are currently being implemented.

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
