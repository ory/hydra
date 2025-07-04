// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package corsx

import "net/url"

// NormalizeOrigins normalizes the CORS origins.
func NormalizeOrigins(origins []url.URL) []string {
	results := make([]string, len(origins))
	for k, o := range origins {
		results[k] = o.Scheme + "://" + o.Host
	}
	return results
}

// NormalizeOriginStrings normalizes the CORS origins from string representation
func NormalizeOriginStrings(origins []string) ([]string, error) {
	results := make([]string, len(origins))
	for k, o := range origins {
		u, err := url.ParseRequestURI(o)
		if err != nil {
			return nil, err
		}
		results[k] = u.Scheme + "://" + u.Host
	}
	return results, nil
}
