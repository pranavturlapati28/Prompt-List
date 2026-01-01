#!/bin/bash

# Test script for Prompt Tree API routes
# Make sure the backend server is running on http://localhost:8080

API_URL="${API_URL:-http://localhost:8080}"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Testing Prompt Tree API Routes"
echo "API URL: $API_URL"
echo "=========================================="
echo ""

# Test counter
PASSED=0
FAILED=0

# Helper function to test endpoints
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -n "Testing $description... "
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$API_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$API_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} (HTTP $http_code)"
        echo "  Response: $body"
        ((FAILED++))
        return 1
    fi
}

# Test 1: Health check
test_endpoint "GET" "/health" "" "Health Check"

# Test 2: Get tree
test_endpoint "GET" "/tree" "" "Get Tree"

# Test 3: Get single prompt (assuming prompt ID 1 exists)
test_endpoint "GET" "/prompts/1" "" "Get Prompt (ID: 1)"

# Test 4: Create a new prompt
PROMPT_DATA='{"title":"Test Prompt","description":"This is a test prompt"}'
test_endpoint "POST" "/prompts/0" "$PROMPT_DATA" "Create Prompt"

# Get the ID of the created prompt (if creation was successful)
# Note: This is a simplified test - in reality you'd parse the response
NEW_PROMPT_ID=2  # Assuming it creates ID 2

# Test 5: Get nodes for a prompt
test_endpoint "GET" "/prompts/1/nodes" "" "Get Prompt Nodes"

# Test 6: Create a node
NODE_DATA='{"name":"Test Node","action":"Test action"}'
test_endpoint "POST" "/prompts/1/nodes" "$NODE_DATA" "Create Node"

# Test 7: Get notes for a prompt
test_endpoint "GET" "/prompts/1/notes" "" "Get Notes"

# Test 8: Create a note
NOTE_DATA='{"content":"This is a test note"}'
test_endpoint "POST" "/prompts/1/notes" "$NOTE_DATA" "Create Note"

# Test 9: Test SSE endpoint (webhook)
echo -n "Testing SSE Webhook (/events)... "
timeout 3 curl -s -N "$API_URL/events" > /dev/null 2>&1
if [ $? -eq 124 ] || [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ PASS${NC} (SSE connection established)"
    ((PASSED++))
else
    echo -e "${RED}✗ FAIL${NC} (Could not connect to SSE endpoint)"
    ((FAILED++))
fi

# Test 10: Test OpenAPI docs
echo -n "Testing OpenAPI Docs (/docs)... "
response=$(curl -s -w "\n%{http_code}" "$API_URL/docs")
http_code=$(echo "$response" | tail -n1)
if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠ WARN${NC} (HTTP $http_code - docs might be at different path)"
fi

echo ""
echo "=========================================="
echo "Test Results:"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo "=========================================="

if [ $FAILED -eq 0 ]; then
    exit 0
else
    exit 1
fi

