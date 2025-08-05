// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package ipx

import (
	"iter"
	"net/netip"
)

func Hosts(prefix netip.Prefix) iter.Seq[netip.Addr] {
	prefix = prefix.Masked()
	return func(yield func(netip.Addr) bool) {
		if !prefix.IsValid() {
			return
		}
		for addr := prefix.Addr().Next(); prefix.Contains(addr); addr = addr.Next() {
			if !yield(addr) {
				return
			}
		}
	}
}
