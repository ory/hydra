// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package urlx

import (
	"context"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
)

// hostCache caches DNS lookup results for hostnames to avoid repeated lookups.
// The cache is thread-safe and stores true/false whether a hostname resolves to public IPs.
type hostCache struct {
	mu    sync.RWMutex
	cache map[string]bool
}

// get retrieves a cached value for a hostname. Returns value and whether it was found.
func (hc *hostCache) get(hostname string) (bool, bool) {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	isPublic, found := hc.cache[hostname]
	return isPublic, found
}

// set stores the lookup result for a hostname.
func (hc *hostCache) set(hostname string, isPublic bool) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.cache[hostname] = isPublic
}

// localCache lives for the lifetime of the main process. The cache
// size is not expected to grow more than a few hundred bytes.
var localCache = &hostCache{
	cache: make(map[string]bool),
}

// ExtractPublicAddress iterates over parameters and extracts the first public
// address found. Parameter values are assumed to be in priority order. Returns
// an empty string if only private addresses are available.
func ExtractPublicAddress(values ...string) string {
	for _, value := range values {
		if value == "" || value == "*" {
			continue
		}
		host := value

		// parse URL addresses
		if u, err := url.Parse(value); err == nil && len(u.Host) > 1 {
			host = removeWildcardsFromHostname(u.Host)
		}

		// strip port on both URL and non-URL addresses
		hostname, _, err := net.SplitHostPort(host)
		if err != nil {
			hostname = host
		}

		// for IP addresses
		if ip := net.ParseIP(hostname); ip != nil {
			if !isPrivateIP(ip) {
				return host
			}
			continue
		}

		// for hostnames, first check cache
		if isPublic, found := localCache.get(hostname); found {
			if isPublic {
				return host
			}
			continue
		}

		// otherwise, perform DNS lookup & cache result
		isPublic := isPublicHostname(hostname)
		localCache.set(hostname, isPublic)
		if isPublic {
			return host
		}
	}

	return ""
}

// isPrivateIP checks if an IP address is private (RFC 1918/4193).
func isPrivateIP(ip net.IP) bool {
	return ip.IsPrivate() ||
		ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified() // 0.0.0.0 or ::
}

// isPublicHostname performs DNS lookup to determine if hostname resolves to public IPs.
// Returns true if at least one resolved IP is public, false if all are private or lookup fails.
func isPublicHostname(hostname string) bool {
	// avoid DNS lookup if localhost
	lower := strings.ToLower(hostname)
	if lower == "localhost" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ips, err := net.DefaultResolver.LookupIPAddr(ctx, hostname)
	if err != nil {
		return false
	}

	for _, ip := range ips {
		if !isPrivateIP(ip.IP) {
			return true
		}
	}

	return false
}

// removeWildcardsFromHostname removes wildcard segments from a hostname string
// by splitting on dots and filtering out asterisk-only segments.
func removeWildcardsFromHostname(hostname string) string {
	sep := strings.Split(hostname, ".")
	clean := make([]string, 0, len(sep))
	for _, s := range sep {
		if s != "*" && s != "" {
			clean = append(clean, s)
		}
	}
	return strings.Join(clean, ".")
}
