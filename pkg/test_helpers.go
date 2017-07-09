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
