# C4 Context Diagram

This diagram shows the system context for the Go Strava Weekly application.

```mermaid
C4Context
    title System Context diagram for Go Strava Weekly

    Person(athlete, "Athlete", "A cyclist or endurance athlete who tracks their workouts")
    
    System(strava_weekly, "Go Strava Weekly", "Fetches workout data from Strava API and saves to CSV files for weekly analysis and tracking")
    
    System_Ext(strava, "Strava API", "Provides access to athlete's workout data including power, heart rate, distance, and duration metrics")
    
    SystemDb_Ext(csv_files, "CSV Files", "Stores workout data locally for analysis")

    Rel(athlete, strava_weekly, "Runs weekly sync", "CLI")
    Rel(strava_weekly, strava, "Fetches workouts", "HTTPS/REST (planned)")
    Rel(strava_weekly, csv_files, "Saves workout data", "File I/O")
    Rel(athlete, strava, "Records workouts", "Mobile/Web App")
    Rel(athlete, csv_files, "Views/analyzes data", "Spreadsheet apps")
```

## Notes

- **Strava API integration** is currently planned but not yet implemented
- The application currently focuses on processing and storing workout data to CSV files
- CSV storage allows for easy analysis using spreadsheet applications
