package configuration

import (
	"fmt"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/tracing"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/stringslice"
	"github.com/ory/x/urlx"
	"github.com/ory/x/viperx"
)

type ViperProvider struct {
	l               logrus.FieldLogger
	ss              [][]byte
	generatedSecret []byte
}

func NewViperProvider(l logrus.FieldLogger) Provider {
	return &ViperProvider{
		l: l,
	}
}

func (v *ViperProvider) WellKnownKeys(include ...string) []string {
	return append(viperx.GetStringSlice(v.l, "oidc.jwks.publish", []string{}), include...)
}

func (v *ViperProvider) SubjectTypesSupported() []string {
	types := stringslice.Filter(
		viperx.GetStringSlice(v.l,
			"oidc.subject_identifiers.enabled",
			[]string{"public"},
			"OIDC_SUBJECT_TYPES_SUPPORTED",
		),
		func(s string) bool {
			return !(s == "public" || s == "pairwise")
		},
	)

	if stringslice.Has(types, "pairwise") {
		if v.AccessTokenStrategy() == "jwt" {
			v.l.Fatalf(`The pairwise subject identifier algorithm is not supported by the JWT OAuth 2.0 Access Token Strategy. Please remove "pairwise" from oidc.subject_identifiers.supported or set strategies.access_token to "opaque".`)
		}
		if len(v.SubjectIdentifierAlgorithmSalt()) < 8 {
			v.l.Fatalf(`The pairwise subject identifier algorithm was set but length of oidc.subject_identifier.salt is too small (%d < 8), please set oidc.subject_identifiers.pairwise.salt to a random string with 8 characters or more.`, len(v.SubjectIdentifierAlgorithmSalt()))
		}
	}

	return types
}

func (v *ViperProvider) DefaultClientScope() []string {
	return viperx.GetStringSlice(v.l,
		"oidc.dynamic_client_registration.default_scope",
		[]string{"offline_access", "offline", "openid"},
		"OIDC_DYNAMIC_CLIENT_REGISTRATION_DEFAULT_SCOPE",
	)
}

func (v *ViperProvider) DSN() string {
	return viperx.GetString(v.l, "dsn", "", "DATABASE_URL")
}

func (v *ViperProvider) DataSourcePlugin() string {
	return viperx.GetString(v.l, "driver.plugin_path", "", "DATABASE_PLUGIN")
}

func (v *ViperProvider) BCryptCost() int {
	return viperx.GetInt(v.l, "hashers.bcrypt.cost", 10, "BCRYPT_COST")
}

func (v *ViperProvider) AdminListenOn() string {
	host := viperx.GetString(v.l, "httpd.admin.host", "", "ADMIN_HOST")
	port := viperx.GetInt(v.l, "httpd.admin.port", 4445, "ADMIN_PORT")
	return fmt.Sprintf("%s:%d", host, port)
}

func (v *ViperProvider) PublicListenOn() string {
	host := viperx.GetString(v.l, "httpd.public.host", "", "PUBLIC_HOST")
	port := viperx.GetInt(v.l, "httpd.public.port", 4444, "PUBLIC_PORT")
	return fmt.Sprintf("%s:%d", host, port)
}

func (v *ViperProvider) ConsentRequestMaxAge() time.Duration {
	return viperx.GetDuration(v.l, "ttl.login_consent_request", time.Minute*30, "LOGIN_CONSENT_REQUEST_LIFESPAN")
}

func (v *ViperProvider) AccessTokenLifespan() time.Duration {
	return viperx.GetDuration(v.l, "ttl.access_token", time.Minute*30, "ACCESS_TOKEN_LIFESPAN")
}

func (v *ViperProvider) RefreshTokenLifespan() time.Duration {
	return viperx.GetDuration(v.l, "ttl.refresh_token", -1, "REFRESH_TOKEN_LIFESPAN")
}

func (v *ViperProvider) IDTokenLifespan() time.Duration {
	return viperx.GetDuration(v.l, "ttl.id_token", time.Hour, "ID_TOKEN_LIFESPAN")
}

func (v *ViperProvider) AuthCodeLifespan() time.Duration {
	return viperx.GetDuration(v.l, "ttl.auth_code", time.Minute*10, "AUTH_CODE_LIFESPAN")
}

func (v *ViperProvider) ScopeStrategy() string {
	return viperx.GetString(v.l, "strategies.scope", "", "SCOPE_STRATEGY")
}

func (v *ViperProvider) TracingServiceName() string {
	return viperx.GetString(v.l, "tracing.service_name", "ORY Hydra")
}

func (v *ViperProvider) TracingProvider() string {
	return viperx.GetString(v.l, "tracing.provider", "", "TRACING_PROVIDER")
}

func (v *ViperProvider) TracingJaegerConfig() *tracing.JaegerConfig {
	return &tracing.JaegerConfig{
		LocalAgentHostPort: viperx.GetString(v.l, "tracing.providers.jaeger.local_agent_address", "", "TRACING_PROVIDER_JAEGER_LOCAL_AGENT_ADDRESS"),
		SamplerType:        viperx.GetString(v.l, "tracing.providers.jaeger.sampling.type", "", "TRACING_PROVIDER_JAEGER_SAMPLING_TYPE"),
		SamplerValue:       viperx.GetFloat64(v.l, "tracing.providers.jaeger.sampling.value", 0, "TRACING_PROVIDER_JAEGER_SAMPLING_VALUE"),
		SamplerServerUrl:   viperx.GetString(v.l, "tracing.providers.jaeger.sampling.server_url", "", "TRACING_PROVIDER_JAEGER_SAMPLING_SERVER_URL"),
	}
}

