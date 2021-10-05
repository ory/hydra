package jwk

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
)

type EdDSAGenerator struct{}

func (g *EdDSAGenerator) Generate(id, use string) (*jose.JSONWebKeySet, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	}

	if id == "" {
		id = uuid.New()
	}

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Key:          privateKey,
				Use:          use,
				KeyID:        Ider("private", id),
				Certificates: []*x509.Certificate{},
			},
			{
				Key:          publicKey,
				Use:          use,
				KeyID:        Ider("public", id),
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}
