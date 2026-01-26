# C4 Context Diagram

This diagram shows the system context for the Go Strava Weekly application.

```mermaid
graph TB
    subgraph External["External Systems"]
        Strava["Strava API<br/>(Planned)<br/>Workout data provider"]
        CSV["CSV Files<br/>Local storage"]
    end
    
    subgraph System["Go Strava Weekly"]
        App["Go Strava Weekly<br/>CLI Application<br/>Processes and stores workout data"]
    end
    
    User["Athlete<br/>Cyclist tracking workouts"]
    
    User -->|"Runs weekly sync<br/>(CLI)"| App
    App -->|"Will fetch workouts<br/>(HTTPS/REST)"| Strava
    App -->|"Saves workout data<br/>(File I/O)"| CSV
    User -->|"Records workouts<br/>(Mobile/Web)"| Strava
    User -->|"Views/analyzes data<br/>(Spreadsheet apps)"| CSV
    
    classDef userStyle fill:#08427b,stroke:#052e56,color:#fff
    classDef systemStyle fill:#1168bd,stroke:#0b4884,color:#fff
    classDef externalStyle fill:#999,stroke:#666,color:#fff
    
    class User userStyle
    class App systemStyle
    class Strava,CSV externalStyle
```

## Notes

- **Strava API integration** is currently planned but not yet implemented
- The application currently focuses on processing and storing workout data to CSV files
- CSV storage allows for easy analysis using spreadsheet applications
