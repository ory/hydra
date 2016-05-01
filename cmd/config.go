package cmd

type config struct {
	BindPort int `mapstructure:"port"`

	BindHost string `mapstructure:"host"`

	SystemSecret string `mapstructure:"system_secret"`

	ConsentURL string `mapstructure:"consent_url"`
}
