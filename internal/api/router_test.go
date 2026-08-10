package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouterIntegration(t *testing.T) {
	// Initialize your fully wired central hub router
	hubRouter := NewRouter()

	// Simulate a full API request targeting the complete combined path.
	// Note: include the trailing slash to hit the root of the mounted sub-router.
	req, err := http.NewRequest("GET", "/api/v1/health/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a recorder to capture the routing outcome
	rr := httptest.NewRecorder()

	// Pass the request into the hub router
	hubRouter.ServeHTTP(rr, req)

	// Verification Assertion
	// If the central hub r.Mount() is wired properly, it passes the request 
	// to the sub-router and returns 200 OK. If your path string is broken, 
	// chi will return 404 Not Found.
	if rr.Code != http.StatusOK {
		t.Errorf("Central Hub routing failed! Expected status 200, but got %d. Verify your r.Mount path prefix matches your request URL.", rr.Code)
	}
}
