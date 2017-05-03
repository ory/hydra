package pkg

import (
	"testing"
	"time"

	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/storage"
	"github.com/ory/fosite/token/hmac"
	"github.com/ory/ladon"
	"github.com/ory/ladon/manager/memory"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var HMACStrategy = &oauth2.HMACSHAStrategy{
	Enigma: &hmac.HMACStrategy{
		GlobalSecret: []byte("1234567890123456789012345678901234567890"),
	},
	AccessTokenLifespan:   time.Hour,
	AuthorizeCodeLifespan: time.Hour,
}

func RequireError(t *testing.T, expectError bool, err error, args ...interface{}) {
	if err != nil && !expectError {
		t.Logf("Unexpected error: %s\n", err.Error())
		t.Logf("Arguments: %v\n", args)
		if e, ok := errors.Cause(err).(stackTracer); ok {
			t.Logf("Stack:\n%+v\n", e.StackTrace())
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
		if e, ok := errors.Cause(err).(stackTracer); ok {
			t.Logf("Stack:\n%+v\n", e.StackTrace())
		}
		t.Logf("\n\n")
	}
}

func LadonWarden(ps map[string]ladon.Policy) ladon.Warden {
	return &ladon.Ladon{
		Manager: &memory.MemoryManager{
			Policies: ps,
		},
	}
}

func FositeStore() *storage.MemoryStore {
	return storage.NewMemoryStore()
}

func Tokens(length int) (res [][]string) {
	for i := 0; i < length; i++ {
		tok, sig, _ := HMACStrategy.Enigma.Generate()
		res = append(res, []string{sig, tok})
	}
	return res
}
