# C4 Container Diagram

This diagram shows the high-level technical architecture of the Go Strava Weekly application.

```mermaid
C4Container
    title Container diagram for Go Strava Weekly

    Person(athlete, "Athlete", "A cyclist who uses the system for workout tracking")
    
    System_Ext(strava, "Strava API", "External workout data provider (planned)")
    
    System_Boundary(strava_weekly, "Go Strava Weekly") {
        Container(main, "Main Application", "Go", "Command-line application that orchestrates the workout sync process")
        
        Container(domain, "Domain Layer", "Go Package", "Contains core business entities (Workout, WorkoutType) and domain logic")
        
        Container(application, "Application Layer", "Go Package", "Contains use cases (SaveWorkout, MergeWorkouts) and contracts (WorkoutProvider, WorkoutRepository)")
        
        Container(infrastructure, "Infrastructure Layer", "Go Package", "Contains CSVWorkoutRepository for data persistence")
    }
    
    ContainerDb_Ext(csv, "CSV Files", "File System", "Stores workout data locally")

    Rel(athlete, main, "Executes", "CLI")
    Rel(main, application, "Uses use cases", "Function calls")
    Rel(application, domain, "Uses entities", "Import")
    Rel(application, infrastructure, "Uses via contracts", "Interface")
    Rel(infrastructure, domain, "Depends on", "Import")
    Rel(infrastructure, strava, "Will fetch workouts", "HTTPS/REST (planned)")
    Rel(infrastructure, csv, "Reads/Writes", "File I/O")
```

## Notes

- **Clean Architecture**: The application follows Clean Architecture principles with clear separation of concerns
- **Domain Layer**: Core business entities that are independent of external frameworks
- **Application Layer**: Use cases orchestrate the business logic
- **Infrastructure Layer**: Currently implements CSV-based storage; Strava API integration is planned
