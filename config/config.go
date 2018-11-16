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
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v1"

	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/hmac"
	"github.com/ory/go-convenience/stringslice"
	"github.com/ory/go-convenience/urlx"
	"github.com/ory/hydra/metrics/prometheus"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/tracing"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/healthx"
)

type Config struct {
	// These are used by client commands
	EndpointURL string `mapstructure:"HYDRA_URL" yaml:"-"`

	// These are used by the host command
	FrontendBindPort                 int     `mapstructure:"PUBLIC_PORT" yaml:"-"`
	FrontendBindHost                 string  `mapstructure:"PUBLIC_HOST" yaml:"-"`
	BackendBindPort                  int     `mapstructure:"ADMIN_PORT" yaml:"-"`
	BackendBindHost                  string  `mapstructure:"ADMIN_HOST" yaml:"-"`
	Issuer                           string  `mapstructure:"OAUTH2_ISSUER_URL" yaml:"-"`
	ClientRegistrationURL            string  `mapstructure:"OAUTH2_CLIENT_REGISTRATION_URL" yaml:"-"`
	SystemSecret                     string  `mapstructure:"SYSTEM_SECRET" yaml:"-"`
	RotatedSystemSecret              string  `mapstructure:"ROTATED_SYSTEM_SECRET" yaml:"-"`
	DatabaseURL                      string  `mapstructure:"DATABASE_URL" yaml:"-"`
	DatabasePlugin                   string  `mapstructure:"DATABASE_PLUGIN" yaml:"-"`
	ConsentURL                       string  `mapstructure:"OAUTH2_CONSENT_URL" yaml:"-"`
	LoginURL                         string  `mapstructure:"OAUTH2_LOGIN_URL" yaml:"-"`
	LogoutRedirectURL                string  `mapstructure:"OAUTH2_LOGOUT_REDIRECT_URL" yaml:"-"`
	DefaultClientScope               string  `mapstructure:"OIDC_DYNAMIC_CLIENT_REGISTRATION_DEFAULT_SCOPE" yaml:"-"`
	ErrorURL                         string  `mapstructure:"OAUTH2_ERROR_URL" yaml:"-"`
	AllowTLSTermination              string  `mapstructure:"HTTPS_ALLOW_TERMINATION_FROM" yaml:"-"`
	BCryptWorkFactor                 int     `mapstructure:"BCRYPT_COST" yaml:"-"`
	AccessTokenLifespan              string  `mapstructure:"ACCESS_TOKEN_LIFESPAN" yaml:"-"`
	ScopeStrategy                    string  `mapstructure:"SCOPE_STRATEGY" yaml:"-"`
	AuthCodeLifespan                 string  `mapstructure:"AUTH_CODE_LIFESPAN" yaml:"-"`
	RefreshTokenLifespan             string  `mapstructure:"REFRESH_TOKEN_LIFESPAN" yaml:"-"`
	IDTokenLifespan                  string  `mapstructure:"ID_TOKEN_LIFESPAN" yaml:"-"`
	LoginConsentRequestLifespan      string  `mapstructure:"LOGIN_CONSENT_REQUEST_LIFESPAN" yaml:"-"`
	CookieSecret                     string  `mapstructure:"COOKIE_SECRET" yaml:"-"`
	LogLevel                         string  `mapstructure:"LOG_LEVEL" yaml:"-"`
	LogFormat                        string  `mapstructure:"LOG_FORMAT" yaml:"-"`
	AccessControlResourcePrefix      string  `mapstructure:"RESOURCE_NAME_PREFIX" yaml:"-"`
	SubjectTypesSupported            string  `mapstructure:"OIDC_SUBJECT_TYPES_SUPPORTED" yaml:"-"`
	SubjectIdentifierAlgorithmSalt   string  `mapstructure:"OIDC_SUBJECT_TYPE_PAIRWISE_SALT" yaml:"-"`
	OpenIDDiscoveryClaimsSupported   string  `mapstructure:"OIDC_DISCOVERY_CLAIMS_SUPPORTED" yaml:"-"`
	OpenIDDiscoveryScopesSupported   string  `mapstructure:"OIDC_DISCOVERY_SCOPES_SUPPORTED" yaml:"-"`
	OpenIDDiscoveryUserinfoEndpoint  string  `mapstructure:"OIDC_DISCOVERY_USERINFO_ENDPOINT" yaml:"-"`
	SendOAuth2DebugMessagesToClients bool    `mapstructure:"OAUTH2_SHARE_ERROR_DEBUG" yaml:"-"`
	OAuth2AccessTokenStrategy        string  `mapstructure:"OAUTH2_ACCESS_TOKEN_STRATEGY" yaml:"-"`
	TracingProvider                  string  `mapstructure:"TRACING_PROVIDER" yaml:"-"`
	TracingServiceName               string  `mapstructure:"TRACING_SERVICE_NAME" yaml:"-"`
	JaegerSamplingServerUrl          string  `mapstructure:"TRACING_PROVIDER_JAEGER_SAMPLING_SERVER_URL" yaml:"-"`
	JaegerLocalAgentHostPort         string  `mapstructure:"TRACING_PROVIDER_JAEGER_LOCAL_AGENT_ADDRESS" yaml:"-"`
	JaegerSamplingType               string  `mapstructure:"TRACING_PROVIDER_JAEGER_SAMPLING_TYPE" yaml:"-"`
	JaegerSamplingValue              float64 `mapstructure:"TRACING_PROVIDER_JAEGER_SAMPLING_VALUE" yaml:"-"`
	ForceHTTP                        bool    `yaml:"-"`

	BuildVersion string                     `yaml:"-"`
	BuildHash    string                     `yaml:"-"`
	BuildTime    string                     `yaml:"-"`
	tracer       *tracing.Tracer            `yaml:"-"`
	logger       *logrus.Logger             `yaml:"-"`
	prometheus   *prometheus.MetricsManager `yaml:"-"`
	cluster      *url.URL                   `yaml:"-"`
	oauth2Client *http.Client               `yaml:"-"`
	context      *Context                   `yaml:"-"`
	systemSecret []byte                     `yaml:"-"`
}

