package pkg

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/ory-am/fosite/enigma/hmac"
	"github.com/ory-am/fosite/fosite-example/store"
	"github.com/ory-am/fosite/handler/core/strategy"
	"github.com/ory-am/ladon"
	"github.com/ory-am/ladon/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RequireError(t *testing.T, expectError bool, err error, args ...interface{}) {
	require.Equal(t, expectError, err != nil)
	if err != nil && !expectError {
		t.Logf("Unexpected error: %s\n", err.Error())
		t.Logf("Args error: %v\n", args)
		if e, ok := err.(*errors.Error); !ok {
			t.Logf("Stack trace: %s\n", e.ErrorStack())
		}
		t.Logf("\n\n")
	}
}

func AssertError(t *testing.T, expectError bool, err error, args ...interface{}) {
	assert.Equal(t, expectError, err != nil)
	if err != nil && !expectError {
		t.Logf("Unexpected error: %s\n", err.Error())
		t.Logf("Args error: %s\n", args)
		if e, ok := err.(*errors.Error); ok {
			t.Logf("Stack trace: %s\n", e.ErrorStack())
		}
		t.Logf("\n\n")
	}
}

func LadonWarden(ps map[string]ladon.Policy) ladon.Warden {
	return &ladon.Ladon{
		Manager: &memory.Manager{
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
	Enigma: &hmac.Enigma{
		GlobalSecret: []byte("1234567890123456789012345678901234567890"),
	},
}