func (v *ViperProvider) GetCookieSecrets() [][]byte {
	return [][]byte{
		[]byte(viperx.GetString(v.l, "secrets.cookie", string(v.GetSystemSecret()), "COOKIE_SECRET")),
	}
}

func (v *ViperProvider) GetRotatedSystemSecrets() [][]byte {
	secrets := viperx.GetStringSlice(v.l, "secrets.system", []string{}, "ROTATED_SYSTEM_SECRET")

	if len(secrets) < 2 {
		return nil
	}

	var rotated [][]byte
	for _, secret := range secrets[1:] {
		rotated = append(rotated, pkg.HashStringSecret(secret))
	}

	return rotated
}

func (v *ViperProvider) GetSystemSecret() []byte {
	secrets := viperx.GetStringSlice(v.l, "secrets.system", []string{}, "SYSTEM_SECRET")

	if len(secrets) == 0 {
		if v.generatedSecret != nil {
			return v.generatedSecret
		}

		v.l.Warnf("Configuration secrets.system is not set, generating a temporary, random secret...")
		secret, err := pkg.GenerateSecret(32)
		cmdx.Must(err, "Could not generate secret: %s", err)

		v.l.Warnf("Generated secret: %s", secret)
		v.generatedSecret = pkg.HashByteSecret(secret)

		v.l.Warnln("Do not use generate secrets in production. The secret will be leaked to the logs.")
		return pkg.HashByteSecret(secret)
	}

	secret := secrets[0]
	if len(secret) >= 16 {
		return pkg.HashStringSecret(secret)
	}

	v.l.Fatalf("System secret must be undefined or have at least 16 characters but only has %d characters.", len(secret))
	return nil
}

func (v *ViperProvider) LogoutRedirectURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, "urls.post_logout_redirect", "", "OAUTH2_LOGOUT_REDIRECT_URL"))
}

func (v *ViperProvider) LoginURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, "urls.logout", "", "OAUTH2_LOGIN_URL"))
}
func (v *ViperProvider) ConsentURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, "urls.consent", "", "OAUTH2_CONSENT_URL"))
}

func (v *ViperProvider) ErrorURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, "urls.error", "", "OAUTH2_ERROR_URL"))
}

func (v *ViperProvider) PublicURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, "urls.self.public", ""))
}

func (v *ViperProvider) AdminURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, "urls.self.public", ""))
}

func (v *ViperProvider) IssuerURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, "urls.self.issuer", v.PublicURL().String(), "OAUTH2_ISSUER_URL", "ISSUER", "ISSUER_URL"))
}

func (v *ViperProvider) OAuth2AuthURL() string {
	return urlx.MustJoin(v.PublicURL().String(), "/oauth2/auth")
}

func (v *ViperProvider) OAuth2ClientRegistrationURL() *url.URL {
	return urlx.ParseOrFatal(v.l, viperx.GetString(v.l, "oidc.discovery.client_registration_url", "", "OAUTH2_CLIENT_REGISTRATION_URL"))
}

func (v *ViperProvider) AllowTLSTerminationFrom() []string {
	return viperx.GetStringSlice(v.l, "httpd.tls.allow_termination_from", []string{}, "HTTPS_ALLOW_TERMINATION_FROM")
}

func (v *ViperProvider) AccessTokenStrategy() string {
	return viperx.GetString(v.l, "strategies.access_token", "opaque", "OAUTH2_ACCESS_TOKEN_STRATEGY")
}

func (v *ViperProvider) SubjectIdentifierAlgorithmSalt() string {
	return viperx.GetString(v.l, "oidc.subject_identifiers.pairwise.salt", "", "OIDC_SUBJECT_TYPE_PAIRWISE_SALT")
}

func (v *ViperProvider) OIDCDiscoverySupportedClaims() []string {
	return viperx.GetStringSlice(v.l, "oidc.discovery.supported_claims", []string{}, "OIDC_DISCOVERY_CLAIMS_SUPPORTED")
}

func (v *ViperProvider) OIDCDiscoverySupportedScope() []string {
	return viperx.GetStringSlice(v.l, "oidc.discovery.supported_scope", []string{}, "OIDC_DISCOVERY_SCOPES_SUPPORTED")
}

func (v *ViperProvider) OIDCDiscoveryUserinfoEndpoint() string {
	return viperx.GetString(v.l, "oidc.discovery.userinfo_url", urlx.AppendPaths(v.PublicURL(), "/userinfo").String(), "OIDC_DISCOVERY_USERINFO_ENDPOINT")
}

func (v *ViperProvider) ShareOAuth2Debug() bool {
	return viperx.GetBool(v.l, "debug.share_oauth2_errors", "OAUTH2_SHARE_ERROR_DEBUG")
}

//func (v *ViperProvider) ServesHTTPS() bool {
//	return
//}
//func (v *ViperProvider) HashSignature() bool {
//	return
//}
//func (v *ViperProvider) IsUsingJWTAsAccessTokens() bool {
//	return
//}
