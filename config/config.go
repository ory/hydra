// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/hmac"
	"github.com/ory/hydra/metrics"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	lmem "github.com/ory/ladon/manager/memory"
	lsql "github.com/ory/ladon/manager/sql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/yaml.v2"
)

type Config struct {
	// These are used by client commands
	ClusterURL   string `mapstructure:"CLUSTER_URL" yaml:"cluster_url"`
	ClientID     string `mapstructure:"CLIENT_ID" yaml:"client_id,omitempty"`
	ClientSecret string `mapstructure:"CLIENT_SECRET" yaml:"client_secret,omitempty"`

	// These are used by the host command
	BindPort                        int    `mapstructure:"PORT" yaml:"-"`
	BindHost                        string `mapstructure:"HOST" yaml:"-"`
	Issuer                          string `mapstructure:"ISSUER" yaml:"-"`
	SystemSecret                    string `mapstructure:"SYSTEM_SECRET" yaml:"-"`
	DatabaseURL                     string `mapstructure:"DATABASE_URL" yaml:"-"`
	DatabasePlugin                  string `mapstructure:"DATABASE_PLUGIN" yaml:"-"`
	ConsentURL                      string `mapstructure:"CONSENT_URL" yaml:"-"`
	AllowTLSTermination             string `mapstructure:"HTTPS_ALLOW_TERMINATION_FROM" yaml:"-"`
	BCryptWorkFactor                int    `mapstructure:"BCRYPT_COST" yaml:"-"`
	AccessTokenLifespan             string `mapstructure:"ACCESS_TOKEN_LIFESPAN" yaml:"-"`
	ScopeStrategy                   string `mapstructure:"SCOPE_STRATEGY" yaml:"-"`
	AuthCodeLifespan                string `mapstructure:"AUTH_CODE_LIFESPAN" yaml:"-"`
	IDTokenLifespan                 string `mapstructure:"ID_TOKEN_LIFESPAN" yaml:"-"`
	ChallengeTokenLifespan          string `mapstructure:"CHALLENGE_TOKEN_LIFESPAN" yaml:"-"`
	CookieSecret                    string `mapstructure:"COOKIE_SECRET" yaml:"-"`
	LogLevel                        string `mapstructure:"LOG_LEVEL" yaml:"-"`
	LogFormat                       string `mapstructure:"LOG_FORMAT" yaml:"-"`
	AccessControlResourcePrefix     string `mapstructure:"RESOURCE_NAME_PREFIX" yaml:"-"`
	OpenIDDiscoveryClaimsSupported  string `mapstructure:"OIDC_DISCOVERY_CLAIMS_SUPPORTED" yaml:"-"`
	OpenIDDiscoveryScopesSupported  string `mapstructure:"OIDC_DISCOVERY_SCOPES_SUPPORTED" yaml:"-"`
	OpenIDDiscoveryUserinfoEndpoint string `mapstructure:"OIDC_DISCOVERY_USERINFO_ENDPOINT" yaml:"-"`
	ForceHTTP                       bool   `yaml:"-"`

	BuildVersion string                  `yaml:"-"`
	BuildHash    string                  `yaml:"-"`
	BuildTime    string                  `yaml:"-"`
	logger       *logrus.Logger          `yaml:"-"`
	metrics      *metrics.MetricsManager `yaml:"-"`
	cluster      *url.URL                `yaml:"-"`
	oauth2Client *http.Client            `yaml:"-"`
	context      *Context                `yaml:"-"`
	systemSecret []byte                  `yaml:"-"`
}

func (c *Config) GetScopeStrategy() fosite.ScopeStrategy {
	if c.ScopeStrategy == "DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY" {
		c.GetLogger().Warn("Using deprecated hierarchical scope strategy, consider upgrading to wildcards.")
		return fosite.HierarchicScopeStrategy
	}

	return fosite.WildcardScopeStrategy
}

func matchesRange(r *http.Request, ranges []string) error {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, rn := range ranges {
		_, cidr, err := net.ParseCIDR(rn)
		if err != nil {
			return errors.WithStack(err)
		}
		addr := net.ParseIP(ip)
		if cidr.Contains(addr) {
			return nil
		}
	}
	return errors.Errorf("Remote address %s does not match cidr ranges %v", ip, ranges)
}

