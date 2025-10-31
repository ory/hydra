// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"net"

	"github.com/go-jose/go-jose/v3"
	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/token/jwt"
	"github.com/ory/x/josex"
)

type JWTSigner interface {
	GetPublicKeyID(ctx context.Context) (string, error)
	GetPublicKey(ctx context.Context) (jose.JSONWebKey, error)
	jwt.Signer
}

type DefaultJWTSigner struct {
	*jwt.DefaultSigner
	r     InternalRegistry
	setID string
}

func NewDefaultJWTSigner(r InternalRegistry, setID string) *DefaultJWTSigner {
	j := &DefaultJWTSigner{r: r, setID: setID, DefaultSigner: &jwt.DefaultSigner{}}
	j.DefaultSigner.GetPrivateKey = j.getPrivateKey
	return j
}

func (j *DefaultJWTSigner) getKeys(ctx context.Context) (private *jose.JSONWebKey, err error) {
	private, err = GetOrGenerateKeys(ctx, j.r, j.r.KeyManager(), j.setID, string(jose.RS256))
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
		WithHintf(`Could not ensure that signing keys for "%s" exists. If you are running against a persistent SQL database this is most likely because your "secrets.system" ("SECRETS_SYSTEM" environment variable) is not set or changed. When running with an SQL database backend you need to make sure that the secret is set and stays the same, unless when doing key rotation. This may also happen when you forget to run "hydra migrate sql up -e".`, j.setID))
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
