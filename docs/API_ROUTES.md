# API Testing Guide

This guide shows you how to test all the API endpoints using curl commands.

**What you need:**
- Backend URL: `https://backend-709459926380.us-central1.run.app`
- API key (provided separately)
- Terminal/command line with `curl` installed

**Important:** Replace `<BACKEND_URL>` with your backend URL and `<YOUR_API_KEY>` with your API key in all commands below.

---

## Quick Tips

1. **Set environment variables** to make testing easier:
   ```bash
   export API_KEY="<YOUR_API_KEY>"
   export BACKEND_URL="<BACKEND_URL>"
   ```
   Then use: `curl -H "Authorization: Bearer $API_KEY" "$BACKEND_URL/tree"`

2. **Pretty print JSON** responses:
   ```bash
   curl -H "Authorization: Bearer <YOUR_API_KEY>" <BACKEND_URL>/tree | jq .
   ```
   (Requires `jq` to be installed)

3. **Test authentication** first:
   ```bash
   # This should fail
   curl <BACKEND_URL>/tree
   
   # This should work
   curl -H "Authorization: Bearer <YOUR_API_KEY>" <BACKEND_URL>/tree
   ```


## Authentication

All endpoints (except `/health` and `/docs`) require an API key. Add this header to every request:

```bash
-H "Authorization: Bearer <YOUR_API_KEY>"
```

## API Endpoints

### Health Check
Check if the API is running.

```bash
curl -H "Authorization: Bearer $API_KEY" $BACKEND_URL/tree
```

**Sample response:**
```json
{"status":"ok"}
```

---

### Get Full Tree
Get the complete prompt tree with all prompts and their nodes.

```bash
curl -H "Authorization: Bearer <YOUR_API_KEY>" <BACKEND_URL>/tree
```

**Sample response:**
```json
{
  "project": "Personal Finance Copilot",
  "mainRequest": "Build a web app...",
  "prompts": [...]
}
```

---

### Get Single Prompt

You can see prompt ID here:

<img width="885" height="401" alt="image" src="https://github.com/user-attachments/assets/e819b020-8a0b-4f7d-81d4-13c9acfa9a8a" />

Get details for a specific prompt by ID.

```bash
curl -H "Authorization: Bearer <YOUR_API_KEY>" <BACKEND_URL>/prompts/1
```

**Sample response:**
```json
{
  "id": 1,
  "title": "Project Setup",
  "description": "Initialize repo...",
  "nodes": [...]
}
```

---

### Get Nodes for a Prompt
Get all nodes (subprompts) for a specific prompt.

```bash
curl -H "Authorization: Bearer <YOUR_API_KEY>" <BACKEND_URL>/prompts/1/nodes
```

**Sample response:**
```json
[
  {
    "id": 1,
    "name": "Scaffold app",
    "action": "Create a Vite + React project..."
  }
]
```

---

### Get Notes for a Prompt
Get all notes/annotations for a specific prompt.

```bash
curl -H "Authorization: Bearer <YOUR_API_KEY>" <BACKEND_URL>/prompts/1/notes
```

**Sample response:**
```json
[
  {
    "id": 1,
    "content": "This is a note about the prompt"
  }
]
```

---

### Export Tree as JSON
Export the current tree structure as JSON.

```bash
curl -H "Authorization: Bearer <YOUR_API_KEY>" <BACKEND_URL>/tree/export
```

**Sample response:**
```json
{
  "project": "Personal Finance Copilot",
  "mainRequest": "...",
  "prompts": [...]
}
```

---

### Import Tree from JSON
Replace the current tree with a new one from JSON.

```bash
curl -X POST <BACKEND_URL>/tree/import \
  -H "Authorization: Bearer <YOUR_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{
    "project": "New Project",
    "mainRequest": "New main request",
    "prompts": [
      {
        "id": 1,
        "title": "New Prompt",
        "description": "Description here",
        "subprompts": [
          {
            "name": "Node Name",
            "action": "Node action"
          }
        ]
      }
    ]
  }'
```

**Sample response:**
```json
{"message":"Tree imported successfully"}
```

---

### Save Current Tree
Save the current tree with a name for later retrieval.

```bash
curl -X POST <BACKEND_URL>/tree/save \
  -H "Authorization: Bearer <YOUR_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"name":"my-saved-tree"}'
```

**Sample response:**
```json
{"message":"Tree saved as my-saved-tree"}
```

---

### List Saved Trees
Get a list of all saved tree names.

```bash
curl -H "Authorization: Bearer <YOUR_API_KEY>" <BACKEND_URL>/tree/saves
```

