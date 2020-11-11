package testhelpers

import (
	"net/http"
	"testing"
	"time"

	"github.com/gobuffalo/httptest"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
	"github.com/ory/viper"
)

func NoOpMiddleware(h http.Handler) http.Handler {
	return h
}

func NewOAuth2Server(t *testing.T, reg driver.Registry) (publicTS, adminTS *httptest.Server) {
	// Lifespan is two seconds to avoid time synchronization issues with SQL.
	viper.Set(configuration.ViperKeySubjectIdentifierAlgorithmSalt, "76d5d2bf-747f-4592-9fbd-d2b895a54b3a")
	viper.Set(configuration.ViperKeyAccessTokenLifespan, time.Second*2)
	viper.Set(configuration.ViperKeyRefreshTokenLifespan, time.Second*3)
	viper.Set(configuration.ViperKeyScopeStrategy, "exact")

	public, admin := x.NewRouterPublic(), x.NewRouterAdmin()

	publicTS = httptest.NewServer(public)
	t.Cleanup(publicTS.Close)

	adminTS = httptest.NewServer(admin)
	t.Cleanup(adminTS.Close)

	viper.Set(configuration.ViperKeyIssuerURL, publicTS.URL)
	// SendDebugMessagesToClients: true,

	internal.MustEnsureRegistryKeys(reg, x.OpenIDConnectKeyName)
	if reg.Config().AccessTokenStrategy() == "jwt" {
		internal.MustEnsureRegistryKeys(reg, x.OAuth2JWTKeyName)
	}

	reg.RegisterRoutes(admin, public)
	return publicTS, adminTS
}
