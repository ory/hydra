// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
)

type veryInsecureRS256Generator struct{}

func (g *veryInsecureRS256Generator) Generate(id, use string) (*jose.JSONWebKeySet, error) {
	/* #nosec G403 - this is ok because this generator is only used in tests. */
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return nil, errors.Errorf("Could not generate key because %s", err)
	} else if err = key.Validate(); err != nil {
		return nil, errors.Errorf("Validation failed because %s", err)
	}

	if id == "" {
		id = uuid.New()
	}

	// jose does not support this...
	key.Precomputed = rsa.PrecomputedValues{}
	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Algorithm:    "RS256",
				Key:          key,
				Use:          use,
				KeyID:        id,
				Certificates: []*x509.Certificate{},
			},
		},
	}, nil
}
