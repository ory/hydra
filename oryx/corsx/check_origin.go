// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package corsx

import "strings"

// CheckOrigin is a function that can be used well with cors.Options.AllowOriginRequestFunc.
// It checks whether the origin is allowed following the same behavior as github.com/rs/cors.
//
// Recommended usage for hot-reloadable origins:
//
//	func (p *Config) cors(ctx context.Context, prefix string) (cors.Options, bool) {
//		opts, enabled := p.GetProvider(ctx).CORS(prefix, cors.Options{
//			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
//			AllowedHeaders:   []string{"Authorization", "Content-Type", "Cookie"},
//			ExposedHeaders:   []string{"Content-Type", "Set-Cookie"},
//			AllowCredentials: true,
//		})
//		opts.AllowOriginRequestFunc = func(r *http.Request, origin string) bool {
//			// load the origins from the config on every request to allow hot-reloading
//			allowedOrigins := p.GetProvider(r.Context()).Strings(prefix + ".cors.allowed_origins")
//			return corsx.CheckOrigin(allowedOrigins, origin)
//		}
//		return opts, enabled
//	}
func CheckOrigin(allowedOrigins []string, origin string) bool {
	if len(allowedOrigins) == 0 {
		return true
	}
	for _, o := range allowedOrigins {
		if o == "*" {
			// allow all origins
			return true
		}
		// Note: for origins and methods matching, the spec requires a case-sensitive matching.
		// As it may be error-prone, we chose to ignore the spec here.
		// https://github.com/rs/cors/blob/066574eebbd0f5f1b6cd1154a160cc292ac1835e/cors.go#L132-L133
		o = strings.ToLower(o)
		prefix, suffix, found := strings.Cut(o, "*")
		if !found {
			// not a pattern, check for equality
			if o == origin {
				return true
			}
			continue
		}
		// inspired by https://github.com/rs/cors/blob/066574eebbd0f5f1b6cd1154a160cc292ac1835e/utils.go#L15
		if len(origin) >= len(prefix)+len(suffix) && strings.HasPrefix(origin, prefix) && strings.HasSuffix(origin, suffix) {
			return true
		}
	}
	return false
}
