package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type Order struct {
	ID             int       `json:"id"`
	ProductID      string    `json:"product_id"`
	Quantity       int       `json:"quantity"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	ProductName    string    `json:"product_name"`
	StockAvailable int       `json:"stock_available"`
	Warehouse      string    `json:"warehouse"`
}

var orderServiceURL string

func TestMain(m *testing.M) {
	orderServiceURL = os.Getenv("ORDER_SERVICE_URL")
	if orderServiceURL == "" {
		orderServiceURL = "http://localhost:8080"
	}
	os.Exit(m.Run())
}

func TestGetOrdersReturnsOK(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/orders", orderServiceURL))
	if err != nil {
		t.Fatalf("failed to call /orders: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(body))
	}
}

func TestGetOrdersReturnsJSON(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/orders", orderServiceURL))
	if err != nil {
		t.Fatalf("failed to call /orders: %v", err)
	}
	defer resp.Body.Close()
	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}

func TestGetOrdersReturnsNonEmpty(t *testing.T) {
	orders := fetchOrders(t)
	if len(orders) == 0 {
		t.Fatal("expected at least 1 order, got 0")
	}
}

func TestGetOrdersHasRequiredFields(t *testing.T) {
	orders := fetchOrders(t)
	if len(orders) == 0 {
		t.Fatal("no orders returned")
	}
	o := orders[0]
	if o.ID == 0 {
		t.Error("order ID should not be 0")
	}
	if o.ProductID == "" {
		t.Error("product_id should not be empty")
	}
	if o.Quantity == 0 {
		t.Error("quantity should not be 0")
	}
	if o.Status == "" {
		t.Error("status should not be empty")
	}
}

func TestGetOrdersEnrichedWithStock(t *testing.T) {
	orders := fetchOrders(t)
	if len(orders) == 0 {
		t.Fatal("no orders returned")
	}
	found := false
	for _, o := range orders {
		if o.ProductName != "" {
			found = true
			if o.Warehouse == "" {
				t.Errorf("order %d has product_name but no warehouse", o.ID)
			}
			break
		}
	}
	if !found {
		t.Error("expected at least one order enriched with product_name from inventory-service")
	}
}

func TestGetOrdersContainsExpectedProducts(t *testing.T) {
	orders := fetchOrders(t)
	productIDs := make(map[string]bool)
	for _, o := range orders {
		productIDs[o.ProductID] = true
	}
	expected := []string{"PROD-001", "PROD-003"}
	for _, pid := range expected {
		if !productIDs[pid] {
			t.Errorf("expected product %s in orders", pid)
		}
	}
}

func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/health", orderServiceURL))
	if err != nil {
		t.Fatalf("health check failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("health check returned %d", resp.StatusCode)
	}
}

func fetchOrders(t *testing.T) []Order {
	t.Helper()
	resp, err := http.Get(fmt.Sprintf("%s/orders", orderServiceURL))
	if err != nil {
		t.Fatalf("failed to call /orders: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var orders []Order
	if err := json.Unmarshal(body, &orders); err != nil {
		t.Fatalf("failed to parse response: %v\nbody: %s", err, string(body))
	}
	return orders
}
