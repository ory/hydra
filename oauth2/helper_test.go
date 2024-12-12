// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"testing"

	"github.com/oleiade/reflections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/hmac"
)

func Tokens(c fosite.Configurator, length int) (res [][]string) {
	s := oauth2.NewHMACSHAStrategy(&hmac.HMACStrategy{Config: c}, c)

	for i := 0; i < length; i++ {
		tok, sig, _ := s.Enigma.Generate(context.Background())
		res = append(res, []string{sig, tok})
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

func AssertObjectKeysNotEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		assert.NotEqual(t, c, d, "%s", k)
	}
}

func RequireObjectKeysEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		require.Equal(t, c, d, "%s", k)
	}
}
func RequireObjectKeysNotEqual(t *testing.T, a, b interface{}, keys ...string) {
	assert.True(t, len(keys) > 0, "No keys provided.")
	for _, k := range keys {
		c, err := reflections.GetField(a, k)
		assert.Nil(t, err)
		d, err := reflections.GetField(b, k)
		assert.Nil(t, err)
		require.NotEqual(t, c, d, "%s", k)
	}
}

func TestAssertObjectsAreEqualByKeys(t *testing.T) {
	type foo struct {
		Name string
		Body int
	}
	a := &foo{"foo", 1}
	b := &foo{"bar", 1}
	c := &foo{"baz", 3}

	AssertObjectKeysEqual(t, a, a, "Name", "Body")
	AssertObjectKeysNotEqual(t, a, b, "Name")
	AssertObjectKeysNotEqual(t, a, c, "Name", "Body")
}
