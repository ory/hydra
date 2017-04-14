package jwk

import "github.com/square/go-jose"

type KeyGenerator interface {
	Generate(id string) (*jose.JsonWebKeySet, error)
}
