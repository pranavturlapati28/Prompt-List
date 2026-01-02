# Local Development Setup Guide

This guide explains how to set up and run the application locally for development.

## Prerequisites

- Go 1.23+
- Node.js 20+
- PostgreSQL (running locally or accessible)
- Docker (optional, for containerized development)

## Backend Setup

### 1. Navigate to Backend Directory

```bash
cd backend
```

### 2. Set Environment Variables

Create a `.env` file or export environment variables:

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/prompttree?sslmode=disable"
export PORT="8080"
export API_KEY="your-api-key-here"  # Optional for local dev
export ENVIRONMENT="development"
```

Or create a `.env` file:
```bash
DATABASE_URL=postgres://user:pass@localhost:5432/prompttree?sslmode=disable
PORT=8080
API_KEY=your-api-key-here
ENVIRONMENT=development
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run the Server

```bash
go run cmd/server/main.go
```

The server will:
- Connect to PostgreSQL
- Run database migrations (create tables)
- Seed initial data
- Start on `http://localhost:8080`

### 5. Verify Backend is Running

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok","message":"Prompt Tree API is running"}
```

## Frontend Setup

### 1. Navigate to Frontend Directory

```bash
cd frontend
```

### 2. Install Dependencies

```bash
npm install
```

### 3. Set Environment Variable (Optional)

The frontend defaults to `http://localhost:8080` for the API. To override:

```bash
export VITE_API_URL="http://localhost:8080"
```

Or create a `.env` file:
```bash
VITE_API_URL=http://localhost:8080
```

### 4. Start Development Server

```bash
npm run dev
```

Frontend will be available at `http://localhost:5173`

### 5. Verify Frontend is Running

Open `http://localhost:5173` in your browser. You should see the prompt tree interface.

## Database Setup

### Using Local PostgreSQL

1. Install PostgreSQL if not already installed
2. Create a database:
```sql
CREATE DATABASE prompttree;
```

3. Update `DATABASE_URL` in your environment variables:
```bash
export DATABASE_URL="postgres://username:password@localhost:5432/prompttree?sslmode=disable"
```

### Using Docker PostgreSQL

```bash
docker run --name prompttree-db \
  -e POSTGRES_USER=promptuser \
  -e POSTGRES_PASSWORD=promptpass \
  -e POSTGRES_DB=prompttree \
  -p 5432:5432 \
  -d postgres:15
```

Then use:
```bash
export DATABASE_URL="postgres://promptuser:promptpass@localhost:5432/prompttree?sslmode=disable"
```

## Testing Locally

### Backend API Tests

```bash
# Health check
curl http://localhost:8080/health

# Get tree (if API_KEY is set)
curl -H "Authorization: Bearer $API_KEY" http://localhost:8080/tree

# Without API key (if API_KEY env var is not set, all requests are allowed)
curl http://localhost:8080/tree
```

### Frontend Development

- Hot reload is enabled - changes to code will automatically refresh
- Open browser DevTools to see console logs and network requests
- The frontend will automatically connect to the local backend

## Troubleshooting

### Backend Issues

**Database connection errors:**
- Verify PostgreSQL is running: `pg_isready`
- Check `DATABASE_URL` is correct
- Ensure database exists

**Port already in use:**
- Change `PORT` environment variable
- Or stop the process using port 8080

### Frontend Issues

**Cannot connect to backend:**
- Verify backend is running on `http://localhost:8080`
- Check `VITE_API_URL` matches backend URL
- Check browser console for CORS errors

**Build errors:**
- Delete `node_modules` and reinstall: `rm -rf node_modules && npm install`
- Clear Vite cache: `rm -rf node_modules/.vite`

## Development Workflow

1. Start PostgreSQL database
2. Start backend: `cd backend && go run cmd/server/main.go`
3. Start frontend: `cd frontend && npm run dev`
4. Make changes to code
5. Test changes in browser at `http://localhost:5173`

## Building for Production

### Backend

```bash
cd backend
go build -o server ./cmd/server
./server
```

### Frontend

```bash
cd frontend
npm run build
```

Build output will be in `frontend/dist/` directory.

## Docker Development

### Backend

```bash
cd backend
docker build -t prompttree-backend .
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  prompttree-backend
```

### Frontend

```bash
cd frontend
docker build -t prompttree-frontend .
docker run -p 5173:80 prompttree-frontend
```

## Notes

- The backend will automatically run migrations and seed data on startup
- If you want to reset the database, drop and recreate it, then restart the backend
- API key authentication is optional for local development (if `API_KEY` is not set, all requests are allowed)

