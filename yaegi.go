// Package traefik_path_param_headers is a Traefik v3 middleware plugin.
package traefik_path_param_headers

import (
	"context"
	"net/http"
)

// Symbols exposes the plugin constructors for Traefik v3.3.3
var Symbols = struct {
	CreateConfig func() *Config
	New          func(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error)
}{
	CreateConfig: CreateConfig,
	New:          New,
}
