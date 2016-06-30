package pkg

import (
	"testing"

	"time"

	"github.com/go-errors/errors"
	"github.com/ory-am/fosite/fosite-example/store"
	"github.com/ory-am/fosite/handler/core/strategy"
	"github.com/ory-am/fosite/token/hmac"
	"github.com/ory-am/ladon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RequireError(t *testing.T, expectError bool, err error, args ...interface{}) {
	if err != nil && !expectError {
		t.Logf("Unexpected error: %s\n", err.Error())
		t.Logf("Arguments: %v\n", args)
		if e, ok := err.(*errors.Error); ok {
			t.Logf("Stack:\n%s\n", e.ErrorStack())
		}
		t.Logf("\n\n")
	}
	require.Equal(t, expectError, err != nil, "%v", args)
}

func AssertError(t *testing.T, expectError bool, err error, args ...interface{}) {
	assert.Equal(t, expectError, err != nil, "%v", args)
	if err != nil && !expectError {
		t.Logf("Unexpected error: %s\n", err.Error())
		t.Logf("Arguments: %s\n", args)
		if e, ok := err.(*errors.Error); ok {
			t.Logf("Stack:\n%s\n", e.ErrorStack())
		}
		t.Logf("\n\n")
	}
}

func LadonWarden(ps map[string]ladon.Policy) ladon.Warden {
	return &ladon.Ladon{
		Manager: &ladon.MemoryManager{
			Policies: ps,
		},
	}
}

func FositeStore() *store.Store {
	return store.NewStore()
}

func Tokens(length int) (res [][]string) {
	for i := 0; i < length; i++ {
		tok, sig, _ := HMACStrategy.Enigma.Generate()
		res = append(res, []string{sig, tok})
	}
	return res
}

var HMACStrategy = &strategy.HMACSHAStrategy{
	Enigma: &hmac.HMACStrategy{
		GlobalSecret: []byte("1234567890123456789012345678901234567890"),
	},
	AccessTokenLifespan:   time.Hour,
	AuthorizeCodeLifespan: time.Hour,
}
