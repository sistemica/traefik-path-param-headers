// Package pathparamheaders is a Traefik v3 middleware plugin that
// extracts path parameters based on a pattern and adds them as HTTP headers.
package pathparamheaders

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Config holds the plugin configuration.
type Config struct {
	// PathPattern is the URL pattern with parameters in curly braces
	// Example: "/products/{category}/{id}"
	PathPattern string `json:"pathPattern,omitempty"`

	// HeaderPrefix is an optional prefix for the header names
	// Example: "X-Path-" would result in headers like "X-Path-Category"
	HeaderPrefix string `json:"headerPrefix,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderPrefix: "X-Path-",
	}
}

// PathParamHeaders is a middleware plugin that extracts path parameters and adds them as headers.
type PathParamHeaders struct {
	next         http.Handler
	pathPattern  string
	headerPrefix string
}

// New creates a new PathParamHeaders middleware plugin.
func New(_ context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	if config.PathPattern == "" {
		return nil, fmt.Errorf("pathPattern cannot be empty")
	}

	return &PathParamHeaders{
		next:         next,
		pathPattern:  config.PathPattern,
		headerPrefix: config.HeaderPrefix,
	}, nil
}

// ServeHTTP implements the http.Handler interface.
func (p *PathParamHeaders) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Extract path parameters
	params, err := ExtractPathParamValues(p.pathPattern, req.URL.Path)
	if err == nil {
		// Add parameters as headers
		for name, value := range params {
			headerName := p.headerPrefix + strings.ToUpper(name[:1]) + name[1:]
			req.Header.Set(headerName, value)
		}
	}

	// Call the next handler
	p.next.ServeHTTP(rw, req)
}

// ExtractPathParams extracts parameter names from a URL path pattern
// For example, given "/products/{category}/{id}", it returns []string{"category", "id"}
func ExtractPathParams(pattern string) []string {
	re := regexp.MustCompile(`\{([^{}]+)\}`)
	matches := re.FindAllStringSubmatch(pattern, -1)

	params := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			paramName := strings.TrimSpace(match[1])
			params = append(params, paramName)
		}
	}

	return params
}

// ExtractPathParamValues extracts parameter values from an actual path based on a pattern
// For example, given pattern "/products/{category}/{id}" and path "/products/electronics/12345",
// it returns map[string]string{"category": "electronics", "id": "12345"}
func ExtractPathParamValues(pattern string, actualPath string) (map[string]string, error) {
	// Extract parameter names from the pattern
	paramNames := ExtractPathParams(pattern)

	// Create a regex pattern by replacing {param} with capture groups
	regexPattern := pattern
	for _, param := range paramNames {
		// Replace {param} with a capture group ([^/]+)
		regexPattern = strings.Replace(regexPattern, "{"+param+"}", "([^/]+)", 1)
	}

	// Escape any regex special characters in the pattern, except the capture groups we just added
	regexPattern = regexp.QuoteMeta(regexPattern)
	// Restore the capture groups which were escaped by QuoteMeta
	for i := 0; i < len(paramNames); i++ {
		regexPattern = strings.Replace(regexPattern, "\\(\\[\\^/\\]\\+\\)", "([^/]+)", 1)
	}

	// Anchor the pattern to match the entire string
	regexPattern = "^" + regexPattern + "$"

	// Compile the regex
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %v", err)
	}

	// Find matches in the actual path
	matches := re.FindStringSubmatch(actualPath)
	if matches == nil {
		return nil, fmt.Errorf("path does not match pattern")
	}

	// Create a map of parameter names to values
	result := make(map[string]string)
	for i, paramName := range paramNames {
		// The first capture group is at index 1
		if i+1 < len(matches) {
			result[paramName] = matches[i+1]
		}
	}

	return result, nil
}
