package jwk

import (
	"testing"

	"github.com/square/go-jose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindKeyByPrefix(t *testing.T) {
	jwks := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{
		{KeyID: "public:foo"},
		{KeyID: "private:foo"},
	}}

	key, err := FindKeyByPrefix(jwks, "public")
	require.NoError(t, err)
	assert.Equal(t, "public:foo", key.KeyID)

	key, err = FindKeyByPrefix(jwks, "private")
	require.NoError(t, err)
	assert.Equal(t, "private:foo", key.KeyID)

	_, err = FindKeyByPrefix(jwks, "asdf")
	require.Error(t, err)

	jwks = &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{
		{KeyID: "public:"},
		{KeyID: "private:"},
	}}

	key, err = FindKeyByPrefix(jwks, "public")
	require.NoError(t, err)
	assert.Equal(t, "public:", key.KeyID)

	key, err = FindKeyByPrefix(jwks, "private")
	require.NoError(t, err)
	assert.Equal(t, "private:", key.KeyID)

	_, err = FindKeyByPrefix(jwks, "asdf")
	require.Error(t, err)

	jwks = &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{
		{KeyID: ""},
	}}
	require.Error(t, err)
}

func TestIder(t *testing.T) {
	assert.True(t, len(ider("public", "")) > len("public:"))
	assert.Equal(t, "public:foo", ider("public", "foo"))
}
