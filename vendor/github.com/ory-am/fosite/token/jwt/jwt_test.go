package jwt

import (
	"strings"
	"testing"

	"time"

	"github.com/ory-am/fosite/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var header = &Headers{
	Extra: map[string]interface{}{
		"foo": "bar",
	},
}

func TestHash(t *testing.T) {
	j := RS256JWTStrategy{
		PrivateKey: internal.MustRSAKey(),
	}
	in := []byte("foo")
	out, err := j.Hash(in)
	assert.Nil(t, err)
	assert.NotEqual(t, in, out)
}

func TestAssign(t *testing.T) {
	for k, c := range [][]map[string]interface{}{
		{
			{"foo": "bar"},
			{"baz": "bar"},
			{"foo": "bar", "baz": "bar"},
		},
		{
			{"foo": "bar"},
			{"foo": "baz"},
			{"foo": "bar"},
		},
		{
			{},
			{"foo": "baz"},
			{"foo": "baz"},
		},
		{
			{"foo": "bar"},
			{"foo": "baz", "bar": "baz"},
			{"foo": "bar", "bar": "baz"},
		},
	} {
		assert.EqualValues(t, c[2], assign(c[0], c[1]), "Case %d", k)
	}
}

func TestGenerateJWT(t *testing.T) {
	claims := &JWTClaims{
		ExpiresAt: time.Now().Add(time.Hour),
	}

	j := RS256JWTStrategy{
		PrivateKey: internal.MustRSAKey(),
	}

	token, sig, err := j.Generate(claims.ToMapClaims(), header)
	require.Nil(t, err, "%s", err)
	require.NotNil(t, token)

	sig, err = j.Validate(token)
	require.Nil(t, err, "%s", err)

	sig, err = j.Validate(token + "." + "0123456789")
	require.NotNil(t, err, "%s", err)

	partToken := strings.Split(token, ".")[2]

	sig, err = j.Validate(partToken)
	require.NotNil(t, err, "%s", err)

	// Reset private key
	j.PrivateKey = internal.MustRSAKey()

	// Lets validate the exp claim
	claims = &JWTClaims{
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	token, sig, err = j.Generate(claims.ToMapClaims(), header)
	require.Nil(t, err, "%s", err)
	require.NotNil(t, token)
	//t.Logf("%s.%s", token, sig)

	sig, err = j.Validate(token)
	require.NotNil(t, err, "%s", err)

	// Lets validate the nbf claim
	claims = &JWTClaims{
		NotBefore: time.Now().Add(time.Hour),
	}
	token, sig, err = j.Generate(claims.ToMapClaims(), header)
	require.Nil(t, err, "%s", err)
	require.NotNil(t, token)
	//t.Logf("%s.%s", token, sig)
	sig, err = j.Validate(token)
	require.NotNil(t, err, "%s", err)
	require.Empty(t, sig, "%s", err)
}

func TestValidateSignatureRejectsJWT(t *testing.T) {
	var err error
	j := RS256JWTStrategy{
		PrivateKey: internal.MustRSAKey(),
	}

	for k, c := range []string{
		"",
		" ",
		"foo.bar",
		"foo.",
		".foo",
	} {
		_, err = j.Validate(c)
		assert.NotNil(t, err, "%s", err)
		t.Logf("Passed test case %d", k)
	}
}
