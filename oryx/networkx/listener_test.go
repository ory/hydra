// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package networkx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressIsUnixSocket(t *testing.T) {
	for k, tc := range []struct {
		a string
		e bool
	}{
		{a: "unix:/var/baz", e: true},
		{a: "https://foo", e: false},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.EqualValues(t, tc.e, AddressIsUnixSocket(tc.a))
		})
	}
}
