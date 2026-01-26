# C4 Container Diagram

This diagram shows the high-level technical architecture of the Go Strava Weekly application.

```mermaid
graph TB
    User["Athlete<br/>Cyclist"]
    Strava["Strava API<br/>(Planned)"]
    CSV["CSV Files"]
    
    subgraph App["Go Strava Weekly Application"]
        Main["Main Application<br/>CLI Entry Point"]
        
        subgraph Layers["Architecture Layers"]
            Domain["Domain Layer<br/>Workout, WorkoutType<br/>Core business entities"]
            Application["Application Layer<br/>SaveWorkout, MergeWorkouts<br/>Use cases & contracts"]
            Infrastructure["Infrastructure Layer<br/>CSVWorkoutRepository<br/>Data persistence"]
        end
    end
    
    User -->|"Executes (CLI)"| Main
    Main -->|"Uses use cases"| Application
    Application -->|"Uses entities"| Domain
    Application -->|"Uses via contracts"| Infrastructure
    Infrastructure -->|"Depends on"| Domain
    Infrastructure -->|"Will fetch (HTTPS/REST)"| Strava
    Infrastructure -->|"Reads/Writes (File I/O)"| CSV
    
    classDef userStyle fill:#08427b,stroke:#052e56,color:#fff
    classDef mainStyle fill:#1168bd,stroke:#0b4884,color:#fff
    classDef layerStyle fill:#438dd5,stroke:#2e6295,color:#fff
    classDef externalStyle fill:#999,stroke:#666,color:#fff
    
    class User userStyle
    class Main mainStyle
    class Domain,Application,Infrastructure layerStyle
    class Strava,CSV externalStyle
```

## Notes

- **Clean Architecture**: The application follows Clean Architecture principles with clear separation of concerns
- **Domain Layer**: Core business entities that are independent of external frameworks
- **Application Layer**: Use cases orchestrate the business logic
- **Infrastructure Layer**: Currently implements CSV-based storage; Strava API integration is planned
