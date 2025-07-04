// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package templatex

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexCompiler(t *testing.T) {
	for k, c := range []struct {
		template       string
		delimiterStart byte
		delimiterEnd   byte
		failCompile    bool
		matchAgainst   string
		failMatch      bool
	}{
		{"urn:foo:{.*}", '{', '}', false, "urn:foo:bar:baz", false},
		{"urn:foo.bar.com:{.*}", '{', '}', false, "urn:foo.bar.com:bar:baz", false},
		{"urn:foo.bar.com:{.*}", '{', '}', false, "urn:foo.com:bar:baz", true},
		{"urn:foo.bar.com:{.*}", '{', '}', false, "foobar", true},
		{"urn:foo.bar.com:{.{1,2}}", '{', '}', false, "urn:foo.bar.com:aa", false},

		{"urn:foo.bar.com:{.*{}", '{', '}', true, "", true},
		{"urn:foo:<.*>", '<', '>', false, "urn:foo:bar:baz", false},

		// Ignoring this case for now...
		//{"urn:foo.bar.com:{.*\\{}", '{', '}', false, "", true},
	} {
		k++
		result, err := CompileRegex(c.template, c.delimiterStart, c.delimiterEnd)
		assert.Equal(t, c.failCompile, err != nil, "Case %d", k)
		if c.failCompile || err != nil {
			continue
		}

		t.Logf("Case %d compiled to: %s", k, result.String())
		ok, err := regexp.MatchString(result.String(), c.matchAgainst)
		assert.Nil(t, err, "Case %d", k)
		assert.Equal(t, !c.failMatch, ok, "Case %d", k)
	}
}
