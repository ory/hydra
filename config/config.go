package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"net/http"
	"net/url"

	"github.com/ory-am/fosite/hash"
	"github.com/ory-am/hydra/pkg"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/yaml.v2"
	"github.com/ory-am/ladon"
	"github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/ory-am/fosite/handler/core/strategy"
	"github.com/ory-am/fosite/token/hmac"
)

type Config struct {
	BindPort     int `mapstructure:"port" yaml:"-"`

	BindHost     string `mapstructure:"host" yaml:"-"`

	Issuer       string `mapstructure:"issuer" yaml:"-"`

	SystemSecret []byte `mapstructure:"system_secret" yaml:"-"`

	ConsentURL   string `mapstructure:"consent_url" yaml:"-"`

	ClusterURL   string `mapstructure:"cluster_url" yaml:"cluster_url"`

	ClientID     string `mapstructure:"client_id" yaml:"client_id"`

	ClientSecret string `mapstructure:"client_secret" yaml:"client_secret"`

	cluster      *url.URL

	oauth2Client *http.Client

	context      *Context
}

func (c *Config) Context() *Context {
	if c.context != nil {
		return c.context
	}

	manager := ladon.NewMemoryManager()
	c.context = &Context{
		Connection: &MemoryConnection{},
		Hasher: &hash.BCrypt{
			WorkFactor: 11,
		},
		LadonManager: manager,
		FositeStrategy: &strategy.HMACSHAStrategy{
			Enigma: &hmac.HMACStrategy{
				GlobalSecret: c.GetSystemSecret(),
			},
		},
	}

	return c.context
}

func (c *Config) Resolve(join ...string) *url.URL {
	if c.cluster == nil {
		cluster, err := url.Parse(c.ClusterURL)
		c.cluster = cluster
		pkg.Must(err, "Could not parse cluster url: %s", err)
	}

	if len(join) == 0 {
		return c.cluster
	}

	return pkg.JoinURL(c.cluster, join...)
}

func (c *Config) OAuth2Client() *http.Client {
	if c.oauth2Client != nil {
		return c.oauth2Client
	}

	oauthConfig := clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     pkg.JoinURLStrings(c.ClusterURL, "/oauth2/token"),
		Scopes:       []string{
			"core",
			"hydra",
			"hydra.clients",
		},
	}

	_, err := oauthConfig.Token(context.Background())
	pkg.Must(err, "Could not authenticate: %s", err)
	c.oauth2Client = oauthConfig.Client(context.Background())
	return c.oauth2Client
}

func (c *Config) GetSystemSecret() []byte {
	if len(c.SystemSecret) >= 8 {
		return c.SystemSecret
	}

	var err error
	c.SystemSecret, err = pkg.GenerateSecret(32)
	pkg.Must(err, "Could not generate global secret: %s", err)
	logrus.Warnln("No system secret specified.")
	logrus.Warnf("Generated system secret: %s", c.SystemSecret)
	logrus.Warnln("Do not auto-generate system secrets in production.")
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

func (c *Config) Persist() error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return errors.New(err)
	}

	if err := ioutil.WriteFile(getConfigPath(), out, 0700); err != nil {
		return errors.New(err)
	}
	return nil
}

func getConfigPath() string {
	path := getUserHome() + "/.hydra.yml"
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	p, err := filepath.Abs(path)
	pkg.Must(err, "Could not fetch configuration path because %s", err)
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
