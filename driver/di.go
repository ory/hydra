// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"github.com/ory/hydra/v2/consent"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
	"github.com/ory/hydra/v2/fosite/handler/openid"
	"github.com/ory/hydra/v2/fosite/handler/pkce"
	"github.com/ory/hydra/v2/jwk"
)

type RegistryModifier func(r *RegistrySQL) error

func WithRegistryModifiers(f ...RegistryModifier) OptionsModifier {
	return func(o *options) {
		o.registryModifiers = append(o.registryModifiers, f...)
	}
}

func RegistryWithHMACSHAStrategy(s func(r *RegistrySQL) oauth2.CoreStrategy) RegistryModifier {
	return func(r *RegistrySQL) error {
		r.hmacs = s(r)
		return nil
	}
}

func RegistryWithJWTStrategy(s func(r *RegistrySQL) oauth2.AccessTokenStrategy) RegistryModifier {
	return func(r *RegistrySQL) error {
		r.jwtStrategy = s(r)
		return nil
	}
}

func RegistryWithAuthorizeCodeStrategy(s func(r *RegistrySQL) oauth2.AuthorizeCodeStrategy) RegistryModifier {
	return func(r *RegistrySQL) error {
		r.authorizeCodeStrategy = s(r)
		return nil
	}
}

func RegistryWithKeyManager(km func(r *RegistrySQL) (jwk.Manager, error)) RegistryModifier {
	return func(r *RegistrySQL) (err error) {
		r.keyManager, err = km(r)
		return err
	}
}

func RegistryWithOAuth2Provider(pr func(r *RegistrySQL) fosite.OAuth2Provider) RegistryModifier {
	return func(r *RegistrySQL) error {
		r.fop = pr(r)
		return nil
	}
}

func RegistryWithAccessTokenStorage(s func(r *RegistrySQL) oauth2.AccessTokenStorage) RegistryModifier {
	return func(r *RegistrySQL) error {
		r.accessTokenStorage = s(r)
		return nil
	}
}

func RegistryWithAuthorizeCodeStorage(s func(r *RegistrySQL) oauth2.AuthorizeCodeStorage) RegistryModifier {
	return func(r *RegistrySQL) error {
		r.authorizeCodeStorage = s(r)
		return nil
	}
}

func RegistryWithPKCERequestStorage(s func(r *RegistrySQL) pkce.PKCERequestStorage) RegistryModifier {
	return func(r *RegistrySQL) error {
		r.pkceRequestStorage = s(r)
		return nil
	}
}

func RegistryWithOpenIDConnectRequestStorage(s func(r *RegistrySQL) openid.OpenIDConnectRequestStorage) RegistryModifier {
	return func(r *RegistrySQL) error {
		r.openIDConnectRequestStorage = s(r)
		return nil
	}
}

func RegistryWithConsentManager(cm func(r *RegistrySQL) (consent.Manager, error)) RegistryModifier {
	return func(r *RegistrySQL) (err error) {
		r.consentManager, err = cm(r)
		return err
	}
}
