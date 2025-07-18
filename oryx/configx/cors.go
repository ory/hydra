// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	_ "embed"

	"github.com/rs/cors"
)

const CORSConfigSchemaID = "ory://cors-config"

//go:embed cors.schema.json
var CORSConfigSchema []byte

func (p *Provider) CORS(prefix string, defaults cors.Options) (cors.Options, bool) {
	prefix = cleanPrefix(prefix)

	return cors.Options{
		AllowedOrigins:     p.StringsF(prefix+"cors.allowed_origins", defaults.AllowedOrigins),
		AllowedMethods:     p.StringsF(prefix+"cors.allowed_methods", defaults.AllowedMethods),
		AllowedHeaders:     p.StringsF(prefix+"cors.allowed_headers", defaults.AllowedHeaders),
		ExposedHeaders:     p.StringsF(prefix+"cors.exposed_headers", defaults.ExposedHeaders),
		AllowCredentials:   p.BoolF(prefix+"cors.allow_credentials", defaults.AllowCredentials),
		OptionsPassthrough: p.BoolF(prefix+"cors.options_passthrough", defaults.OptionsPassthrough),
		MaxAge:             p.IntF(prefix+"cors.max_age", defaults.MaxAge),
		Debug:              p.BoolF(prefix+"cors.debug", defaults.Debug),
	}, p.Bool(prefix + "cors.enabled")
}
