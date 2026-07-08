// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package corsx

import (
	"strings"

	"golang.org/x/net/publicsuffix"
)

// CheckOrigin is a function that can be used well with cors.Options.AllowOriginRequestFunc.
// It checks whether the origin is allowed following the same behavior as github.com/rs/cors.
//
// When legacyAllowInsecureOrigins is false (the default), wildcard patterns are
// only honored when ClassifyOrigin reports them as bounded at a registrable
// domain. Pass true to opt into legacy (trusting) behavior for unbounded
// wildcards.
//
// TODO: legacyAllowInsecureOrigins grandfathers a fixed set of projects through
// a time-boxed migration window (feature_flags.legacy_allow_insecure_origins).
// Once those projects move to bounded wildcards and the entitlement is revoked,
// drop this parameter and always enforce the boundary.
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
//			return corsx.CheckOrigin(allowedOrigins, origin, false)
//		}
//		return opts, enabled
//	}
func CheckOrigin(allowedOrigins []string, origin string, legacyAllowInsecureOrigins bool) bool {
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
		// Only honor wildcards bounded at a registrable domain unless the caller
		// explicitly opts into legacy insecure matching. See ClassifyOrigin.
		if !legacyAllowInsecureOrigins && ClassifyOrigin(o).IsUnsafeWildcard() {
			continue
		}
		// inspired by https://github.com/rs/cors/blob/066574eebbd0f5f1b6cd1154a160cc292ac1835e/utils.go#L15
		if len(origin) >= len(prefix)+len(suffix) && strings.HasPrefix(origin, prefix) && strings.HasSuffix(origin, suffix) {
			return true
		}
	}
	return false
}

// OriginPattern describes a CORS origin or return-URL host pattern: whether it
// uses a wildcard and, if so, whether that wildcard is safely bounded at a
// registrable domain.
type OriginPattern struct {
	// HasWildcard reports whether the pattern contains a "*".
	HasWildcard bool

	// BoundedWildcard reports whether the "*" is confined to a subdomain label
	// and the fixed domain that follows it is a registrable domain (an eTLD+1,
	// e.g. "example.com" or "example.co.uk"). Every host the pattern can match
	// then shares that one customer-owned registrable domain, so an attacker
	// cannot register a matching host. Always false when HasWildcard is false.
	BoundedWildcard bool

	// Base is the fixed domain that follows the wildcard label — "example.com"
	// for "*.example.com", "com" for "*.com". It is empty for non-wildcards and
	// for bare or trailing wildcards where no domain follows the "*". When
	// BoundedWildcard is false, Base names the offending suffix, which is the
	// actionable signal for reporting why a wildcard was rejected.
	Base string
}

// IsUnsafeWildcard reports whether the pattern is a wildcard that is NOT bounded
// at a registrable domain. Such a wildcard would match an attacker-registrable
// host (e.g. "https://*foo.com" matches "https://evilfoo.com"), so it must be
// rejected unless the caller explicitly opts into legacy insecure matching. This
// is the dominant question at call sites that gate, drop, or reject wildcard
// origins and return URLs.
func (p OriginPattern) IsUnsafeWildcard() bool {
	return p.HasWildcard && !p.BoundedWildcard
}

// ClassifyOrigin inspects a CORS origin or bare host pattern and reports whether
// it is a wildcard and, if so, whether the wildcard is safely bounded at a
// registrable domain. Only the text from the last "*" onward is inspected, so
// the result is identical whether pattern carries a scheme or is a bare host; a
// trailing ":port" is ignored. Examples:
//
//   - "https://*.example.com" → {HasWildcard: true, BoundedWildcard: true, Base: "example.com"}
//   - "https://*foo.com"      → {HasWildcard: true, Base: "com"}   (dot-less; base is a public suffix)
//   - "https://*.com"         → {HasWildcard: true, Base: "com"}   (public suffix, not registrable)
//   - "https://www.ory.*"     → {HasWildcard: true}                (trailing; no domain follows)
//   - "https://exact.foo.com" → {}                                 (no wildcard)
func ClassifyOrigin(pattern string) OriginPattern {
	i := strings.LastIndexByte(pattern, '*')
	if i < 0 {
		return OriginPattern{}
	}
	p := OriginPattern{HasWildcard: true}
	// The fixed domain is everything after the first "." that follows the last
	// "*" — i.e. the label containing the wildcard is dropped. Without such a dot
	// the wildcard is trailing or bare, so no registrable domain follows.
	_, base, found := strings.Cut(pattern[i:], ".")
	if !found {
		return p
	}
	base = strings.ToLower(base)
	if host, _, found := strings.Cut(base, ":"); found {
		base = host // Drop a trailing ":port".
	}
	p.Base = base
	// EffectiveTLDPlusOne returns an error when base is itself a public suffix
	// (e.g. "com", "co.uk", "vercel.app") or otherwise has no registrable domain.
	_, err := publicsuffix.EffectiveTLDPlusOne(base)
	p.BoundedWildcard = err == nil
	return p
}
