/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package config

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/hmac"
	"github.com/ory/go-convenience/urlx"
	"github.com/ory/hydra/health"
	"github.com/ory/hydra/metrics/prometheus"
	"github.com/ory/hydra/metrics/telemetry"
	"github.com/ory/hydra/pkg"
	"github.com/ory/sqlcon"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	// These are used by client commands
	EndpointURL string `mapstructure:"HYDRA_URL" yaml:"-"`

	// These are used by the host command
	BindPort                         int    `mapstructure:"PORT" yaml:"-"`
	BindHost                         string `mapstructure:"HOST" yaml:"-"`
	Issuer                           string `mapstructure:"OAUTH2_ISSUER_URL" yaml:"-"`
	SystemSecret                     string `mapstructure:"SYSTEM_SECRET" yaml:"-"`
	DatabaseURL                      string `mapstructure:"DATABASE_URL" yaml:"-"`
	DatabasePlugin                   string `mapstructure:"DATABASE_PLUGIN" yaml:"-"`
	ConsentURL                       string `mapstructure:"OAUTH2_CONSENT_URL" yaml:"-"`
	LoginURL                         string `mapstructure:"OAUTH2_LOGIN_URL" yaml:"-"`
	ErrorURL                         string `mapstructure:"OAUTH2_ERROR_URL" yaml:"-"`
	AllowTLSTermination              string `mapstructure:"HTTPS_ALLOW_TERMINATION_FROM" yaml:"-"`
	BCryptWorkFactor                 int    `mapstructure:"BCRYPT_COST" yaml:"-"`
	AccessTokenLifespan              string `mapstructure:"ACCESS_TOKEN_LIFESPAN" yaml:"-"`
	ScopeStrategy                    string `mapstructure:"SCOPE_STRATEGY" yaml:"-"`
	AuthCodeLifespan                 string `mapstructure:"AUTH_CODE_LIFESPAN" yaml:"-"`
	IDTokenLifespan                  string `mapstructure:"ID_TOKEN_LIFESPAN" yaml:"-"`
	ChallengeTokenLifespan           string `mapstructure:"CHALLENGE_TOKEN_LIFESPAN" yaml:"-"`
	CookieSecret                     string `mapstructure:"COOKIE_SECRET" yaml:"-"`
	LogLevel                         string `mapstructure:"LOG_LEVEL" yaml:"-"`
	LogFormat                        string `mapstructure:"LOG_FORMAT" yaml:"-"`
	AccessControlResourcePrefix      string `mapstructure:"RESOURCE_NAME_PREFIX" yaml:"-"`
	OpenIDDiscoveryClaimsSupported   string `mapstructure:"OIDC_DISCOVERY_CLAIMS_SUPPORTED" yaml:"-"`
	OpenIDDiscoveryScopesSupported   string `mapstructure:"OIDC_DISCOVERY_SCOPES_SUPPORTED" yaml:"-"`
	OpenIDDiscoveryUserinfoEndpoint  string `mapstructure:"OIDC_DISCOVERY_USERINFO_ENDPOINT" yaml:"-"`
	SendOAuth2DebugMessagesToClients bool   `mapstructure:"OAUTH2_SHARE_ERROR_DEBUG" yaml:"-"`
	ForceHTTP                        bool   `yaml:"-"`

	BuildVersion string                     `yaml:"-"`
	BuildHash    string                     `yaml:"-"`
	BuildTime    string                     `yaml:"-"`
	logger       *logrus.Logger             `yaml:"-"`
	telemetry    *telemetry.MetricsManager  `yaml:"-"`
	prometheus   *prometheus.MetricsManager `yaml:"-"`
	cluster      *url.URL                   `yaml:"-"`
	oauth2Client *http.Client               `yaml:"-"`
	context      *Context                   `yaml:"-"`
	systemSecret []byte                     `yaml:"-"`
}

func (c *Config) GetClusterURLWithoutTailingSlash(cmd *cobra.Command) string {
	if endpoint, _ := cmd.Flags().GetString("endpoint"); endpoint != "" {
		return strings.TrimRight(endpoint, "/")
	}
	return strings.TrimRight(c.EndpointURL, "/")
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

func (c *Config) GetTelemetryMetrics() *telemetry.MetricsManager {
	if c.telemetry == nil {
		c.telemetry = telemetry.NewMetricsManager(c.Issuer, c.DatabaseURL, c.GetLogger(), c.BuildVersion, c.BuildHash, c.BuildTime)
	}

	return c.telemetry
}

func (c *Config) GetPrometheusMetrics() *prometheus.MetricsManager {
	if c.prometheus == nil {
		c.GetLogger().Info("Setting up Prometheus metrics")
		c.prometheus = prometheus.NewMetricsManager(c.BuildVersion, c.BuildHash, c.BuildTime)
	}

	return c.prometheus
}

func (c *Config) DoesRequestSatisfyTermination(r *http.Request) error {
	if c.AllowTLSTermination == "" {
		return errors.New("TLS termination is not enabled")
	}

	if r.URL.Path == health.AliveCheckPath || r.URL.Path == health.ReadyCheckPath {
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
			connection, err = sqlcon.NewSQLConnection(c.DatabaseURL, c.GetLogger())
			if err != nil {
				c.GetLogger().WithError(err).Fatalf(`Unable to initialize SQL connection`)
			}
			break
		default:
			c.GetLogger().Fatalf(`Unknown DSN "%s" in DATABASE_URL: %s`, u.Scheme, c.DatabaseURL)
		}
	}

	c.context = &Context{
		Connection: connection,
		Hasher: &fosite.BCrypt{
			WorkFactor: c.BCryptWorkFactor,
		},
		FositeStrategy: &foauth2.HMACSHAStrategy{
			Enigma: &hmac.HMACStrategy{
				GlobalSecret: c.GetSystemSecret(),
			},
			AccessTokenLifespan:   c.GetAccessTokenLifespan(),
			AuthorizeCodeLifespan: c.GetAuthCodeLifespan(),
		},
	}

	return c.context
}

func (c *Config) Resolve(join ...string) *url.URL {
	if c.cluster == nil {
		cluster, err := url.Parse(c.EndpointURL)
		c.cluster = cluster
		pkg.Must(err, "Could not parse cluster url: %s", err)
	}

	if len(join) == 0 {
		return c.cluster
	}

	return urlx.AppendPaths(c.cluster, join...)
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
