package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
)

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	// Collect headers
	headers := make(map[string]string)
	pathParamHeaders := make(map[string]string)

	// Get all headers
	for name, values := range r.Header {
		headers[name] = strings.Join(values, ", ")

		// Collect X-Path-* headers separately
		if strings.HasPrefix(name, "X-Path-") {
			paramName := strings.TrimPrefix(name, "X-Path-")
			pathParamHeaders[paramName] = strings.Join(values, ", ")
		}
	}

	// Create response
	response := map[string]interface{}{
		"method":      r.Method,
		"path":        r.URL.Path,
		"headers":     headers,
		"pathParams":  pathParamHeaders,
		"queryParams": r.URL.Query(),
	}

	// Log path parameter headers
	log.Println("Path parameter headers:")
	var paramNames []string
	for name := range pathParamHeaders {
		paramNames = append(paramNames, name)
	}
	sort.Strings(paramNames)

	for _, name := range paramNames {
		log.Printf("  %s: %s", name, pathParamHeaders[name])
	}

	// Return response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
