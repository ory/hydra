// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/fosite/token/jwt"
	"github.com/ory/hydra/v2/x"
)

var _ fosite.GlobalSecretProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetGlobalSecret(ctx context.Context) ([]byte, error) {
	secrets := p.getProvider(ctx).Strings(KeyGetSystemSecret)

	if len(secrets) == 0 {
		p.l.Error("The system secret is not configured. Please provide one in the configuration file or environment variables.")
		return nil, errors.New("global secret is not configured")
	}

	secret := secrets[0]
	if len(secret) < 16 {
		p.l.Errorf("System secret must be undefined or have at least 16 characters but only has %d characters.", len(secret))
		return nil, errors.New("global secret is too short")
	}

	return x.HashStringSecret(secret), nil
}

var _ fosite.RotatedGlobalSecretsProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetRotatedGlobalSecrets(ctx context.Context) ([][]byte, error) {
	secrets := p.getProvider(ctx).Strings(KeyGetSystemSecret)

	if len(secrets) < 2 {
		return nil, nil
	}

	rotated := make([][]byte, len(secrets)-1)
	for i, secret := range secrets[1:] {
		rotated[i] = x.HashStringSecret(secret)
	}

	return rotated, nil
}

var _ fosite.BCryptCostProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetBCryptCost(ctx context.Context) int {
	return p.getProvider(ctx).IntF(KeyBCryptCost, 10)
}

var _ fosite.AccessTokenLifespanProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetAccessTokenLifespan(ctx context.Context) time.Duration {
	return p.getProvider(ctx).DurationF(KeyAccessTokenLifespan, time.Hour)
}

var _ fosite.RefreshTokenLifespanProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetRefreshTokenLifespan(ctx context.Context) time.Duration {
	return p.getProvider(ctx).DurationF(KeyRefreshTokenLifespan, time.Hour*720)
}

var _ fosite.VerifiableCredentialsNonceLifespanProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetVerifiableCredentialsNonceLifespan(ctx context.Context) time.Duration {
	return p.getProvider(ctx).DurationF(KeyVerifiableCredentialsNonceLifespan, time.Hour)
}

var _ fosite.IDTokenLifespanProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetIDTokenLifespan(ctx context.Context) time.Duration {
	return p.getProvider(ctx).DurationF(KeyIDTokenLifespan, time.Hour)
}

var _ fosite.AuthorizeCodeLifespanProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetAuthorizeCodeLifespan(ctx context.Context) time.Duration {
	return p.getProvider(ctx).DurationF(KeyAuthCodeLifespan, time.Minute*10)
}

var _ fosite.ScopeStrategyProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetScopeStrategy(ctx context.Context) fosite.ScopeStrategy {
	if strings.ToLower(p.getProvider(ctx).String(KeyScopeStrategy)) == "wildcard" {
		return fosite.WildcardScopeStrategy
	}
	return fosite.ExactScopeStrategy
}

var _ fosite.JWTScopeFieldProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetJWTScopeField(ctx context.Context) jwt.JWTScopeFieldEnum {
	switch strings.ToLower(p.getProvider(ctx).String(KeyJWTScopeClaimStrategy)) {
	case "string":
		return jwt.JWTScopeFieldString
	case "both":
		return jwt.JWTScopeFieldBoth
	case "list":
		return jwt.JWTScopeFieldList
	default:
		return jwt.JWTScopeFieldUnset
	}
}

func (p *DefaultProvider) GetUseLegacyErrorFormat(context.Context) bool {
	return false
}
