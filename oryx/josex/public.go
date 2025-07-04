package josex

import (
	"crypto"

	"github.com/go-jose/go-jose/v3"
)

// ToPublicKey returns the public key of the given private key.
func ToPublicKey(k *jose.JSONWebKey) jose.JSONWebKey {
	if key := k.Public(); key.Key != nil {
		return key
	}

	// HSM workaround - jose does not understand crypto.Signer / HSM so we need to manually
	// extract the public key.
	switch key := k.Key.(type) {
	case crypto.Signer:
		newKey := *k
		newKey.Key = key.Public()
		return newKey
	case jose.OpaqueSigner:
		newKey := *k
		newKey.Key = key.Public().Key
		return newKey
	}

	return jose.JSONWebKey{}
}
