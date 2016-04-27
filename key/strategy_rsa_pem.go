package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/go-errors/errors"
)

type RSAPEMStrategy struct {}

func (s *RSAPEMStrategy) AsymmetricKey(id string) (*AsymmetricKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	}

	if err = priv.Validate(); err != nil {
		return  nil,errors.Errorf("Validation failed because %s", err)
	}

	privDer := x509.MarshalPKCS1PrivateKey(priv)
	privBlk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDer,
	}

	pub := priv.PublicKey
	pubDer, err := x509.MarshalPKIXPublicKey(&pub)
	if err != nil {
		return  nil, errors.Errorf("Failed to get der format for PublicKey because %s", err)
	}

	pubBlk := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubDer,
	}

	return &AsymmetricKey{
		ID: id,
		Private: pem.EncodeToMemory(&privBlk),
		Public: pem.EncodeToMemory(&pubBlk),
	}, nil
}