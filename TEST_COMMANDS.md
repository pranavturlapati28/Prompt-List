# Test Commands for Prompt Tree API

Make sure your backend server is running on `http://localhost:8080`

## Basic Endpoints

```bash
# 1. Health check
curl http://localhost:8080/health

# 2. Get full tree
curl http://localhost:8080/tree

# 3. Get single prompt (ID: 1)
curl http://localhost:8080/prompts/1

# 4. Get nodes for prompt 1
curl http://localhost:8080/prompts/1/nodes

# 5. Get notes for prompt 1
curl http://localhost:8080/prompts/1/notes
```

## Create Operations

```bash
# 6. Create a new prompt
curl -X POST http://localhost:8080/prompts/0 \
  -H "Content-Type: application/json" \
  -d '{"title":"New Prompt","description":"Test description"}'

# 7. Create a node for prompt 1
curl -X POST http://localhost:8080/prompts/1/nodes \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Node","action":"Test action"}'

# 8. Create a note for prompt 1
curl -X POST http://localhost:8080/prompts/1/notes \
  -H "Content-Type: application/json" \
  -d '{"content":"This is a test note"}'
```

## SSE Webhook Testing

```bash
# 9. Test SSE webhook (will stream events - press Ctrl+C to stop)
curl -N http://localhost:8080/events

# 10. Test SSE with better formatting (shows events as they come)
curl -N -H "Accept: text/event-stream" http://localhost:8080/events
```

## Testing Webhook in Action

**Terminal 1** - Listen for events:
```bash
curl -N http://localhost:8080/events
```

**Terminal 2** - Trigger events by creating data:
```bash
# This should trigger a "tree_changed" event
curl -X POST http://localhost:8080/prompts/0 \
  -H "Content-Type: application/json" \
  -d '{"title":"Webhook Test","description":"Testing real-time updates"}'

# This should trigger "tree_changed" and "node_changed" events
curl -X POST http://localhost:8080/prompts/1/nodes \
  -H "Content-Type: application/json" \
  -d '{"name":"Webhook Node","action":"Test webhook"}'

# This should trigger a "note_changed" event
curl -X POST http://localhost:8080/prompts/1/notes \
  -H "Content-Type: application/json" \
  -d '{"content":"Webhook test note"}'
```

## OpenAPI Documentation

```bash
# View API documentation
curl http://localhost:8080/docs
# Or open in browser: http://localhost:8080/docs
```

