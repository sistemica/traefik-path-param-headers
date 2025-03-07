package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/traefik-path-param-headers"
)

func TestPathParamHeaders(t *testing.T) {
	tests := []struct {
		name           string
		pathPattern    string
		requestPath    string
		expectedParams map[string]string
		shouldMatch    bool
	}{
		{
			name:        "Simple path with two parameters",
			pathPattern: "/products/{category}/{id}",
			requestPath: "/products/electronics/12345",
			expectedParams: map[string]string{
				"X-Path-Category": "electronics",
				"X-Path-Id":       "12345",
			},
			shouldMatch: true,
		},
		{
			name:        "Path with nested parameters",
			pathPattern: "/api/users/{userId}/posts/{postId}",
			requestPath: "/api/users/john123/posts/987",
			expectedParams: map[string]string{
				"X-Path-UserId": "john123",
				"X-Path-PostId": "987",
			},
			shouldMatch: true,
		},
		{
			name:        "Path with mixed static and parameter segments",
			pathPattern: "/api/{version}/users/{id}/profile",
			requestPath: "/api/v1/users/123/profile",
			expectedParams: map[string]string{
				"X-Path-Version": "v1",
				"X-Path-Id":      "123",
			},
			shouldMatch: true,
		},
		{
			name:        "Path that doesn't match pattern",
			pathPattern: "/products/{category}/{id}",
			requestPath: "/services/support/ticket",
			shouldMatch: false,
		},
		{
			name:        "Path with missing parameter",
			pathPattern: "/products/{category}/{id}",
			requestPath: "/products/electronics",
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create config
			config := pathparamheaders.CreateConfig()
			config.PathPattern = tt.pathPattern
			config.HeaderPrefix = "X-Path-"

			// Setup a test handler that captures headers
			var capturedHeaders http.Header
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedHeaders = r.Header.Clone()
				w.WriteHeader(http.StatusOK)
			})

			// Create middleware
			middleware, err := pathparamheaders.New(context.Background(), next, config, "path-param-headers")
			if err != nil {
				t.Fatalf("Error creating middleware: %v", err)
			}

			// Create test request
			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			recorder := httptest.NewRecorder()

			// Execute middleware
			middleware.ServeHTTP(recorder, req)

			// Check results
			if tt.shouldMatch {
				for name, expectedValue := range tt.expectedParams {
					actualValue := capturedHeaders.Get(name)
					if actualValue != expectedValue {
						t.Errorf("Header %s = %s, want %s", name, actualValue, expectedValue)
					}
				}
			} else {
				// Check that no path parameter headers were added
				hasPathHeaders := false
				for name := range capturedHeaders {
					if len(name) >= 7 && name[:7] == "X-Path-" {
						hasPathHeaders = true
						break
					}
				}
				if hasPathHeaders {
					t.Errorf("Expected no path parameter headers, but found some")
				}
			}
		})
	}
}
