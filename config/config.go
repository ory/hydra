package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/yaml.v2"
	"net/http"
	"golang.org/x/oauth2/clientcredentials"
	"github.com/ory-am/hydra/pkg"
	"golang.org/x/net/context"
	"github.com/ory-am/fosite/hash"
)

type Config struct {
	BindPort int `mapstructure:"port" yaml:"-"`

	BindHost string `mapstructure:"host" yaml:"-"`

	Issuer string `mapstructure:"issuer" yaml:"-"`

	SystemSecret []byte `mapstructure:"system_secret" yaml:"-"`

	ConsentURL string `mapstructure:"consent_url" yaml:"-"`

	ClusterURL string `mapstructure:"cluster_url" yaml:"cluster_url"`

	ClientID string `mapstructure:"client_id" yaml:"client_id"`

	ClientSecret string `mapstructure:"client_secret" yaml:"client_secret"`
}

func (c *Config) Context() *Context {
	return &Context{
		Connection: &MemoryConnection{},
		Hasher: &hash.BCrypt{
			WorkFactor: 11,
		},
	}
}

func (c *Config) OAuth2Client() *http.Client {
	oauthConfig := clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     pkg.JoinURL(c.ClusterURL, "oauth2/token"),
		Scopes:       []string{"core", "hydra"},
	}

	if _, err := oauthConfig.Token(context.Background()); err != nil {
		return err
	}

	return oauthConfig.Client(context.Background())
}

func (c *Config) GetSystemSecret() []byte {
	if len(c.SystemSecret) >= 8 {
		return c.SystemSecret
	}

	fmt.Println("No global secret was set. Generating a random one...")
	c.SystemSecret = generateSecret(32)
	fmt.Printf("A global secret was generated:\n%s\n", c.SystemSecret)
	return c.SystemSecret
}

func (c *Config) GetAddress() string {
	if c.BindPort == 0 {
		c.BindPort = 4444
	}
	return fmt.Sprintf("%s:%d", c.BindHost, c.BindPort)
}

func (c *Config) GetIssuer() string {
	if c.Issuer == "" {
		c.Issuer = "hydra"
	}
	return c.Issuer
}

func (c *Config) GetAccessTokenLifespan() time.Duration {
	return time.Hour
}

func (c *Config) Save() error {
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
