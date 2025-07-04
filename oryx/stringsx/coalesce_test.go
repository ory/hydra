// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringsx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoalesce(t *testing.T) {
	for k, tc := range []struct {
		in     []string
		expect string
	}{
		{
			in:     []string{"", "", "foo"},
			expect: "foo",
		},
		{
			in:     []string{"bar", "", "foo"},
			expect: "bar",
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			assert.EqualValues(t, tc.expect, Coalesce(tc.in...))
		})
	}
}
