package jwk

import (
	"crypto/rsa"

	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

func MustRSAPublic(key *jose.JsonWebKey) *rsa.PublicKey {
	res, err := ToRSAPublic(key)
	if err != nil {
		panic(err.Error())
	}
	return res

}

func ToRSAPublic(key *jose.JsonWebKey) (*rsa.PublicKey, error) {
	res, ok := key.Key.(*rsa.PublicKey)
	if !ok {
		return res, errors.New("Could not convert key to RSA Private Key.")
	}
	return res, nil
}

func MustRSAPrivate(key *jose.JsonWebKey) *rsa.PrivateKey {
	res, err := ToRSAPrivate(key)
	if err != nil {
		panic(err.Error())
	}
	return res
}

func ToRSAPrivate(key *jose.JsonWebKey) (*rsa.PrivateKey, error) {
	res, ok := key.Key.(*rsa.PrivateKey)
	if !ok {
		return res, errors.New("Could not convert key to RSA Private Key.")
	}
	return res, nil
}
