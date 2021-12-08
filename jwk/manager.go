/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package jwk

import (
	"context"
	"net/http"
	"time"

	"github.com/ory/fosite"

	jose "gopkg.in/square/go-jose.v2"
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

	SQLData struct {
		ID        int       `db:"pk"`
		Set       string    `db:"sid"`
		KID       string    `db:"kid"`
		Version   int       `db:"version"`
		CreatedAt time.Time `db:"created_at"`
		Key       string    `db:"keydata"`
	}
)

func (d SQLData) TableName() string {
	return "hydra_jwk"
}
