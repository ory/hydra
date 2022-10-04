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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package jwk

import (
	"context"
	"net"

	"github.com/ory/x/josex"

	"github.com/gofrs/uuid"
	"gopkg.in/square/go-jose.v2"

	"github.com/ory/fosite"
	"github.com/ory/hydra/driver/config"

	"github.com/pkg/errors"

	"github.com/ory/fosite/token/jwt"
)

type JWTSigner interface {
	GetPublicKeyID(ctx context.Context) (string, error)
	GetPublicKey(ctx context.Context) (jose.JSONWebKey, error)
	jwt.Signer
}

type DefaultJWTSigner struct {
	*jwt.DefaultSigner
	r     InternalRegistry
	c     *config.DefaultProvider
	setID string
}

func NewDefaultJWTSigner(c *config.DefaultProvider, r InternalRegistry, setID string) *DefaultJWTSigner {
	j := &DefaultJWTSigner{c: c, r: r, setID: setID, DefaultSigner: &jwt.DefaultSigner{}}
	j.DefaultSigner.GetPrivateKey = j.getPrivateKey
	return j
}

func (j *DefaultJWTSigner) getKeys(ctx context.Context) (private *jose.JSONWebKey, err error) {
	private, err = GetOrGenerateKeys(ctx, j.r, j.r.KeyManager(), j.setID, uuid.Must(uuid.NewV4()).String(), string(jose.RS256))
	if err == nil {
		return private, nil
	}

	var netError net.Error
	if errors.As(err, &netError) {
		return nil, errors.WithStack(fosite.ErrServerError.
			WithHintf(`Could not ensure that signing keys for "%s" exists. A network error occurred, see error for specific details.`, j.setID))
	}

	return nil, errors.WithStack(fosite.ErrServerError.
		WithWrap(err).
		WithHintf(`Could not ensure that signing keys for "%s" exists. If you are running against a persistent SQL database this is most likely because your "secrets.system" ("SECRETS_SYSTEM" environment variable) is not set or changed. When running with an SQL database backend you need to make sure that the secret is set and stays the same, unless when doing key rotation. This may also happen when you forget to run "hydra migrate sql..`, j.setID))
}

func (j *DefaultJWTSigner) GetPublicKeyID(ctx context.Context) (string, error) {
	private, err := j.getKeys(ctx)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return josex.ToPublicKey(private).KeyID, nil
}

func (j *DefaultJWTSigner) GetPublicKey(ctx context.Context) (jose.JSONWebKey, error) {
	private, err := j.getKeys(ctx)
	if err != nil {
		return jose.JSONWebKey{}, errors.WithStack(err)
	}
	return josex.ToPublicKey(private), nil
}

func (j *DefaultJWTSigner) getPrivateKey(ctx context.Context) (interface{}, error) {
	private, err := j.getKeys(ctx)
	if err != nil {
		return nil, err
	}

	return private, nil
}
