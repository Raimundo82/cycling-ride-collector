# Architecture Diagrams

This directory contains PlantUML diagrams as code that document the architecture of the Go Strava Weekly application.

## Diagrams

### C4 Context Diagram
**File:** `c4-context.puml`

Shows the system context - how the Go Strava Weekly application fits into the broader ecosystem, including:
- Athletes (users)
- Strava API (external system)
- Google Sheets (external system)

[View Online](https://www.plantuml.com/plantuml/uml/c4-context.puml)

### C4 Container Diagram
**File:** `c4-container.puml`

Shows the high-level technical architecture - the main containers/components of the application:
- Main Application (CLI entry point)
- Domain Layer (core business entities)
- Application Layer (use cases and business logic)
- Infrastructure Layer (external integrations)

[View Online](https://www.plantuml.com/plantuml/uml/c4-container.puml)

### Domain Model Diagram
**File:** `domain-model.puml`

Shows the domain model - the core business entities and their relationships:
- Workout entity with all metrics
- WorkoutType enumeration
- Application contracts (interfaces)
- Infrastructure implementations

[View Online](https://www.plantuml.com/plantuml/uml/domain-model.puml)

## Viewing the Diagrams

### Option 1: PlantUML Online Server

You can view any diagram by using the PlantUML online server. For each diagram file:

1. Copy the content of the `.puml` file
2. Go to [PlantUML Online Server](http://www.plantuml.com/plantuml/uml/)
3. Paste the content into the text area
4. The diagram will be automatically rendered

### Option 2: Direct URL Encoding

You can generate a direct link by encoding the PlantUML file:

```bash
# Example for c4-context.puml
cat c4-context.puml | plantuml -encodeurl
```

Then visit: `http://www.plantuml.com/plantuml/svg/[encoded-string]`

### Option 3: VS Code Extension

If you're using Visual Studio Code:

1. Install the "PlantUML" extension
2. Open any `.puml` file
3. Press `Alt+D` to preview the diagram

### Option 4: Local PlantUML Installation

Install PlantUML locally to generate images:

```bash
# Install PlantUML (requires Java)
# On macOS with Homebrew:
brew install plantuml

# Generate PNG images
plantuml c4-context.puml
plantuml c4-container.puml
plantuml domain-model.puml

# Generate SVG images (better quality)
plantuml -tsvg c4-context.puml
plantuml -tsvg c4-container.puml
plantuml -tsvg domain-model.puml
```

## Updating the Diagrams

To update a diagram:

1. Edit the corresponding `.puml` file
2. View the changes using one of the methods above
3. Commit the updated `.puml` file to version control

## Resources

- [PlantUML Official Site](https://plantuml.com/)
- [C4 Model](https://c4model.com/)
- [PlantUML C4 Extension](https://github.com/plantuml-stdlib/C4-PlantUML)
- [PlantUML Language Reference](https://plantuml.com/guide)
