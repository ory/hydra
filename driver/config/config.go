// Copyright © 2022 Ory Corp

package config

type Provider interface {
	Config() *DefaultProvider
}
