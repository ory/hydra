// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	_ "embed"
	"net/http"

	"github.com/ory/x/corsx"
	"github.com/rs/cors"
)

const CORSConfigSchemaID = "ory://cors-config"

//go:embed cors.schema.json
var CORSConfigSchema []byte

func (p *Provider) CORS(prefix string, defaults cors.Options) (cors.Options, bool) {
	prefix = cleanPrefix(prefix)
	allowedOrigins := p.StringsF(prefix+"cors.allowed_origins", defaults.AllowedOrigins)
	allowInsecureOrigins := p.BoolF("feature_flags.legacy_allow_insecure_origins", false)
	return cors.Options{
		// Populated even though rs/cors ignores it once AllowOriginVaryRequestFunc
		// is set: some consumers read AllowedOrigins directly — notably hydra's
		// oauth2cors.Middleware (which builds its own matcher and treats an empty
		// list as "allow all", a CORS bypass) and oathkeeper's address detection.
		AllowedOrigins:     p.StringsF(prefix+"cors.allowed_origins", defaults.AllowedOrigins),
		AllowedMethods:     p.StringsF(prefix+"cors.allowed_methods", defaults.AllowedMethods),
		AllowedHeaders:     p.StringsF(prefix+"cors.allowed_headers", defaults.AllowedHeaders),
		ExposedHeaders:     p.StringsF(prefix+"cors.exposed_headers", defaults.ExposedHeaders),
		AllowCredentials:   p.BoolF(prefix+"cors.allow_credentials", defaults.AllowCredentials),
		OptionsPassthrough: p.BoolF(prefix+"cors.options_passthrough", defaults.OptionsPassthrough),
		MaxAge:             p.IntF(prefix+"cors.max_age", defaults.MaxAge),
		Debug:              p.BoolF(prefix+"cors.debug", defaults.Debug),
		AllowOriginVaryRequestFunc: func(_ *http.Request, origin string) (bool, []string) {
			return corsx.CheckOrigin(allowedOrigins, origin, allowInsecureOrigins), nil
		},
	}, p.Bool(prefix + "cors.enabled")
}
