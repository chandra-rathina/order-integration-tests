# order-integration-tests

Integration tests for order-service. Tests the GET /orders endpoint end-to-end including enrichment from inventory-service.

## Prerequisites
Before running tests, start:
1. Postgres on :5432
2. inventory-service on :8081
3. order-service on :8080

## Run Tests
```bash
export ORDER_SERVICE_URL=http://localhost:8080
go test -v ./...
```

## Test Cases
- `TestGetOrdersReturnsOK` — /orders returns HTTP 200
- `TestGetOrdersReturnsJSON` — response Content-Type is application/json
- `TestGetOrdersReturnsNonEmpty` — at least 1 order returned
- `TestGetOrdersHasRequiredFields` — orders have id, product_id, quantity, status
- `TestGetOrdersEnrichedWithStock` — orders include product_name + warehouse from inventory-service
- `TestGetOrdersContainsExpectedProducts` — PROD-001 and PROD-003 present
- `TestHealthEndpoint` — /health returns 200
