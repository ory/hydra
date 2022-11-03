// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"strings"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/hydra/x"
	"github.com/ory/x/cmdx"
)

var _ fosite.GlobalSecretProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetGlobalSecret(ctx context.Context) []byte {
	secrets := p.getProvider(ctx).Strings(KeyGetSystemSecret)

	if len(secrets) == 0 {
		if p.generatedSecret != nil {
			return p.generatedSecret
		}

		p.l.Warnf("Configuration secrets.system is not set, generating a temporary, random secret...")
		secret, err := x.GenerateSecret(32)
		cmdx.Must(err, "Could not generate secret: %s", err)

		p.l.Warnf("Generated secret: %s", secret)
		p.generatedSecret = x.HashByteSecret(secret)

		p.l.Warnln("Do not use generate secrets in production. The secret will be leaked to the logs.")
		return x.HashByteSecret(secret)
	}

	secret := secrets[0]
	if len(secret) >= 16 {
		return x.HashStringSecret(secret)
	}

	p.l.Fatalf("System secret must be undefined or have at least 16 characters but only has %d characters.", len(secret))
	return nil
}

var _ fosite.RotatedGlobalSecretsProvider = (*DefaultProvider)(nil)

func (p *DefaultProvider) GetRotatedGlobalSecrets(ctx context.Context) [][]byte {
	secrets := p.getProvider(ctx).Strings(KeyGetSystemSecret)

	if len(secrets) < 2 {
		return nil
	}

	var rotated [][]byte
	for _, secret := range secrets[1:] {
		rotated = append(rotated, x.HashStringSecret(secret))
	}

	return rotated
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

func (p *DefaultProvider) GetUseLegacyErrorFormat(context.Context) bool {
	return false
}
