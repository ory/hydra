package jwk

import "github.com/square/go-jose"


type Manager interface {
	AddKeys(set string, keys *jose.JsonWebKeySet) error

	AddKey(set string, key *jose.JsonWebKey) error

	RemoveKey(set, kid string) error

	GetKey(set, kid string) (error, *jose.JsonWebKey)

	GetKeys(set string) (error, *jose.JsonWebKeySet)
}