**Sample response:**
```json
{
  "trees": [
    {"name": "my-saved-tree", "savedAt": "2024-01-01T12:00:00Z"},
    {"name": "another-tree", "savedAt": "2024-01-02T10:00:00Z"}
  ]
}
```

---

### Load Saved Tree
Load a previously saved tree by name.

```bash
curl -X POST <BACKEND_URL>/tree/load/my-saved-tree \
  -H "Authorization: Bearer <YOUR_API_KEY>"
```

**Sample response:**
```json
{"message":"Tree loaded successfully"}
```

---

### Delete Saved Tree
Delete a saved tree by name.

```bash
curl -X DELETE <BACKEND_URL>/tree/saves/my-saved-tree \
  -H "Authorization: Bearer <YOUR_API_KEY>"
```

**Sample response:**
```json
{"message":"Tree deleted successfully"}
```

---

### Create New Prompt
Create a new prompt in the tree.

```bash
curl -X POST <BACKEND_URL>/prompts/0 \
  -H "Authorization: Bearer <YOUR_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"title":"New Prompt","description":"This is a new prompt"}'
```

**Sample response:**
```json
{
  "id": 7,
  "title": "New Prompt",
  "description": "This is a new prompt"
}
```

---

### Update Prompt
Update an existing prompt by ID.

```bash
curl -X PUT <BACKEND_URL>/prompts/1 \
  -H "Authorization: Bearer <YOUR_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","description":"Updated description"}'
```

**Sample response:**
```json
{
  "id": 1,
  "title": "Updated Title",
  "description": "Updated description"
}
```

---

### Delete Prompt
Delete a prompt by ID (also deletes all associated nodes and notes).

```bash
curl -X DELETE <BACKEND_URL>/prompts/1 \
  -H "Authorization: Bearer <YOUR_API_KEY>"
```

**Sample response:**
```json
{"message":"Prompt deleted successfully"}
```

---

### Create Node
Create a new node (subprompt) for a specific prompt.

```bash
curl -X POST <BACKEND_URL>/prompts/1/nodes \
  -H "Authorization: Bearer <YOUR_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"name":"New Node","action":"Action description here"}'
```

**Sample response:**
```json
{
  "id": 10,
  "name": "New Node",
  "action": "Action description here"
}
```

---

### Update Node
Update an existing node by ID.

```bash
curl -X PUT <BACKEND_URL>/prompts/1/nodes/5 \
  -H "Authorization: Bearer <YOUR_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Node Name","action":"Updated action"}'
```

**Sample response:**
```json
{
  "id": 5,
  "name": "Updated Node Name",
  "action": "Updated action"
}
```

---

### Delete Node
Delete a node by ID.

```bash
curl -X DELETE <BACKEND_URL>/prompts/1/nodes/5 \
  -H "Authorization: Bearer <YOUR_API_KEY>"
```

**Sample response:**
```json
{"message":"Node deleted successfully"}
```

---

### Create Note
Add a note/annotation to a prompt.

```bash
curl -X POST <BACKEND_URL>/prompts/1/notes \
  -H "Authorization: Bearer <YOUR_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"content":"This is my note about this prompt"}'
```

**Sample response:**
```json
{
  "id": 3,
  "content": "This is my note about this prompt"
}
```

---

### Update Note
Update an existing note by ID.

```bash
curl -X PUT <BACKEND_URL>/prompts/1/notes/3 \
  -H "Authorization: Bearer <YOUR_API_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"content":"Updated note content"}'
```

**Sample response:**
```json
{
  "id": 3,
  "content": "Updated note content"
}
```

---

### Delete Note
Delete a note by ID.

```bash
curl -X DELETE <BACKEND_URL>/prompts/1/notes/3 \
  -H "Authorization: Bearer <YOUR_API_KEY>"
```

**Sample response:**
```json
{"message":"Note deleted successfully"}
```

---

### View API Documentation
View interactive API documentation (no authentication needed).

```bash
curl <BACKEND_URL>/docs
```

Or open in your browser: `<BACKEND_URL>/docs`

---

## Common Errors

### 401 Unauthorized
You're missing the API key or it's incorrect.

```json
{"error":"API key required. Use Authorization: Bearer <your-api-key>"}
```

**Fix:** Add the `Authorization` header with your API key.

### 404 Not Found
The resource doesn't exist (wrong ID, etc.).

**Fix:** Check that the ID exists.

### 422 Unprocessable Entity
The request body format is invalid.

**Fix:** Check your JSON format matches the examples above.

