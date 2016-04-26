package config

import "net/url"

type Config struct {
	Viper
	DBAL
}

type Viper struct {
	BindPort int `mapstructure:"port"`

	BindHost string `mapstructure:"host"`

	SystemSecret string `mapstructure:"system_secret"`

	SelfURL string `mapstructure:"self_url"`
}

type DBAL struct {
	IdentityProviders map[string]*EndpointConfig

	ConnectionProviders map[string]*EndpointConfig

	ConsentEndpoint *EndpointConfig
}

type EndpointConfig struct {
	ID  string
	URL *url.URL
}
