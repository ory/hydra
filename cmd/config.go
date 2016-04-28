package cmd

type config struct {
	BindPort int `mapstructure:"port"`

	BindHost string `mapstructure:"host"`

	SystemSecret string `mapstructure:"system_secret"`

	SelfURL string `mapstructure:"self_url"`

	ConsentURL string
}
