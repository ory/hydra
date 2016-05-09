package jwk

import "github.com/square/go-jose"

type KeyGenerator interface {
	Generate(id string) (set *jose.JsonWebKeySet, error)
}
