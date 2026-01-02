# Architecture Overview

This document explains how the system is structured and how components interact.

## System Architecture

```mermaid
graph TB
    subgraph "Client"
        Browser[Web Browser]
    end
    
    subgraph "Frontend"
        App[App.jsx]
        TreeView[TreeView]
        SidePanel[SidePanel]
        API[API Client]
    end
    
    subgraph "Backend"
        Router[Router]
        Handlers[Handlers]
        Services[Services]
        Repository[Repository]
    end
    
    subgraph "Database"
        DB[(PostgreSQL)]
    end
    
    Browser --> App
    App --> TreeView
    App --> SidePanel
    App --> API
    API --> Router
    Router --> Handlers
    Handlers --> Services
    Services --> Repository
    Repository --> DB
```

## How Requests Flow

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant Backend
    participant Database
    
    User->>Frontend: User Action
    Frontend->>Backend: HTTP Request
    Backend->>Database: Query
    Database-->>Backend: Data
    Backend-->>Frontend: JSON Response
    Frontend-->>User: Updated UI
```

## Backend Layers

The backend uses a **layered architecture** for clear separation:

```mermaid
graph TD
    Request[HTTP Request] --> API[API Layer<br/>Handlers & Routes]
    API --> Service[Service Layer<br/>Business Logic]
    Service --> Repo[Repository Layer<br/>Data Access]
    Repo --> DB[(Database)]
```

**Layer Responsibilities:**
- **API Layer**: Handles HTTP requests, validates input, formats responses
- **Service Layer**: Contains business logic and rules
- **Repository Layer**: Manages database queries and data mapping

## Frontend Components

```mermaid
graph TD
    App[App.jsx<br/>Main State] --> TreeView[TreeView<br/>Tree Display]
    App --> SidePanel[SidePanel<br/>Details]
    App --> TreeManager[TreeManager<br/>Import/Export]
    
    TreeView --> Modals[Edit Modals]
    SidePanel --> NotesList[Notes List]
```

**Key Components:**
- **App.jsx**: Manages application state and coordinates components
- **TreeView**: Displays the interactive tree structure
- **SidePanel**: Shows details for selected prompts
- **TreeManager**: Handles import/export and save/load operations

## Database Structure

```mermaid
erDiagram
    prompts ||--o{ nodes : "has"
    prompts ||--o{ notes : "has"
    
    prompts {
        int id
        string title
        text description
    }
    
    nodes {
        int id
        int prompt_id
        string name
        text action
    }
    
    notes {
        int id
        int prompt_id
        text content
    }
```

**Relationships:**
- Each prompt can have multiple nodes (subprompts)
- Each prompt can have multiple notes (annotations)
- Deletions cascade (deleting a prompt deletes its nodes and notes)

## Security Flow

```mermaid
graph LR
    Request[Request] --> Check{Origin Check}
    Check -->|Frontend| Allow[Allow]
    Check -->|External| Auth{API Key?}
    Auth -->|Valid| Allow
    Auth -->|Invalid| Reject[Reject]
```

**Security:**
- Frontend requests (from browser) are automatically allowed
- External requests (curl, etc.) require API key
- API key passed in `Authorization: Bearer <key>` header

## Deployment

```mermaid
graph TB
    User[User] --> Frontend[Cloud Run<br/>Frontend]
    Frontend --> Backend[Cloud Run<br/>Backend]
    Backend --> DB[(Cloud SQL<br/>PostgreSQL)]
```

**Deployment:**
- Frontend and backend deployed separately on Cloud Run
- Database hosted on Cloud SQL
- All communication over HTTPS
