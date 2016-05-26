package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"net/http"
	"net/url"

	"crypto/tls"
	"os"
	"strconv"

	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/ory-am/fosite/handler/core/strategy"
	"github.com/ory-am/fosite/hash"
	"github.com/ory-am/fosite/token/hmac"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/yaml.v2"
)

type Config struct {
	BindPort     int `mapstructure:"port" yaml:"port,omitempty"`

	BindHost     string `mapstructure:"host" yaml:"host,omitempty"`

	Issuer       string `mapstructure:"issuer" yaml:"issuer,omitempty"`

	SystemSecret []byte `mapstructure:"system_secret" yaml:"-"`

	DatabaseURL  string `mapstructure:"database_url" yaml:"database_url,omitempty"`

	ConsentURL   string `mapstructure:"consent_url" yaml:"consent_url,omitempty"`

	ClusterURL   string `mapstructure:"cluster_url" yaml:"cluster_url,omitempty"`

	ClientID     string `mapstructure:"client_id" yaml:"client_id,omitempty"`

	ClientSecret string `mapstructure:"client_secret" yaml:"client_secret,omitempty"`

	ForceHTTP    bool `mapstructure:"foolishly_force_http" yaml:"-"`

	cluster      *url.URL

	oauth2Client *http.Client

	context      *Context

	sync.Mutex
}

func (c *Config) GetClusterURL() string {
	c.Lock()
	defer c.Unlock()

	if c.ClusterURL == "" {
		bindHost := c.BindHost
		if bindHost == "" {
			bindHost = "localhost"
		}

		schema := "https"
		if c.ForceHTTP {
			schema = "http"
		}

		port := strconv.Itoa(c.BindPort)
		if c.BindPort == 0 {
			port = "4444"
		}

		c.ClusterURL = schema + "://" + bindHost + ":" + port
	}

	return c.ClusterURL
}

func (c *Config) Context() *Context {
	secret := c.GetSystemSecret()
	c.Lock()
	defer c.Unlock()

	if c.context != nil {
		return c.context
	}

	var connection interface{} = &MemoryConnection{}
	if c.DatabaseURL != "" {
		u, err := url.Parse(c.DatabaseURL)
		if err != nil {
			panic(errors.New(err))
		}

		switch u.Scheme {
		case "rethinkdb":
			connection = &RethinkDBConnection{URL: u}
			break
		}
	}

	var manager ladon.Manager
	switch con := connection.(type) {
	case *MemoryConnection:
		manager = ladon.NewMemoryManager()
	case *RethinkDBConnection:
		con.CreateTableIfNotExists("hydra_policies")
		m := &ladon.RethinkManager{Session: con.GetSession()}
		m.ColdStart()
		m.Watch(context.Background())
		manager = m
	default:
		panic("Unknown connection type.")
	}

	c.context = &Context{
		Connection: connection,
		Hasher: &hash.BCrypt{
			WorkFactor: 11,
		},
		LadonManager: manager,
		FositeStrategy: &strategy.HMACSHAStrategy{
			Enigma: &hmac.HMACStrategy{
				GlobalSecret: secret,
			},
		},
	}

	return c.context
}

func (c *Config) Resolve(join ...string) *url.URL {
	c.Lock()
	defer c.Unlock()

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

func (c *Config) OAuth2Client(cmd *cobra.Command) *http.Client {
	c.Lock()
	defer c.Unlock()

	if c.oauth2Client != nil {
		return c.oauth2Client
	}

	oauthConfig := clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     pkg.JoinURLStrings(c.ClusterURL, "/oauth2/token"),
		Scopes: []string{
			"core",
			"hydra",
		},
	}

	ctx := context.Background()
	if ok, _ := cmd.Flags().GetBool("skip-tls-verify"); ok {
		fmt.Println("Warning: Skipping TLS Certificate Verification.")
		ctx = context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}})
	}

	_, err := oauthConfig.Token(ctx)
	if err != nil {
		fmt.Printf("Could not authenticate, because: %s\n", err)
		fmt.Println("Did you forget to log on? Run `hydra connect`.")
		fmt.Println("Did you run Hydra without a valid TLS certificate? Make sure to use the `--skip-tls-verify` flag.")
		fmt.Println("Did you know you can skip `hydra connect` when running `hydra host --dangerous-auto-logon`? DO NOT use this flag in production!")
		os.Exit(1)
	}

	c.oauth2Client = oauthConfig.Client(ctx)
	return c.oauth2Client
}

func (c *Config) GetSystemSecret() []byte {
	c.Lock()
	defer c.Unlock()

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
	c.Lock()
	defer c.Unlock()

	if c.BindPort == 0 {
		c.BindPort = 4444
	}
	return fmt.Sprintf("%s:%d", c.BindHost, c.BindPort)
}

func (c *Config) GetIssuer() string {
	c.Lock()
	defer c.Unlock()

	if c.Issuer == "" {
		c.Issuer = "hydra"
	}
	return c.Issuer
}

func (c *Config) GetAccessTokenLifespan() time.Duration {
	c.Lock()
	defer c.Unlock()

	return time.Hour
}

func (c *Config) Persist() error {
	_ = c.GetIssuer()
	_ = c.GetAddress()
	_ = c.GetClusterURL()

	out, err := yaml.Marshal(c)
	if err != nil {
		return errors.New(err)
	}

	if err := ioutil.WriteFile(viper.ConfigFileUsed(), out, 0700); err != nil {
		return errors.Errorf(`Could not write to "%s" because: %s`, viper.ConfigFileUsed(), err)
	}
	return nil
}
