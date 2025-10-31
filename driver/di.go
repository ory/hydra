// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
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

func RegistryWithKeyManager(km func(r *RegistrySQL) (jwk.Manager, error)) RegistryModifier {
	return func(r *RegistrySQL) (err error) {
		r.keyManager, err = km(r)
		return err
	}
}
