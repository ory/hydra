package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/yaml.v2"
)

var config = new(configuration)

type configuration struct {
	BindPort     int `mapstructure:"port" yaml:"-"`

	BindHost     string `mapstructure:"host" yaml:"-"`

	Issuer       string `mapstructure:"issuer" yaml:"-"`

	SystemSecret []byte `mapstructure:"system_secret" yaml:"-"`

	ConsentURL   string `mapstructure:"consent_url" yaml:"-"`

	ClusterURL   string `mapstructure:"cluster_url" yaml:"cluster_url"`

	ClientID     string `mapstructure:"client_id" yaml:"client_id"`

	ClientSecret string `mapstructure:"client_secret" yaml:"client_secret"`
}

func (c *configuration) GetSystemSecret() []byte {
	if len(c.SystemSecret) >= 8 {
		return c.SystemSecret
	}

	fmt.Println("No global secret was set. Generating a random one...")
	c.SystemSecret = generateSecret(32)
	fmt.Printf("A global secret was generated:\n%s\n", c.SystemSecret)
	return c.SystemSecret
}

func (c *configuration) GetAddress() string {
	if c.BindPort == 0 {
		c.BindPort = 4444
	}
	return fmt.Sprintf("%s:%d", c.BindHost, c.BindPort)
}

func (c *configuration) GetIssuer() string {
	if c.Issuer == "" {
		c.Issuer = "hydra"
	}
	return c.Issuer
}

func (c *configuration) GetAccessTokenLifespan() time.Duration {
	return time.Hour
}

func (c *configuration) Save() error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(getConfigPath(), out, 0700); err != nil {
		return err
	}
	return nil
}

func getConfigPath() string {
	path := getUserHome() + "/.hydra.yml"
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	p, err := filepath.Abs(path)
	if err != nil {
		fatal("Could not fetch configuration path because %s", err)
	}
	return filepath.Clean(p)
}

func getUserHome() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
