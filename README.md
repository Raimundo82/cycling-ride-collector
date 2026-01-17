# Go Strava Weekly

A Go application that syncs Strava activities to Google Sheets using Clean Architecture principles.

## Architecture

This project follows Clean Architecture principles with clear separation of concerns:

```
.
├── main.go                          # Application entry point (dependency injection)
└── internal/
    ├── domain/                      # Domain Layer (innermost)
    │   ├── activity.go              # Activity entity
    │   └── sync_service.go          # Business logic interfaces
    ├── usecase/                     # Use Cases Layer
    │   └── sync_activities.go       # Sync activities business logic
    ├── infrastructure/              # Infrastructure Layer (adapters)
    │   ├── strava/
    │   │   └── client.go            # Strava API adapter
    │   └── sheets/
    │       └── client.go            # Google Sheets API adapter
    ├── interfaces/                  # Interface Adapters Layer
    │   └── cli/
    │       └── app.go               # CLI interface
    └── config/                      # Configuration
        └── config.go                # Config management
```

### Layers

1. **Domain Layer** (`internal/domain/`)
   - Contains business entities and interfaces
   - No dependencies on other layers
   - Defines contracts (interfaces) for external dependencies

2. **Use Case Layer** (`internal/usecase/`)
   - Contains application-specific business logic
   - Orchestrates data flow between repositories and presenters
   - Depends only on domain layer

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - Contains implementations of domain interfaces
   - Handles external dependencies (APIs, databases, etc.)
   - Adapts external services to domain interfaces

4. **Interface Adapters Layer** (`internal/interfaces/`)
   - Contains interfaces to the outside world (CLI, API, etc.)
   - Converts data from external format to use case format

5. **Configuration** (`internal/config/`)
   - Handles application configuration
   - Loads environment variables

## Configuration

The application uses environment variables for configuration:

- `STRAVA_API_KEY`: Your Strava API key
- `SHEETS_SPREADSHEET_ID`: Your Google Sheets spreadsheet ID

## Building

```bash
go build -o strava-weekly main.go
```

## Running

```bash
./strava-weekly
```

Or with configuration:

```bash
STRAVA_API_KEY=your_key SHEETS_SPREADSHEET_ID=your_sheet_id ./strava-weekly
```

## Testing

```bash
make test
```

## Linting

```bash
make lint
```

## Development

The project uses:
- Go 1.25.3
- golangci-lint for linting
- gofumpt and gci for formatting

### Make Commands

- `make format` - Format code
- `make check-format` - Check code formatting
- `make lint` - Run linter
- `make test` - Run tests
- `make precommit` - Run all checks before commit
- `make prepush` - Run all checks before push

## Clean Architecture Benefits

1. **Independence of Frameworks**: The business logic doesn't depend on external frameworks
2. **Testability**: Business logic can be tested without UI, database, or external services
3. **Independence of UI**: The UI can change without affecting business logic
4. **Independence of Database**: Business logic is independent of data storage
5. **Independence of External Agencies**: Business logic doesn't know about external services

## Next Steps

The current implementation provides placeholder adapters. To complete the application:

1. Implement real Strava API client in `internal/infrastructure/strava/client.go`
2. Implement real Google Sheets client in `internal/infrastructure/sheets/client.go`
3. Add authentication handling for both services
4. Add error handling and logging
5. Add unit tests for each layer
6. Add integration tests
