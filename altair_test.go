package altair

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHandler(t *testing.T) {
	// Define a test configuration
	cfg := &Config{
		Force:              false,
		DefaultWindowTitle: "Altair",
		Endpoint:           "http://example.com/graphql",
		Headers: []Header{
			{Key: "Authorization", Value: "Bearer Token"},
		},
	}

	// Create a fake HTTP request
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()

	// Call the Handler function
	handler := Handler(cfg)
	handler(rec, req)

	// Check the response status code
	if rec.Code != http.StatusOK {
		t.Errorf("Handler returned unexpected status code: got %d, want %d", rec.Code, http.StatusOK)
	}

	// You can also check the response body if needed
	// responseBody := rec.Body.String()
	// ...

	// Add more assertions as needed
}

func TestGetSubscriptionAbsoluteEndpoint(t *testing.T) {
	tests := []struct {
		endpoint string
		expected string
	}{
		{"http://example.com/graphql", "ws://example.com/graphql"},
		{"https://example.com/graphql", "wss://example.com/graphql"},
		{"invalid-url", ""},
		// Add more test cases as needed
	}

	for _, test := range tests {
		result := getSubscriptionAbsoluteEndpoint(test.endpoint)
		if result != test.expected {
			t.Errorf("getSubscriptionAbsoluteEndpoint(%s) returned %s, expected %s", test.endpoint, result, test.expected)
		}
	}
}

func TestEndpointToAbsolute(t *testing.T) {
	tests := []struct {
		endpoint       string
		initialHost    string
		initialScheme  string
		expectedResult string
	}{
		{"http://example.com/graphql", "example.com", "http", "http://example.com/graphql"},
		{"//example.com/graphql", "example.com", "http", "http://example.com/graphql"},
		{"example.com/graphql", "", "http", "http://example.com/graphql"},
		{"ws://example.com/graphql", "example.com", "http", "ws://example.com/graphql"},
		{"http://example.com/graphql", "", "", "http://example.com/graphql"},
		{"http://example.com/graphql", "example.com", "", "http://example.com/graphql"},
		{"http://example.com/graphql", "", "https", "http://example.com/graphql"},
	}

	for _, test := range tests {
		result := endpointToAbsolute(test.endpoint, test.initialHost, test.initialScheme)
		if result != test.expectedResult {
			t.Errorf("endpointToAbsolute(%s, %s, %s) returned %s, expected %s",
				test.endpoint, test.initialHost, test.initialScheme, result, test.expectedResult)
		}
	}
}

func TestEchoHandler(t *testing.T) {
	// Define a test configuration
	cfg := &Config{
		Force:              false,
		DefaultWindowTitle: "Altair",
		Endpoint:           "http://example.com/graphql",
		Headers: []Header{
			{Key: "Authorization", Value: "Bearer Token"},
		},
	}

	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request using Echo's `NewRequest` method
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the EchoHandler function
	err := EchoHandler(cfg)(c)

	// Check for any errors
	if err != nil {
		t.Errorf("EchoHandler returned an error: %v", err)
	}

	// Check the response status code
	if rec.Code != http.StatusOK {
		t.Errorf("EchoHandler returned unexpected status code: got %d, want %d", rec.Code, http.StatusOK)
	}
}
