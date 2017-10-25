package jwk

import (
	"crypto/rand"
	"crypto/x509"
	"io"

	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

type HS256Generator struct{}

func (g *HS256Generator) Generate(id string) (*jose.JSONWebKeySet, error) {
	// Taken from NewHMACKey
	key := &[32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var sliceKey = key[:]

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Algorithm:    "HS256",
				Key:          sliceKey,
				KeyID:        id,
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}
