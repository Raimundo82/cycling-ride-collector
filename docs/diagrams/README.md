# Architecture Diagrams

This directory contains Mermaid diagrams as code that document the architecture of the Go Strava Weekly application.

## Diagrams

### C4 Context Diagram
**File:** [`c4-context.md`](./c4-context.md)

Shows the system context - how the Go Strava Weekly application fits into the broader ecosystem, including:
- Athletes (users)
- Strava API (external system - planned)
- CSV Files (local storage)

### C4 Container Diagram
**File:** [`c4-container.md`](./c4-container.md)

Shows the high-level technical architecture - the main containers/components of the application:
- Main Application (CLI entry point)
- Domain Layer (core business entities)
- Application Layer (use cases and business logic)
- Infrastructure Layer (CSV-based storage)

### Domain Model Diagram
**File:** [`domain-model.md`](./domain-model.md)

Shows the domain model - the core business entities and their relationships:
- Workout entity with all metrics
- WorkoutType enumeration
- Application contracts (interfaces)
- Infrastructure implementations

### Auth C4 Component Diagram
**File:** [`auth/c4-component.md`](./auth/c4-component.md)

Shows auth subsystem components and integrations:
- TokenService orchestration
- OAuth provider/client flow
- File token repository
- Strava OAuth API interaction

### Auth Class Diagram
**File:** [`auth/class-diagram.md`](./auth/class-diagram.md)

Shows auth interfaces, implementations, models, and key dependencies.

### Auth Sequence Diagram
**File:** [`auth/sequence-diagram.md`](./auth/sequence-diagram.md)

Shows token resolution flow, including cache hit, refresh path, and persistence.

## Viewing the Diagrams

### GitHub Rendering

All Mermaid diagrams render automatically in GitHub's markdown viewer. Simply click on any of the `.md` files above to view the diagrams directly in your browser.

### VS Code Extension

If you're using Visual Studio Code:

1. Install the "Markdown Preview Mermaid Support" extension
2. Open any `.md` file
3. Use the markdown preview (`Ctrl+Shift+V` or `Cmd+Shift+V`)

### Other Markdown Viewers

Most modern markdown viewers support Mermaid diagrams natively, including:
- GitLab
- Bitbucket
- Notion
- Obsidian
- Many static site generators (Hugo, Jekyll, etc.)

### Mermaid Live Editor

For editing and previewing diagrams:

1. Visit [Mermaid Live Editor](https://mermaid.live/)
2. Copy and paste the Mermaid code from any diagram
3. Edit and preview in real-time

## Updating the Diagrams

To update a diagram:

1. Edit the corresponding `.md` file
2. Modify the Mermaid code within the triple backticks
3. Preview using GitHub or your markdown viewer
4. Commit the updated file to version control

## Resources

- [Mermaid Official Documentation](https://mermaid.js.org/)
- [C4 Model](https://c4model.com/)
- [Mermaid Live Editor](https://mermaid.live/)
- [Mermaid Syntax Guide](https://mermaid.js.org/intro/)
