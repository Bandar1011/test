# Testing Guide for PATCH /items/:id Endpoint

## 1. Running Unit Tests

### Run all tests
```bash
go test ./...
```

### Run tests with verbose output
```bash
go test ./... -v
```

### Run only the PATCH endpoint tests
```bash
go test ./internal/interfaces/controller/items/... -v -run TestItemHandler_PatchItem
```

### Run tests with coverage
```bash
go test ./... -cover
```

## 2. Manual API Testing

### Prerequisites
1. Start the database:
```bash
docker-compose up -d mysql
```

2. Set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=items_db
```

3. Start the server:
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

### Test Cases

#### ✅ Success Case 1: Update name only
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Item Name"
  }'
```

**Expected Response:** 200 OK with updated item

#### ✅ Success Case 2: Update purchase_price only
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "purchase_price": 2000000
  }'
```

**Expected Response:** 200 OK with updated item

#### ✅ Success Case 3: Update multiple fields
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Name",
    "brand": "New Brand",
    "purchase_price": 1500000
  }'
```

**Expected Response:** 200 OK with updated item

#### ❌ Error Case 1: Item not found (404)
```bash
curl -X PATCH http://localhost:8080/items/9999 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Name"
  }'
```

**Expected Response:** 404 Not Found
```json
{
  "error": "item not found"
}
```

#### ❌ Error Case 2: Invalid price (400)
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "purchase_price": -100
  }'
```

**Expected Response:** 400 Bad Request
```json
{
  "error": "validation failed",
  "details": ["purchase_price must be >= 0"]
}
```

#### ❌ Error Case 3: Immutable field - id (400)
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "id": 999,
    "name": "Updated Name"
  }'
```

**Expected Response:** 400 Bad Request
```json
{
  "error": "validation failed",
  "details": ["id is immutable"]
}
```

#### ❌ Error Case 4: Immutable field - created_at (400)
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "created_at": "2023-01-01T00:00:00Z",
    "name": "Updated Name"
  }'
```

**Expected Response:** 400 Bad Request
```json
{
  "error": "validation failed",
  "details": ["created_at is immutable"]
}
```

#### ❌ Error Case 5: Immutable field - updated_at (400)
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "updated_at": "2023-01-01T00:00:00Z",
    "name": "Updated Name"
  }'
```

**Expected Response:** 400 Bad Request
```json
{
  "error": "validation failed",
  "details": ["updated_at is immutable"]
}
```

#### ❌ Error Case 6: Multiple immutable fields (400)
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "id": 999,
    "created_at": "2023-01-01T00:00:00Z",
    "name": "Updated Name"
  }'
```

**Expected Response:** 400 Bad Request
```json
{
  "error": "validation failed",
  "details": ["id is immutable", "created_at is immutable"]
}
```

#### ❌ Error Case 7: Name too long (400)
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "This is a very long name that exceeds one hundred characters and should fail validation because it is too long for the database field"
  }'
```

**Expected Response:** 400 Bad Request
```json
{
  "error": "validation failed",
  "details": ["name must be 100 characters or less"]
}
```

#### ❌ Error Case 8: Brand too long (400)
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "brand": "This is a very long brand name that exceeds one hundred characters and should fail validation because it is too long for the database field"
  }'
```

**Expected Response:** 400 Bad Request
```json
{
  "error": "validation failed",
  "details": ["brand must be 100 characters or less"]
}
```

#### ❌ Error Case 9: Invalid item ID format (400)
```bash
curl -X PATCH http://localhost:8080/items/invalid \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Name"
  }'
```

**Expected Response:** 400 Bad Request
```json
{
  "error": "invalid item ID"
}
```

## 3. Complete Test Workflow

### Step 1: Create an item first
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Item",
    "category": "時計",
    "brand": "ROLEX",
    "purchase_price": 1000000,
    "purchase_date": "2023-01-01"
  }'
```

Note the `id` from the response (e.g., `{"id": 1, ...}`)

### Step 2: Verify the item exists
```bash
curl -X GET http://localhost:8080/items/1
```

### Step 3: Update the item
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Test Item",
    "purchase_price": 1500000
  }'
```

### Step 4: Verify the update
```bash
curl -X GET http://localhost:8080/items/1
```

You should see:
- `name` changed to "Updated Test Item"
- `purchase_price` changed to 1500000
- `updated_at` timestamp updated
- `id`, `created_at` remain unchanged

## 4. Using Docker Compose (Full Stack)

### Start everything
```bash
docker-compose up -d
```

### Test the endpoint
```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Name"}'
```

### View logs
```bash
docker-compose logs -f app
```

### Stop everything
```bash
docker-compose down
```

## 5. Testing with jq (Pretty JSON Output)

If you have `jq` installed, use it for better output:

```bash
curl -X PATCH http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Name"}' | jq
```

## 6. Testing Script

Create a test script `test_patch.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "Testing PATCH /items/:id endpoint"
echo "=================================="

# Test 1: Update name
echo -e "\n1. Testing update name..."
curl -X PATCH "$BASE_URL/items/1" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Name"}' \
  -w "\nStatus: %{http_code}\n"

# Test 2: Update price
echo -e "\n2. Testing update price..."
curl -X PATCH "$BASE_URL/items/1" \
  -H "Content-Type: application/json" \
  -d '{"purchase_price": 2000000}' \
  -w "\nStatus: %{http_code}\n"

# Test 3: Invalid price
echo -e "\n3. Testing invalid price..."
curl -X PATCH "$BASE_URL/items/1" \
  -H "Content-Type: application/json" \
  -d '{"purchase_price": -100}' \
  -w "\nStatus: %{http_code}\n"

# Test 4: Immutable field
echo -e "\n4. Testing immutable field..."
curl -X PATCH "$BASE_URL/items/1" \
  -H "Content-Type: application/json" \
  -d '{"id": 999, "name": "Updated"}' \
  -w "\nStatus: %{http_code}\n"

echo -e "\nDone!"
```

Make it executable and run:
```bash
chmod +x test_patch.sh
./test_patch.sh
```

