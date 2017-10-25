package jwk

import (
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
	"io"
	"crypto/rand"
)

type HS256Generator struct {
}

func (g *HS256Generator) Generate(id string) (*jose.JsonWebKeySet, error) {
	// Taken from NewHMACKey
	key := &[32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var sliceKey []byte
	copy(sliceKey, key[:])

	return &jose.JsonWebKeySet{
		Keys: []jose.JsonWebKey{
			{
				Algorithm:    "HS256",
				Key:          sliceKey,
				KeyID:        id,
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}
