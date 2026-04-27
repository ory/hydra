// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package ipx

import (
	"context"
	stderrors "errors"
	"net"
	"net/netip"
	"net/url"
	"time"

	"code.dny.dev/ssrf"
	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
)

// IsAssociatedIPAllowedWhenSet is a wrapper for IsAssociatedIPAllowed which returns valid
// when ipOrHostnameOrURL is empty.
func IsAssociatedIPAllowedWhenSet(ctx context.Context, ipOrHostnameOrURL string) error {
	if ipOrHostnameOrURL == "" {
		return nil
	}
	return IsAssociatedIPAllowed(ctx, ipOrHostnameOrURL)
}

// AreAllAssociatedIPsAllowed fails if one of the pairs is failing.
func AreAllAssociatedIPsAllowed(ctx context.Context, pairs map[string]string) error {
	g, ctx := errgroup.WithContext(ctx)
	for key, ipOrHostnameOrURL := range pairs {
		g.Go(func() error {
			return errors.Wrapf(IsAssociatedIPAllowed(ctx, ipOrHostnameOrURL), "key %s validation is failing", key)
		})
	}
	return g.Wait()
}

// IsAssociatedIPAllowed returns nil for a domain (with NS lookup), IP, or IPv6 address if it
// does not resolve to a private IP subnet. This is a first level of defense against
// SSRF attacks by disallowing any domain or IP to resolve to a private network range.
//
// Please keep in mind that validations for domains is valid only when looking up.
// A malicious actor could easily update the DSN record post validation to point
// to an internal IP
func IsAssociatedIPAllowed(ctx context.Context, ipOrHostnameOrURL string) error {
	ipOrHostname := ipOrHostnameOrURL
	if parsed, err := url.ParseRequestURI(ipOrHostnameOrURL); err == nil {
		ipOrHostname = parsed.Hostname()
	}

	if ip, err := netip.ParseAddr(ipOrHostname); err == nil {
		if !allowed(ip) {
			return errors.Errorf("ip %s is not a permitted destination", ip)
		}
		return nil
	}

	if addr, err := netip.ParseAddrPort(ipOrHostnameOrURL); err == nil {
		if !allowed(addr.Addr()) {
			return errors.Errorf("ip %s is not a permitted destination", addr.Addr())
		}
		return nil
	}

	ctx, cancel := context.WithTimeoutCause(ctx, 2*time.Second, errors.New("DNS lookup timed out"))
	defer cancel()
	ips, err := resolver.LookupNetIP(ctx, "ip", ipOrHostname)
	if err != nil {
		if dnsErr, ok := stderrors.AsType[*net.DNSError](err); ok {
			// Copy the `*net.DNSError` before masking `Server` to avoid a data
			// race: the DNS resolver uses `singleflight` to deduplicate
			// concurrent lookups, so multiple goroutines may receive the same
			// `*net.DNSError` pointer. Mutating it in place races with concurrent
			// readers (e.g. the `otelhttp` `dnsDone` trace hook).
			maskedDNS := *dnsErr
			maskedDNS.Server = "" // Mask our DNS server's IP address.
			return errors.Wrapf(&maskedDNS, "failed to resolve %s", ipOrHostnameOrURL)
		}
	}

	for _, ip := range ips {
		if !allowed(ip) {
			return errors.Wrapf(&net.DNSError{
				Err:         "no such host",
				Name:        ipOrHostname,
				Server:      "",
				IsTimeout:   false,
				IsTemporary: false,
				IsNotFound:  true,
			}, "failed to resolve %s", ipOrHostnameOrURL)
		}
	}

	return nil
}

var resolver = &net.Resolver{PreferGo: true}

func allowed(ip netip.Addr) bool {
	if !ip.IsGlobalUnicast() {
		return false
	}

	if ip.Is4() {
		for _, net := range ssrf.IPv4DeniedPrefixes {
			if net.Contains(ip) {
				return false
			}
		}
	} else { // IPv6
		// ip.IsGlobalUnicast returns true for IPv6 addresses which fall outside
		// of the current IANA-allocated 2000::/3 global unicast space. Hence,
		// we need to check for ourselves if the IPv6 address is within the
		// currently allocated block.
		if !ssrf.IPv6GlobalUnicast.Contains(ip) {
			return false
		}
		for _, net := range ssrf.IPv6DeniedPrefixes {
			if net.Contains(ip) {
				return false
			}
		}
	}

	return true
}
