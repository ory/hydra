// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"fmt"
	"net"
	"net/netip"
	"net/url"

	"code.dny.dev/ssrf"
	"github.com/pkg/errors"
)

// ErrPrivateIPAddressDisallowed is returned when a private IP address is disallowed.
type ErrPrivateIPAddressDisallowed error

// DisallowIPPrivateAddresses returns nil for a domain (with NS lookup), IP, or IPv6 address if it
// does not resolve to a private IP subnet. This is a first level of defense against
// SSRF attacks by disallowing any domain or IP to resolve to a private network range.
//
// Please keep in mind that validations for domains is valid only when looking up.
// A malicious actor could easily update the DSN record post validation to point
// to an internal IP
func DisallowIPPrivateAddresses(ipOrHostnameOrURL string) error {
	lookup := func(hostname string) ([]net.IP, error) {
		lookup, err := net.LookupIP(hostname)
		if err != nil {
			if dnsErr := new(net.DNSError); errors.As(err, &dnsErr) && (dnsErr.IsNotFound || dnsErr.IsTemporary) {
				// If the hostname does not resolve, we can't validate it. So yeah,
				// I guess we're allowing it.
				return nil, nil
			}
			return nil, errors.WithStack(err)
		}
		return lookup, nil
	}

	var ips []net.IP
	ip := net.ParseIP(ipOrHostnameOrURL)
	if ip == nil {
		if result, err := lookup(ipOrHostnameOrURL); err != nil {
			return err
		} else if result != nil {
			ips = append(ips, result...)
		}

		if parsed, err := url.Parse(ipOrHostnameOrURL); err == nil {
			if result, err := lookup(parsed.Hostname()); err != nil {
				return err
			} else if result != nil {
				ips = append(ips, result...)
			}
		}
	} else {
		ips = append(ips, ip)
	}

	for _, ip := range ips {
		ip, err := netip.ParseAddr(ip.String())
		if err != nil {
			return ErrPrivateIPAddressDisallowed(errors.WithStack(err)) // should be unreacheable
		}

		if ip.Is4() {
			for _, deny := range ssrf.IPv4DeniedPrefixes {
				if deny.Contains(ip) {
					return ErrPrivateIPAddressDisallowed(fmt.Errorf("%s is not a public IP address", ip))
				}
			}
		} else {
			if !ssrf.IPv6GlobalUnicast.Contains(ip) {
				return ErrPrivateIPAddressDisallowed(fmt.Errorf("%s is not a public IP address", ip))
			}
			for _, net := range ssrf.IPv6DeniedPrefixes {
				if net.Contains(ip) {
					return ErrPrivateIPAddressDisallowed(fmt.Errorf("%s is not a public IP address", ip))
				}
			}
		}
	}

	return nil
}
