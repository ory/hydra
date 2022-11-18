// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package config

type Provider interface {
	Config() *DefaultProvider
}
