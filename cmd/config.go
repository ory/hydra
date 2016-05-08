package cmd

import "fmt"

type config struct {
	BindPort int `mapstructure:"port"`

	BindHost string `mapstructure:"host"`

	Issuer string `mapstructure:"issuer"`

	SystemSecret []byte `mapstructure:"system_secret"`

	ConsentURL string `mapstructure:"consent_url"`

	BackendURL string `mapstructure:"backend_url"`
}

func (c *config) Addr() string {
	if c.BindPort == 0 {
		c.BindPort = 4444
	}
	return fmt.Sprintf("%s:%s", c.BindHost, c.BindPort)
}

func (c *config) Iss() string {
	if c.Issuer == "" {
		c.Issuer = "hydra"
	}
	return c.Issuer
}
