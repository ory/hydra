// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/x"
)

var ErrUnsupportedKeyAlgorithm = &fosite.RFC6749Error{
	CodeField:        http.StatusBadRequest,
	ErrorField:       http.StatusText(http.StatusBadRequest),
	DescriptionField: "Unsupported key algorithm",
}

var ErrUnsupportedEllipticCurve = &fosite.RFC6749Error{
	CodeField:        http.StatusBadRequest,
	ErrorField:       http.StatusText(http.StatusBadRequest),
	DescriptionField: "Unsupported elliptic curve",
}

var ErrMinimalRsaKeyLength = &fosite.RFC6749Error{
	CodeField:        http.StatusBadRequest,
	ErrorField:       http.StatusText(http.StatusBadRequest),
	DescriptionField: "Unsupported RSA key length",
}

type (
	Manager interface {
		GenerateAndPersistKeySet(ctx context.Context, set, kid, alg, use string) (*jose.JSONWebKeySet, error)

		AddKey(ctx context.Context, set string, key *jose.JSONWebKey) error

		AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) error

		UpdateKey(ctx context.Context, set string, key *jose.JSONWebKey) error

		UpdateKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) error

		GetKey(ctx context.Context, set, kid string) (*jose.JSONWebKeySet, error)

		GetKeySet(ctx context.Context, set string) (*jose.JSONWebKeySet, error)

		DeleteKey(ctx context.Context, set, kid string) error

		DeleteKeySet(ctx context.Context, set string) error
	}
	ManagerProvider interface {
		KeyManager() Manager
	}

	SQLData struct {
		ID  uuid.UUID `db:"pk"`
		NID uuid.UUID `json:"-" db:"nid"`
		// This field is deprecated and will be removed
		PKDeprecated int64     `json:"-" db:"pk_deprecated"`
		Set          string    `db:"sid"`
		KID          string    `db:"kid"`
		Version      int       `db:"version"`
		CreatedAt    time.Time `db:"created_at"`
		Key          string    `db:"keydata"`
	}

	SQLDataRows []SQLData
)

func (d SQLData) TableName() string { return "hydra_jwk" }

func (d SQLDataRows) ToJWK(ctx context.Context, aes *aead.AESGCM) (keys *jose.JSONWebKeySet, err error) {
	if len(d) == 0 {
		return nil, errors.Wrap(x.ErrNotFound, "")
	}

	keys = &jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, len(d)),
	}
	for i, d := range d {
		key, err := aes.Decrypt(ctx, d.Key, nil)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if err := json.Unmarshal(key, &keys.Keys[i]); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return keys, nil
}
