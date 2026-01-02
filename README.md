# Prompt Tree Management System

A full-stack application for visualizing, editing, and managing hierarchical prompt structures. Built with React and Go, deployed on Google Cloud Run.

## Overview

This system allows users to:
- Visualize prompt trees in an interactive interface
- Create, edit, and delete prompts and their subprompts (nodes)
- Add annotations (notes) to prompts
- Import/export prompt trees as JSON
- Save and load multiple named tree configurations

## Live Application

### Frontend
**URL**: `https://frontend-709459926380.us-central1.run.app`

The frontend provides a web interface for managing prompt trees. You can:
- View the interactive tree structure
- Click on prompts to see details in the side panel
- Right-click on nodes to edit or delete them
- Add notes to prompts
- Import/export trees as JSON
- Save and load multiple tree configurations

<img width="1437" height="800" alt="image" src="https://github.com/user-attachments/assets/8e9171e7-6227-4552-9fba-6e95196b47c4" />

Importing Trees

<img width="1437" height="800" alt="image" src="https://github.com/user-attachments/assets/56479847-2824-4230-bf1a-e9ac1f38747f" />

Exporting Trees

<img width="1437" height="800" alt="image" src="https://github.com/user-attachments/assets/8959c1bf-0617-4ccb-9c53-a5801640f082" />

Save/Load Trees

<img width="1437" height="800" alt="image" src="https://github.com/user-attachments/assets/7203a061-adc6-4df1-8a6e-d3b0b6d157d9" />

Editing Tree Subprompts

<img width="1437" height="800" alt="image" src="https://github.com/user-attachments/assets/d0ff117f-0523-4028-8965-616106c2293d" />

Editing Tree Prompts

<img width="1437" height="800" alt="image" src="https://github.com/user-attachments/assets/1c6898b2-b737-4eec-8a38-6afe769733f8" />

Side Panel

<img width="360" height="636" alt="image" src="https://github.com/user-attachments/assets/e909bb97-a0fa-4b74-bad9-3ef0d2d2723a" />

Edit Actions

<img width="1440" height="810" alt="image" src="https://github.com/user-attachments/assets/8c485920-2d8e-444d-ba0e-bde15fb4de68" />

Add Notes

<img width="361" height="185" alt="image" src="https://github.com/user-attachments/assets/baf1999e-f348-4a9c-b0f6-46eb5e66ff0e" />







### Backend API
**URL**: `https://backend-709459926380.us-central1.run.app`

The backend provides a RESTful API for all operations. You can test all routes without any local setup - just use the API key provided.

