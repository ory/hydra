// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"testing"

	"github.com/oleiade/reflections"
	"github.com/stretchr/testify/assert"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/hmac"
)

func Tokens(c fosite.Configurator, length int) []struct{ sig, tok string } {
	s := oauth2.NewHMACSHAStrategy(&hmac.HMACStrategy{Config: c}, c)

	res := make([]struct{ sig, tok string }, length)
	for i := range res {
		res[i].tok, res[i].sig, _ = s.Enigma.Generate(context.Background())
	}
	return res
}

func AssertObjectKeysEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		assert.Equal(t, c, d, "%s", k)
	}
}

func TestAssertObjectsAreEqualByKeys(t *testing.T) {
	type foo struct {
		Name string
		Body int
	}
	a := &foo{"foo", 1}

	AssertObjectKeysEqual(t, a, a, "Name", "Body")
}
