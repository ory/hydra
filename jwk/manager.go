package jwk

import "github.com/square/go-jose"

type Manager interface {
	AddKey(set string, key *jose.JSONWebKey) error

	AddKeySet(set string, keys *jose.JSONWebKeySet) error

	GetKey(set, kid string) (*jose.JSONWebKeySet, error)

	GetKeySet(set string) (*jose.JSONWebKeySet, error)

	DeleteKey(set, kid string) error

	DeleteKeySet(set string) error
}
