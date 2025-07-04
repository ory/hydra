// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package ipx

import (
	"context"
	"net"
	"net/url"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
)

// IsAssociatedIPAllowedWhenSet is a wrapper for IsAssociatedIPAllowed which returns valid
// when ipOrHostnameOrURL is empty.
func IsAssociatedIPAllowedWhenSet(ipOrHostnameOrURL string) error {
	if ipOrHostnameOrURL == "" {
		return nil
	}
	return IsAssociatedIPAllowed(ipOrHostnameOrURL)
}

// AreAllAssociatedIPsAllowed fails if one of the pairs is failing.
func AreAllAssociatedIPsAllowed(pairs map[string]string) error {
	g := new(errgroup.Group)
	for key, ipOrHostnameOrURL := range pairs {
		key := key
		ipOrHostnameOrURL := ipOrHostnameOrURL
		g.Go(func() error {
			return errors.Wrapf(IsAssociatedIPAllowed(ipOrHostnameOrURL), "key %s validation is failing", key)
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
func IsAssociatedIPAllowed(ipOrHostnameOrURL string) error {
	lookup := func(hostname string) []net.IP {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		lookup, err := net.DefaultResolver.LookupIPAddr(ctx, hostname)
		if err != nil {
			return nil
		}
		ips := make([]net.IP, len(lookup))
		for i, ip := range lookup {
			ips[i] = ip.IP
		}
		return ips
	}

	var ips []net.IP
	ip := net.ParseIP(ipOrHostnameOrURL)
	if ip == nil {
		if result := lookup(ipOrHostnameOrURL); result != nil {
			ips = append(ips, result...)
		}

		if parsed, err := url.Parse(ipOrHostnameOrURL); err == nil {
			if result := lookup(parsed.Hostname()); result != nil {
				ips = append(ips, result...)
			}
		}
	} else {
		ips = append(ips, ip)
	}

	for _, disabled := range []string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fd47:1ed0:805d:59f0::/64",
		"fc00::/7",
		"::1/128",
	} {
		_, cidr, err := net.ParseCIDR(disabled)
		if err != nil {
			return err
		}

		for _, ip := range ips {
			if cidr.Contains(ip) {
				return errors.Errorf("ip %s is in the %s range", ip, disabled)
			}
		}
	}

	return nil
}
