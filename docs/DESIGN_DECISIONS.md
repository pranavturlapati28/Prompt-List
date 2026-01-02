# Design Decisions

This document explains the key technical decisions made in building this application and the reasoning behind them.

## Backend Language: Go

**What:** Backend is written in Go 1.23

**Why:**
- **Performance**: Go compiles to native binaries, providing excellent performance for API servers
- **Simplicity**: Clean syntax and minimal boilerplate make code easy to read and maintain
- **Concurrency**: Built-in goroutines handle concurrent requests efficiently
- **Deployment**: Single binary deployment simplifies containerization and reduces image size
- **Type Safety**: Strong typing catches errors at compile time

## Backend Architecture: Layered Architecture

**What:** Three-layer architecture: API → Service → Repository

**Why:**
- **Separation of Concerns**: Each layer has a single, clear responsibility
- **Testability**: Each layer can be tested independently
- **Maintainability**: Changes in one layer don't affect others
- **Scalability**: Easy to swap implementations (e.g., different database)
- **Industry Standard**: Well-understood pattern that any developer can follow

**Layers:**
1. **API Layer**: Handles HTTP requests/responses, validation
2. **Service Layer**: Contains business logic and rules
3. **Repository Layer**: Manages database operations

## API Framework: Huma

**What:** Huma v2 framework for API definition and OpenAPI generation

**Why:**
- **Auto-Documentation**: Automatically generates OpenAPI/Swagger docs
- **Type Safety**: Request/response validation built-in
- **Developer Experience**: Clear API contracts reduce errors
- **Standards Compliance**: Follows OpenAPI standards
- **Minimal Boilerplate**: Less code than manual validation

## HTTP Router: Chi

**What:** Chi router for HTTP routing

**Why:**
- **Lightweight**: Minimal overhead, fast routing
- **Middleware Support**: Easy to add CORS, auth, logging
- **Standard Library**: Built on `net/http`, no vendor lock-in
- **Flexibility**: Works well with Huma framework
- **Simplicity**: Easy to understand and debug

## Database: PostgreSQL

**What:** PostgreSQL relational database

**Why:**
- **Reliability**: ACID compliance ensures data integrity
- **Relationships**: Foreign keys maintain referential integrity
- **JSON Support**: JSONB type for flexible saved tree storage
- **Mature**: Battle-tested, widely supported
- **Cloud SQL**: Easy integration with Google Cloud SQL

## Authentication: API Key

**What:** API key authentication using Bearer tokens

**Why:**
- **Simplicity**: Easy to implement and use
- **Stateless**: No session management needed
- **External Access**: Perfect for programmatic API access
- **Frontend Bypass**: Whitelisted origins don't need keys (better UX)
- **Security**: Sufficient for take-home project scope

**Note:** For production, consider OAuth2 or JWT for more advanced scenarios.

## Frontend Framework: React

**What:** React 18 with functional components

**Why:**
- **Component Reusability**: Modular components reduce code duplication
- **State Management**: Built-in hooks handle application state
- **Ecosystem**: Large community and extensive libraries
- **Performance**: Virtual DOM for efficient updates
- **Developer Experience**: Hot reload, great tooling

## Frontend Build Tool: Vite

**What:** Vite for building and development

**Why:**
- **Speed**: Lightning-fast development server
- **Modern**: Native ES modules, no bundling in dev
- **Simple**: Minimal configuration needed
- **Optimized Builds**: Efficient production builds
- **Better DX**: Faster feedback loop during development

## Frontend Architecture: Component-Based

**What:** Small, focused React components

**Why:**
- **Maintainability**: Each component has a single responsibility
- **Reusability**: Components can be reused across the app
- **Testability**: Easy to test individual components
- **Readability**: Clear component hierarchy
- **Scalability**: Easy to add new features

## Deployment: Google Cloud Run

**What:** Containerized deployment on Cloud Run

**Why:**
- **Serverless**: No server management, auto-scaling
- **Cost-Effective**: Pay only for what you use
- **Easy Deployment**: Simple `gcloud` commands
- **Container-Based**: Docker containers ensure consistency
- **HTTPS**: Automatic SSL certificates
- **Integration**: Easy connection to Cloud SQL

## Database Migrations: Manual SQL

**What:** SQL migrations in Go code

**Why:**
- **Simplicity**: No external migration tool needed
- **Version Control**: Migrations live with code
- **Transparency**: Easy to see what changes are made
- **Control**: Full control over migration process
- **Portability**: Works anywhere PostgreSQL runs

## State Management: React Hooks

**What:** useState and useEffect hooks for state

**Why:**
- **Built-in**: No external libraries needed
- **Simple**: Easy to understand and use
- **Sufficient**: Meets current requirements
- **Lightweight**: No additional bundle size
- **Standard**: Part of React, well-documented

## HTTP Client: Axios

**What:** Axios for API calls

**Why:**
- **Promise-Based**: Clean async/await syntax
- **Interceptors**: Easy to add auth headers globally
- **Error Handling**: Built-in error handling
- **Browser Support**: Works in all modern browsers
- **Familiar**: Widely used, well-documented

## CORS Strategy: Whitelist Origins

**What:** Only allow specific frontend origins

**Why:**
- **Security**: Prevents unauthorized access
- **Flexibility**: Frontend requests bypass API key (better UX)
- **Control**: Know exactly who can access the API
- **Production-Ready**: Proper security practice

## File Organization: Internal Packages

**What:** Go `internal/` directory structure

**Why:**
- **Encapsulation**: Prevents external imports
- **Clear Boundaries**: Makes dependencies explicit
- **Maintainability**: Easy to understand project structure
- **Go Best Practice**: Standard Go project layout

## Error Handling: JSON Error Responses

**What:** Consistent JSON error format

**Why:**
- **Consistency**: Same format across all endpoints
- **Client-Friendly**: Easy for frontend to handle
- **Debugging**: Clear error messages
- **Standards**: Follows REST API best practices

## Data Models: Separate DB and API Models

**What:** Different structs for database and API responses

**Why:**
- **Flexibility**: Can change API without changing database
- **Security**: Don't expose internal database structure
- **Versioning**: Easy to version API independently
- **Transformation**: Can format data for frontend needs

## Summary

These decisions prioritize:
1. **Simplicity**: Easy to understand and maintain
2. **Performance**: Fast and efficient
3. **Developer Experience**: Good tooling and documentation
4. **Scalability**: Can grow with needs
5. **Best Practices**: Following industry standards

Each choice balances current needs with future flexibility, making the codebase maintainable and extensible.