**Interactive API Documentation**: [https://backend-709459926380.us-central1.run.app/docs](https://backend-709459926380.us-central1.run.app/docs)

## Testing the API

You can test all API routes directly using curl commands. No local setup required!

### Quick Test

```bash
# Set your API key and backend URL
export API_KEY="<YOUR_API_KEY>"
export BACKEND_URL="https://backend-709459926380.us-central1.run.app"

# Test health endpoint (no auth required)
curl $BACKEND_URL/health

# Test getting the tree (requires API key)
curl -H "Authorization: Bearer $API_KEY" $BACKEND_URL/tree
```

### Complete Testing Guide

See [API_ROUTES.md](docs/API_ROUTES.md) for:
- Complete list of all API endpoints
- Sample curl commands for each route
- Expected responses
- Error handling examples

All routes can be tested against the deployed backend: no local setup needed!

## Technology Stack

### Backend
- **Language**: Go 1.23
- **Framework**: 
  - Chi (HTTP router)
  - Huma (API framework with OpenAPI documentation)
- **Database**: PostgreSQL (Cloud SQL)
- **Deployment**: Google Cloud Run
- **Container**: Docker

### Frontend
- **Framework**: React 18
- **Build Tool**: Vite
- **HTTP Client**: Axios
- **Deployment**: Google Cloud Run
- **Web Server**: Nginx (production)

## API Endpoints

### Base URL
`https://backend-709459926380.us-central1.run.app`

### Authentication
All endpoints (except `/health` and `/docs`) require API key authentication:

```bash
curl -H "Authorization: Bearer <YOUR_API_KEY>" $BACKEND_URL/tree
```

### Available Endpoints

**Tree Management:**
- `GET /tree` - Get full prompt tree
- `GET /tree/export` - Export tree as JSON
- `POST /tree/import` - Import tree from JSON
- `POST /tree/save` - Save current tree
- `GET /tree/saves` - List saved trees
- `POST /tree/load/{name}` - Load saved tree
- `DELETE /tree/saves/{name}` - Delete saved tree

**Prompts:**
- `GET /prompts/{id}` - Get single prompt
- `POST /prompts/{id}` - Create prompt
- `PUT /prompts/{id}` - Update prompt
- `DELETE /prompts/{id}` - Delete prompt

**Nodes:**
- `GET /prompts/{id}/nodes` - Get nodes for a prompt
- `POST /prompts/{id}/nodes` - Create node
- `PUT /prompts/{id}/nodes/{nodeId}` - Update node
- `DELETE /prompts/{id}/nodes/{nodeId}` - Delete node

**Notes:**
- `GET /prompts/{id}/notes` - Get notes for a prompt
- `POST /prompts/{id}/notes` - Create note
- `PUT /prompts/{id}/notes/{noteId}` - Update note
- `DELETE /prompts/{id}/notes/{noteId}` - Delete note

See [API_ROUTES.md](docs/API_ROUTES.md) for detailed examples and sample responses.

## Deployment

The application is deployed on Google Cloud Platform:

- **Frontend**: Cloud Run (Nginx serving React build)
- **Backend**: Cloud Run (Go application)
- **Database**: Cloud SQL (PostgreSQL)

### Deployment Process

**Backend:**
```bash
cd backend
gcloud builds submit --config cloudbuild.yaml
```

**Frontend:**
```bash
cd frontend
gcloud builds submit --config cloudbuild.yaml
```

### Environment Variables

**Backend (Cloud Run):**
- `INSTANCE_CONNECTION_NAME` - Cloud SQL instance
- `DB_USER`, `DB_PASS`, `DB_NAME` - Database credentials
- `PORT` - Server port (default: 8080)
- `API_KEY` - API key for authentication
- `ENVIRONMENT` - Environment name (production)

**Frontend (Cloud Run):**
- `VITE_API_URL` - Backend API URL (set during build)

## Architecture

For detailed architecture documentation, see [ARCHITECTURE.md](docs/ARCHITECTURE.md).

High-level overview:
- **Backend**: Layered architecture (API → Service → Repository → Database)
- **Frontend**: Component-based React architecture
- **Communication**: REST API with JSON
- **Deployment**: Containerized on Google Cloud Run

## Local Development

If you want to run the application locally for development or testing, see [LOCAL_SETUP.md](docs/LOCAL_SETUP.md) for detailed setup instructions.

**Note**: All API routes can be tested against the deployed backend - local setup is optional and only needed if you want to modify the code.

## Documentation

- [ARCHITECTURE.md](docs/ARCHITECTURE.md) - System architecture and design diagrams
- [API_ROUTES.md](docs/API_ROUTES.md) - Complete API testing guide with examples
- [LOCAL_SETUP.md](docs/LOCAL_SETUP.md) - Local development setup guide
- [DESIGN_DECISIONS.md](docs/DESIGN_DECISIONS.md) - Key technical decisions and reasoning

## Project Structure

```
merget-takehome/
├── backend/              # Go backend API
│   ├── cmd/server/      # Application entry point
│   ├── internal/        # Internal packages
│   │   ├── api/         # HTTP handlers & routes
│   │   ├── services/    # Business logic
│   │   ├── repository/  # Data access layer
│   │   ├── models/      # Data models
│   │   ├── database/    # DB connection, migrations, seeding
│   │   └── config/      # Configuration
│   ├── Dockerfile
│   └── cloudbuild.yaml
│
├── frontend/             # React frontend
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── api/         # API client
│   │   └── App.jsx      # Main component
│   ├── Dockerfile
│   └── cloudbuild.yaml
│
└── docs/                 # Documentation
    ├── ARCHITECTURE.md  # Architecture diagrams
    ├── API_ROUTES.md    # API testing guide
    ├── LOCAL_SETUP.md   # Local development guide
    └── DESIGN_DECISIONS.md # Technical decisions
```

## Security

- **API Key Authentication**: Required for external API requests
- **CORS Protection**: Whitelisted origins for frontend access
- **HTTPS**: All communication encrypted in production
- **Frontend Requests**: Automatically whitelisted (no API key needed)

## License

This project is part of a take-home interview assessment.
