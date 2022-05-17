package config

type Provider interface {
	Config() *DefaultProvider
}
