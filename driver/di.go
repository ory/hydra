// Copyright © 2024 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"github.com/ory/fosite/handler/oauth2"
)

type RegistryModifier func(r *RegistrySQL) error

func WithRegistryModifiers(f ...RegistryModifier) OptionsModifier {
	return func(o *options) {
		o.registryModifiers = f
	}
}

func RegistryWithHMACSHAStrategy(s func(r *RegistrySQL) oauth2.CoreStrategy) RegistryModifier {
	return func(r *RegistrySQL) error {
		r.hmacs = s(r)
		return nil
	}
}
