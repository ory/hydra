package jwk

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"

	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

type ECDSA256Generator struct{}

func (g *ECDSA256Generator) Generate(id string) (*jose.JSONWebKeySet, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	}

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Key:          key,
				KeyID:        ider("private", id),
				Certificates: []*x509.Certificate{},
			},
			{
				Key:          &key.PublicKey,
				KeyID:        ider("public", id),
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}
