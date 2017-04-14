package jwk

import "github.com/square/go-jose"

type Manager interface {
	AddKey(set string, key *jose.JsonWebKey) error

	AddKeySet(set string, keys *jose.JsonWebKeySet) error

	GetKey(set, kid string) (*jose.JsonWebKeySet, error)

	GetKeySet(set string) (*jose.JsonWebKeySet, error)

	DeleteKey(set, kid string) error

	DeleteKeySet(set string) error
}