func newLogger(c *Config) *logrus.Logger {
	var (
		err    error
		logger = logrus.New()
	)

	if c.LogFormat == "json" {
		logger.Formatter = new(logrus.JSONFormatter)
	}

	logger.Level, err = logrus.ParseLevel(c.LogLevel)
	if err != nil {
		logger.Errorf("Couldn't parse log level: %s", c.LogLevel)
		logger.Level = logrus.InfoLevel
	}

	return logger
}

func (c *Config) GetLogger() *logrus.Logger {
	if c.logger == nil {
		c.logger = newLogger(c)
	}

	return c.logger
}

func (c *Config) GetMetrics() *metrics.MetricsManager {
	if c.metrics == nil {
		c.metrics = metrics.NewMetricsManager(c.Issuer, c.DatabaseURL, c.GetLogger())
	}

	return c.metrics
}

func (c *Config) DoesRequestSatisfyTermination(r *http.Request) error {
	if c.AllowTLSTermination == "" {
		return errors.New("TLS termination is not enabled")
	}

	if r.URL.Path == "/health" {
		return nil
	}

	ranges := strings.Split(c.AllowTLSTermination, ",")
	if err := matchesRange(r, ranges); err != nil {
		return err
	}

	proto := r.Header.Get("X-Forwarded-Proto")
	if proto == "" {
		return errors.New("X-Forwarded-Proto header is missing")
	} else if proto != "https" {
		return errors.Errorf("Expected X-Forwarded-Proto header to be https, got %s", proto)
	}

	return nil
}

func (c *Config) GetChallengeTokenLifespan() time.Duration {
	d, err := time.ParseDuration(c.ChallengeTokenLifespan)
	if err != nil {
		c.GetLogger().Warnf("Could not parse challenge token lifespan value (%s). Defaulting to 10m", c.AccessTokenLifespan)
		return time.Minute * 10
	}
	return d
}

func (c *Config) GetAccessTokenLifespan() time.Duration {
	d, err := time.ParseDuration(c.AccessTokenLifespan)
	if err != nil {
		c.GetLogger().Warnf("Could not parse access token lifespan value (%s). Defaulting to 1h", c.AccessTokenLifespan)
		return time.Hour
	}
	return d
}

func (c *Config) GetAuthCodeLifespan() time.Duration {
	d, err := time.ParseDuration(c.AuthCodeLifespan)
	if err != nil {
		c.GetLogger().Warnf("Could not parse auth code lifespan value (%s). Defaulting to 10m", c.AuthCodeLifespan)
		return time.Minute * 10
	}
	return d
}

func (c *Config) GetIDTokenLifespan() time.Duration {
	d, err := time.ParseDuration(c.IDTokenLifespan)
	if err != nil {
		c.GetLogger().Warnf("Could not parse id token lifespan value (%s). Defaulting to 1h", c.IDTokenLifespan)
		return time.Hour
	}
	return d
}

