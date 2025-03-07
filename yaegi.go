// Package pathparamheaders is a Traefik v3 middleware plugin.
package pathparamheaders

// Yaegi is required for Traefik v3 plugins.
// It exports the plugin constructors.
var Yaegi = map[string]interface{}{
	"CreateConfig": CreateConfig,
	"New":          New,
}

// Declare exported symbols to avoid dyncheck complaining.
var (
	_ = CreateConfig
	_ = New
)
