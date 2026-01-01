# Webhook Testing Commands

## Test 1: Listen to Webhook Events

**Terminal 1** - Listen for events (leave this running):
```bash
curl -N http://localhost:8080/events
```

You should see:
```
data: {"type":"connected"}
```

## Test 2: Trigger Events

**Terminal 2** - Create data to trigger webhook events:

### Create a Prompt (triggers `tree_changed`)
```bash
curl -X POST http://localhost:8080/prompts/0 \
  -H "Content-Type: application/json" \
  -d '{"title":"Webhook Test Prompt","description":"Testing real-time updates"}'
```

**Expected:** Terminal 1 should show:
```
data: {"type":"tree_changed"}
```

### Create a Node (triggers `tree_changed`)
```bash
curl -X POST http://localhost:8080/prompts/1/nodes \
  -H "Content-Type: application/json" \
  -d '{"name":"Webhook Test Node","action":"Test webhook functionality"}'
```

**Expected:** Terminal 1 should show:
```
data: {"type":"tree_changed"}
```

### Create a Note (triggers `note_changed`)
```bash
curl -X POST http://localhost:8080/prompts/1/notes \
  -H "Content-Type: application/json" \
  -d '{"content":"This is a webhook test note"}'
```

**Expected:** Terminal 1 should show:
```
data: {"type":"note_changed","prompt_id":1}
```

## Test 3: Frontend Auto-Update

1. Open your frontend in a browser (http://localhost:5173 or your dev server URL)
2. Open browser console (F12)
3. In Terminal 2, create a prompt:
   ```bash
   curl -X POST http://localhost:8080/prompts/0 \
     -H "Content-Type: application/json" \
     -d '{"title":"Frontend Test","description":"Should auto-update"}'
   ```
4. **Expected:** The frontend tree should automatically refresh and show the new prompt!

## Quick Test Script

Run this to test all at once:
```bash
# Terminal 1: Listen
curl -N http://localhost:8080/events

# Terminal 2: Trigger events
curl -X POST http://localhost:8080/prompts/0 \
  -H "Content-Type: application/json" \
  -d '{"title":"Test 1","description":"Test"}'

sleep 2

curl -X POST http://localhost:8080/prompts/1/nodes \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Node","action":"Test"}'

sleep 2

curl -X POST http://localhost:8080/prompts/1/notes \
  -H "Content-Type: application/json" \
  -d '{"content":"Test note"}'
```

