package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

)

func TestHealthEndpoint(t *testing.T) {
	// get the isolated router instance 
	router := Health()

	// create a fake Http Get Request targetting the root for this router
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request :%v ", err)
	}

	// create a new recorder 
	rr := httptest.NewRecorder()

	//send a fake request straight into the router's handler logic
	router.ServeHTTP(rr, req)

	// Assertions: verify everything 
	// check the HTTP status code 
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", rr.Code)
	}

	// Check the Content-Type header 
	expectedcontentType := "application/json"
	if actualContentType := rr.Header().Get("Content-Type"); actualContentType != expectedcontentType {
		t.Errorf("Expected Content-Type %q, but got %q", expectedcontentType, actualContentType)
	}

	// check the actualmessage body 
	expectedBody := `{"status": "UP"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body %q but got %q", expectedBody, rr.Body.String())
	}
}