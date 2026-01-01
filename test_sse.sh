#!/bin/bash

# Test script specifically for SSE webhook endpoint
# This will connect to the /events endpoint and display incoming events

API_URL="${API_URL:-http://localhost:8080}"

echo "=========================================="
echo "Testing SSE Webhook Endpoint"
echo "Connecting to: $API_URL/events"
echo "=========================================="
echo ""
echo "This will listen for events for 30 seconds..."
echo "In another terminal, try creating a prompt, node, or note to see events!"
echo ""
echo "Press Ctrl+C to stop"
echo ""

# Connect to SSE endpoint and display events
curl -N -H "Accept: text/event-stream" "$API_URL/events" 2>/dev/null | while IFS= read -r line; do
    if [[ $line == data:* ]]; then
        # Extract JSON data from SSE format
        json_data=$(echo "$line" | sed 's/^data: //')
        echo "[$(date +%H:%M:%S)] Event received: $json_data"
    elif [[ $line == :* ]]; then
        # This is a comment/keepalive
        echo "[$(date +%H:%M:%S)] Keepalive ping"
    fi
done