func (c *Config) Context() *Context {
	if c.context != nil {
		return c.context
	}

	var connection interface{} = &MemoryConnection{}
	if c.DatabaseURL == "" {
		c.GetLogger().Fatalf(`DATABASE_URL is not set, use "export DATABASE_URL=memory" for an in memory storage or the documented database adapters.`)
	} else if c.DatabasePlugin != "" {
		c.GetLogger().Infof("Database plugin set to %s", c.DatabasePlugin)
		pc := &PluginConnection{Config: c, Logger: c.GetLogger()}
		if err := pc.Connect(); err != nil {
			c.GetLogger().Fatalf("Could not connect via database plugin: %s", err)
		}
		connection = pc
	} else if c.DatabaseURL != "memory" {
		u, err := url.Parse(c.DatabaseURL)
		if err != nil {
			c.GetLogger().Fatalf("Could not parse DATABASE_URL: %s", err)
		}

		switch u.Scheme {
		case "postgres":
			fallthrough
		case "mysql":
			connection = &SQLConnection{
				URL: u,
				L:   c.GetLogger(),
			}
			break
		default:
			c.GetLogger().Fatalf(`Unknown DSN "%s" in DATABASE_URL: %s`, u.Scheme, c.DatabaseURL)
		}
	}

	var groupManager group.Manager
	var manager ladon.Manager
	switch con := connection.(type) {
	case *MemoryConnection:
		c.GetLogger().Printf("DATABASE_URL set to memory, connecting to ephermal in-memory database.")
		manager = lmem.NewMemoryManager()
		groupManager = group.NewMemoryManager()
		break
	case *SQLConnection:
		manager = lsql.NewSQLManager(con.GetDatabase(), nil)
		groupManager = &group.SQLManager{
			DB: con.GetDatabase(),
		}
		break
	case *PluginConnection:
		var err error
		manager, err = con.NewPolicyManager()
		if err != nil {
			c.GetLogger().Fatalf("Could not load policy manager plugin %s", err)
		}

		groupManager, err = con.NewGroupManager()
		if err != nil {
			c.GetLogger().Fatalf("Could not load group manager plugin %s", err)
		}
		break
	default:
		panic("Unknown connection type.")
	}

	c.context = &Context{
		Connection: connection,
		Hasher: &fosite.BCrypt{
			WorkFactor: c.BCryptWorkFactor,
		},
		LadonManager: manager,
		FositeStrategy: &foauth2.HMACSHAStrategy{
			Enigma: &hmac.HMACStrategy{
				GlobalSecret: c.GetSystemSecret(),
			},
			AccessTokenLifespan:   c.GetAccessTokenLifespan(),
			AuthorizeCodeLifespan: c.GetAuthCodeLifespan(),
		},
		GroupManager: groupManager,
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

type transporter struct {
	*http.Transport
	FakeTLSTermination bool
}

func (t *transporter) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.FakeTLSTermination {
		req.Header.Set("X-Forwarded-Proto", "https")
	}

	return t.Transport.RoundTrip(req)
}

func (c *Config) OAuth2Client(cmd *cobra.Command) *http.Client {
	if c.oauth2Client != nil {
		return c.oauth2Client
	}

	oauthConfig := clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     pkg.JoinURLStrings(c.ClusterURL, "/oauth2/token"),
		Scopes:       []string{"hydra", "hydra.*"},
	}

	fakeTlsTermination, _ := cmd.Flags().GetBool("fake-tls-termination")
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
		Transport: &transporter{
			FakeTLSTermination: fakeTlsTermination,
			Transport:          &http.Transport{},
		},
	})

	if ok, _ := cmd.Flags().GetBool("skip-tls-verify"); ok {
		// fmt.Println("Warning: Skipping TLS Certificate Verification.")
		ctx = context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
			Transport: &transporter{
				FakeTLSTermination: fakeTlsTermination,
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			},
		})
	}

	c.oauth2Client = oauthConfig.Client(ctx)
	if _, err := c.oauth2Client.Get(c.ClusterURL); err != nil {
		fmt.Printf("Could not authenticate, because: %s\n", err)
		fmt.Println("This can have multiple reasons, like a wrong cluster or wrong credentials. To resolve this, run `hydra connect`.")
		fmt.Println("You can disable TLS verification using the `--skip-tls-verify` flag.")
		os.Exit(1)
	}

	return c.oauth2Client
}

func (c *Config) GetCookieSecret() []byte {
	if c.CookieSecret != "" {
		return []byte(c.CookieSecret)
	}
	return c.GetSystemSecret()
}

func (c *Config) GetSystemSecret() []byte {
	if len(c.systemSecret) > 0 {
		return c.systemSecret
	}

	var secret = []byte(c.SystemSecret)
	if len(secret) >= 16 {
		hash := sha256.Sum256(secret)
		secret = hash[:]
		c.systemSecret = secret
		return secret
	}

	c.GetLogger().Warnf("Expected system secret to be at least %d characters long, got %d characters.", 32, len(c.SystemSecret))
	c.GetLogger().Infoln("Generating a random system secret...")
	var err error
	secret, err = pkg.GenerateSecret(32)
	pkg.Must(err, "Could not generate global secret: %s", err)
	c.GetLogger().Infof("Generated system secret: %s", secret)
	hash := sha256.Sum256(secret)
	secret = hash[:]
	c.systemSecret = secret
	c.GetLogger().Warnln("WARNING: DO NOT generate system secrets in production. The secret will be leaked to the logs.")
	return secret
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.BindHost, c.BindPort)
}

func (c *Config) Persist() error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return errors.WithStack(err)
	}

	c.GetLogger().Infof("Persisting config in file %s", viper.ConfigFileUsed())
	if err := ioutil.WriteFile(viper.ConfigFileUsed(), out, 0700); err != nil {
		return errors.Errorf(`Could not write to "%s" because: %s`, viper.ConfigFileUsed(), err)
	}

	return nil
}
