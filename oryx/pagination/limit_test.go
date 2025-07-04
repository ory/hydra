// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package pagination

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	for k, c := range []struct {
		s      []string
		offset int
		limit  int
		e      []string
	}{
		{
			s:      []string{"a", "b", "c"},
			offset: 0,
			limit:  100,
			e:      []string{"a", "b", "c"},
		},
		{
			s:      []string{"a", "b", "c"},
			offset: 0,
			limit:  2,
			e:      []string{"a", "b"},
		},
		{
			s:      []string{"a", "b", "c"},
			offset: 1,
			limit:  10,
			e:      []string{"b", "c"},
		},
		{
			s:      []string{"a", "b", "c"},
			offset: 1,
			limit:  2,
			e:      []string{"b", "c"},
		},
		{
			s:      []string{"a", "b", "c"},
			offset: 2,
			limit:  2,
			e:      []string{"c"},
		},
		{
			s:      []string{"a", "b", "c"},
			offset: 3,
			limit:  10,
			e:      []string{},
		},
		{
			s:      []string{"a", "b", "c"},
			offset: 2,
			limit:  10,
			e:      []string{"c"},
		},
		{
			s:      []string{"a", "b", "c"},
			offset: 1,
			limit:  10,
			e:      []string{"b", "c"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			start, end := Index(c.limit, c.offset, len(c.s))
			assert.EqualValues(t, c.e, c.s[start:end])
		})
	}
}
