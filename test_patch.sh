#!/bin/bash

BASE_URL="http://localhost:8080"

echo "🧪 Testing PATCH /items/:id endpoint"
echo "===================================="

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Update name
echo -e "\n${YELLOW}Test 1: Update name only${NC}"
RESPONSE=$(curl -s -X PATCH "$BASE_URL/items/1" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Item Name"}' \
  -w "\nHTTP_CODE:%{http_code}")

HTTP_CODE=$(echo "$RESPONSE" | grep -oP 'HTTP_CODE:\K\d+')
BODY=$(echo "$RESPONSE" | sed 's/HTTP_CODE:.*//')

if [ "$HTTP_CODE" = "200" ]; then
  echo -e "${GREEN}✓ PASS${NC} - Status: $HTTP_CODE"
  echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
  echo -e "${RED}✗ FAIL${NC} - Status: $HTTP_CODE"
  echo "$BODY"
fi

# Test 2: Update price
echo -e "\n${YELLOW}Test 2: Update purchase_price only${NC}"
RESPONSE=$(curl -s -X PATCH "$BASE_URL/items/1" \
  -H "Content-Type: application/json" \
  -d '{"purchase_price": 2000000}' \
  -w "\nHTTP_CODE:%{http_code}")

HTTP_CODE=$(echo "$RESPONSE" | grep -oP 'HTTP_CODE:\K\d+')
BODY=$(echo "$RESPONSE" | sed 's/HTTP_CODE:.*//')

if [ "$HTTP_CODE" = "200" ]; then
  echo -e "${GREEN}✓ PASS${NC} - Status: $HTTP_CODE"
  echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
  echo -e "${RED}✗ FAIL${NC} - Status: $HTTP_CODE"
  echo "$BODY"
fi

# Test 3: Invalid price (should fail)
echo -e "\n${YELLOW}Test 3: Invalid price (negative) - should return 400${NC}"
RESPONSE=$(curl -s -X PATCH "$BASE_URL/items/1" \
  -H "Content-Type: application/json" \
  -d '{"purchase_price": -100}' \
  -w "\nHTTP_CODE:%{http_code}")

HTTP_CODE=$(echo "$RESPONSE" | grep -oP 'HTTP_CODE:\K\d+')
BODY=$(echo "$RESPONSE" | sed 's/HTTP_CODE:.*//')

if [ "$HTTP_CODE" = "400" ]; then
  echo -e "${GREEN}✓ PASS${NC} - Status: $HTTP_CODE (expected)"
  echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
  echo -e "${RED}✗ FAIL${NC} - Expected 400, got $HTTP_CODE"
  echo "$BODY"
fi

# Test 4: Immutable field - id (should fail)
echo -e "\n${YELLOW}Test 4: Immutable field (id) - should return 400${NC}"
RESPONSE=$(curl -s -X PATCH "$BASE_URL/items/1" \
  -H "Content-Type: application/json" \
  -d '{"id": 999, "name": "Updated"}' \
  -w "\nHTTP_CODE:%{http_code}")

HTTP_CODE=$(echo "$RESPONSE" | grep -oP 'HTTP_CODE:\K\d+')
BODY=$(echo "$RESPONSE" | sed 's/HTTP_CODE:.*//')

if [ "$HTTP_CODE" = "400" ]; then
  echo -e "${GREEN}✓ PASS${NC} - Status: $HTTP_CODE (expected)"
  echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
  echo -e "${RED}✗ FAIL${NC} - Expected 400, got $HTTP_CODE"
  echo "$BODY"
fi

# Test 5: Item not found (should fail)
echo -e "\n${YELLOW}Test 5: Item not found - should return 404${NC}"
RESPONSE=$(curl -s -X PATCH "$BASE_URL/items/9999" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Name"}' \
  -w "\nHTTP_CODE:%{http_code}")

HTTP_CODE=$(echo "$RESPONSE" | grep -oP 'HTTP_CODE:\K\d+')
BODY=$(echo "$RESPONSE" | sed 's/HTTP_CODE:.*//')

if [ "$HTTP_CODE" = "404" ]; then
  echo -e "${GREEN}✓ PASS${NC} - Status: $HTTP_CODE (expected)"
  echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
  echo -e "${RED}✗ FAIL${NC} - Expected 404, got $HTTP_CODE"
  echo "$BODY"
fi

# Test 6: Update multiple fields
echo -e "\n${YELLOW}Test 6: Update multiple fields${NC}"
RESPONSE=$(curl -s -X PATCH "$BASE_URL/items/1" \
  -H "Content-Type: application/json" \
  -d '{"name": "Multi Update", "brand": "New Brand", "purchase_price": 1500000}' \
  -w "\nHTTP_CODE:%{http_code}")

HTTP_CODE=$(echo "$RESPONSE" | grep -oP 'HTTP_CODE:\K\d+')
BODY=$(echo "$RESPONSE" | sed 's/HTTP_CODE:.*//')

if [ "$HTTP_CODE" = "200" ]; then
  echo -e "${GREEN}✓ PASS${NC} - Status: $HTTP_CODE"
  echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
  echo -e "${RED}✗ FAIL${NC} - Status: $HTTP_CODE"
  echo "$BODY"
fi

echo -e "\n${GREEN}Testing complete!${NC}"

