// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringslice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	for i, tc := range []struct {
		i, e []string
	}{
		{
			i: []string{"a", "b", "c"},
			e: []string{"c", "b", "a"},
		},
		{
			i: []string{"foo"},
			e: []string{"foo"},
		},
		{
			i: []string{"foo", "bar"},
			e: []string{"bar", "foo"},
		},
		{
			i: []string{},
			e: []string{},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/input:%v expected:%v", i, tc.i, tc.e), func(t *testing.T) {
			assert.Equal(t, tc.e, Reverse(tc.i))
		})
	}
}