func (c *Config) MustValidate() {
	if stringslice.Has(c.GetSubjectTypesSupported(), "pairwise") && c.OAuth2AccessTokenStrategy == "jwt" {
		c.GetLogger().Fatalf(`The pairwise subject identifier algorithm is not supported by the JWT OAuth 2.0 Access Token Strategy. Please remove "pairwise" from OIDC_SUBJECT_TYPES_SUPPORTED or set OAUTH2_ACCESS_TOKEN_STRATEGY to "opaque"`)
	}

	if stringslice.Has(c.GetSubjectTypesSupported(), "pairwise") && len(c.SubjectIdentifierAlgorithmSalt) < 8 {
		c.GetLogger().Fatalf(`The pairwise subject identifier algorithm was set but length of OIDC_SUBJECT_TYPE_PAIRWISE_SALT is too small (%d < 8), please set OIDC_SUBJECT_TYPE_PAIRWISE_SALT to a random string with 8 characters or more`, len(c.SubjectIdentifierAlgorithmSalt))
	}
}

func (c *Config) GetSubjectTypesSupported() []string {
	types := strings.Split(c.SubjectTypesSupported, ",")
	if len(types) == 0 {
		return []string{"public"}
	}
	return types
}

func (c *Config) GetClusterURLWithoutTailingSlashOrFail(cmd *cobra.Command) string {
	endpoint := c.GetClusterURLWithoutTailingSlash(cmd)
	if endpoint == "" {
		fmt.Println("To execute this command, the endpoint URL must point to the URL where ORY Hydra is located. To set the endpoint URL, use flag --endpoint or environment variable HYDRA_URL or HYDRA_ADMIN_URL if an administrative command was used.")
		os.Exit(1)
	}
	return endpoint
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

func (c *Config) GetTracer() (*tracing.Tracer, error) {
	if c.tracer == nil {
		c.GetLogger().Info("Setting up tracing middleware")

		c.tracer = &tracing.Tracer{
			ServiceName: c.TracingServiceName,
			JaegerConfig: &tracing.JaegerConfig{
				LocalAgentHostPort: c.JaegerLocalAgentHostPort,
				SamplerType:        c.JaegerSamplingType,
				SamplerValue:       c.JaegerSamplingValue,
				SamplerServerUrl:   c.JaegerSamplingServerUrl,
			},
			Provider: c.TracingProvider,
			Logger:   c.GetLogger(),
		}

		return c.tracer, c.tracer.Setup()
	}

	return c.tracer, nil
}

func (c *Config) WithTracing() bool {
	if tracer, err := c.GetTracer(); err == nil && tracer.IsLoaded() {
		return true
	} else {
		return false
	}
}

func (c *Config) GetPrometheusMetrics() *prometheus.MetricsManager {
	c.GetLogger().Info("Setting up Prometheus middleware")

	if c.prometheus == nil {
		c.prometheus = prometheus.NewMetricsManager(c.BuildVersion, c.BuildHash, c.BuildTime)
	}

	return c.prometheus
}

func (c *Config) DoesRequestSatisfyTermination(r *http.Request) error {
	if c.AllowTLSTermination == "" {
		return errors.New("TLS termination is not enabled")
	}

	if r.URL.Path == healthx.AliveCheckPath || r.URL.Path == healthx.ReadyCheckPath {
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

func (c *Config) GetLoginConsentRequestLifespan() time.Duration {
	d, err := time.ParseDuration(c.LoginConsentRequestLifespan)
	if err != nil {
		c.GetLogger().Warnf("Could not parse login and consent request lifespan value (%s). Defaulting to 15m", c.LoginConsentRequestLifespan)
		return time.Minute * 15
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

func (c *Config) GetRefreshTokenLifespan() time.Duration {
	if c.RefreshTokenLifespan == "-1" {
		return 0
	}

	d, err := time.ParseDuration(c.RefreshTokenLifespan)
	if err != nil {
		c.GetLogger().Warnf("Could not parse refresh token lifespan value (%s). Defaulting to 720h", c.RefreshTokenLifespan)
		return time.Hour * 720
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

	var hasher fosite.Hasher = &fosite.BCrypt{
		WorkFactor: c.BCryptWorkFactor,
	}

	if c.DatabaseURL == "" {
		c.GetLogger().Fatalf(`DATABASE_URL is not set, use "export DATABASE_URL=memory" for an in memory storage or the documented database adapters.`)
	} else if c.DatabasePlugin != "" {
		c.GetLogger().Infof("Database plugin set to %s", c.DatabasePlugin)
		pc := &PluginConnection{Config: c, Logger: c.GetLogger()}
		if err := pc.Load(); err != nil {
			c.GetLogger().Fatalf("Could not connect via database plugin: %s", err)
		}
	}

	var connection BackendConnector
	scheme := "memory"
	if c.DatabaseURL != "memory" {
		u, err := url.Parse(c.DatabaseURL)
		if err != nil {
			c.GetLogger().Fatalf("Could not parse DATABASE_URL: %s", err)
		}

		scheme = u.Scheme
	}

	if backend, ok := backends[scheme]; ok {
		options := []ConnectorOptions{}
		if c.WithTracing() {
			hasher = &tracing.TracedBCrypt{
				WorkFactor: c.BCryptWorkFactor,
			}
			options = append(options, WithTracing(), withOmitSQLArgsFromSpans())
		}

		if err := backend.Init(c.DatabaseURL, c.GetLogger(), options...); err != nil {
			c.GetLogger().Fatalf(`Could not connect to database backend: %s`, err)
		}
		connection = backend
	} else {
		c.GetLogger().Fatalf(`Unknown DSN scheme "%s" in DATABASE_URL "%s", schemes %v supported`, scheme, c.DatabaseURL, supportedSchemes())
	}

	c.context = &Context{
		Connection: connection,
		Hasher:     hasher,
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
		cmdx.Must(err, "Could not parse cluster url: %s", err)
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

func (c *Config) GetRotatedSystemSecrets() [][]byte {
	return [][]byte{
		pkg.HashStringSecret(c.RotatedSystemSecret),
	}
}

func (c *Config) GetSystemSecret() []byte {
	if len(c.systemSecret) > 0 {
		return c.systemSecret
	}

	if len(c.SystemSecret) >= 16 {
		c.systemSecret = pkg.HashStringSecret(c.SystemSecret)
		return pkg.HashStringSecret(c.SystemSecret)
	}

	if len(c.SystemSecret) > 0 {
		c.GetLogger().Fatalf("System secret must be undefined or have at least 16 characters, but it has %d characters.", len(c.SystemSecret))
		return nil
	}

	c.GetLogger().Warnf("No system secret was set, generating a random system secret...")
	secret, err := pkg.GenerateSecret(32)
	cmdx.Must(err, "Could not generate global secret: %s", err)
	c.GetLogger().Infof("Generated system secret: %s", secret)
	c.systemSecret = pkg.HashByteSecret(secret)
	c.GetLogger().Warnln("WARNING: DO NOT generate system secrets in production. The secret will be leaked to the logs.")
	return pkg.HashByteSecret(secret)
}

func (c *Config) GetFrontendAddress() string {
	return fmt.Sprintf("%s:%d", c.FrontendBindHost, c.FrontendBindPort)
}

func (c *Config) GetBackendAddress() string {
	return fmt.Sprintf("%s:%d", c.BackendBindHost, c.BackendBindPort)
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
