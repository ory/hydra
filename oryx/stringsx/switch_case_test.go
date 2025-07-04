// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package stringsx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisteredCases(t *testing.T) {
	t.Run("case=adds values", func(t *testing.T) {
		v1, v2 := "value 1", "value 2"

		e := RegisteredCases{}
		e.AddCase(v1)
		e.AddCase(v2)

		p := RegisteredPrefixes{}
		p.HasPrefix(v1)
		p.HasPrefix(v2)

		assert.Equal(t, []string{v1, v2}, e.cases)
		assert.Equal(t, []string{v1, v2}, p.prefixes)
	})

	t.Run("case=returns equality on add", func(t *testing.T) {
		v1, v2 := "value 1", "value 2"

		cs := SwitchExact(v1)
		assert.True(t, cs.AddCase(v1))
		assert.False(t, cs.AddCase(v2))
	})

	t.Run("case=converts to correct error", func(t *testing.T) {
		c1, c2, actual := "case 1", "case 2", "actual"

		e := SwitchExact(actual)
		p := SwitchPrefix(actual)
		e.AddCase(c1)
		p.HasPrefix(c1)
		e.AddCase(c2)
		p.HasPrefix(c2)

		ee := e.ToUnknownCaseErr()
		pe := p.ToUnknownPrefixErr()

		assert.True(t, errors.Is(ee, ErrUnknownCase))
		assert.True(t, errors.Is(pe, ErrUnknownPrefix))

		for _, v := range []string{c1, c2, actual} {
			assert.Contains(t, ee.Error(), v)
			assert.Contains(t, pe.Error(), v)
		}
	})

	t.Run("case=switch integration", func(t *testing.T) {
		var err error

		switch f := SwitchExact("foo"); {
		case f.AddCase("bar"):
			t.FailNow()
		case f.AddCase("baz"):
			t.FailNow()
		default:
			err = f.ToUnknownCaseErr()
		}

		assert.True(t, errors.Is(err, ErrUnknownCase))

		switch p := SwitchPrefix("foobarbaz"); {
		case p.HasPrefix("foobaz"):
			t.FailNow()
		case p.HasPrefix("unknown"):
			t.FailNow()
		default:
			err = p.ToUnknownPrefixErr()
		}

		assert.True(t, errors.Is(err, ErrUnknownPrefix))
	})
}
