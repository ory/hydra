package internal

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/ory/viper"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/jwk"
	"github.com/ory/x/logrusx"
)

func resetConfig() {
	viper.Set(configuration.ViperKeyWellKnownKeys, nil)
	viper.Set(configuration.ViperKeySubjectTypesSupported, nil)
	viper.Set(configuration.ViperKeyDefaultClientScope, nil)
	viper.Set(configuration.ViperKeyDSN, nil)
	viper.Set(configuration.ViperKeyEncryptSessionData, true)
	viper.Set(configuration.ViperKeyBCryptCost, nil)
	viper.Set(configuration.ViperKeyAdminListenOnHost, nil)
	viper.Set(configuration.ViperKeyAdminListenOnPort, nil)
	viper.Set(configuration.ViperKeyPublicListenOnHost, nil)
	viper.Set(configuration.ViperKeyPublicListenOnPort, nil)
	viper.Set(configuration.ViperKeyConsentRequestMaxAge, nil)
	viper.Set(configuration.ViperKeyAccessTokenLifespan, nil)
	viper.Set(configuration.ViperKeyRefreshTokenLifespan, nil)
	viper.Set(configuration.ViperKeyIDTokenLifespan, nil)
	viper.Set(configuration.ViperKeyAuthCodeLifespan, nil)
	viper.Set(configuration.ViperKeyScopeStrategy, nil)
	viper.Set(configuration.ViperKeyGetCookieSecrets, nil)
	viper.Set(configuration.ViperKeyGetSystemSecret, nil)
	viper.Set(configuration.ViperKeyLogoutRedirectURL, nil)
	viper.Set(configuration.ViperKeyLoginURL, nil)
	viper.Set(configuration.ViperKeyConsentURL, nil)
	viper.Set(configuration.ViperKeyErrorURL, nil)
	viper.Set(configuration.ViperKeyPublicURL, nil)
	viper.Set(configuration.ViperKeyIssuerURL, nil)
	viper.Set(configuration.ViperKeyOAuth2ClientRegistrationURL, nil)
	viper.Set(configuration.ViperKeyAllowTLSTerminationFrom, nil)
	viper.Set(configuration.ViperKeyAccessTokenStrategy, nil)
	viper.Set(configuration.ViperKeySubjectIdentifierAlgorithmSalt, nil)
	viper.Set(configuration.ViperKeyOIDCDiscoverySupportedClaims, nil)
	viper.Set(configuration.ViperKeyOIDCDiscoverySupportedScope, nil)
	viper.Set(configuration.ViperKeyOIDCDiscoveryUserinfoEndpoint, nil)

	viper.Set(configuration.ViperKeyBCryptCost, "4")
	viper.Set(configuration.ViperKeySubjectIdentifierAlgorithmSalt, "00000000")
	viper.Set(configuration.ViperKeyGetSystemSecret, []string{"000000000000000000000000000000000000000000000000"})
	viper.Set(configuration.ViperKeyGetCookieSecrets, []string{"000000000000000000000000000000000000000000000000"})

	viper.Set("LOG_LEVEL", "debug")
}

func NewConfigurationWithDefaults() *configuration.ViperProvider {
	resetConfig()
	return configuration.NewViperProvider(logrusx.New(), true, nil).(*configuration.ViperProvider)
}

func NewConfigurationWithDefaultsAndHTTPS() *configuration.ViperProvider {
	resetConfig()
	return configuration.NewViperProvider(logrusx.New(), false, nil).(*configuration.ViperProvider)
}

func NewRegistry(c *configuration.ViperProvider) *driver.RegistryMemory {
	viper.Set("LOG_LEVEL", "debug")
	r := driver.NewRegistryMemory().WithConfig(c)
	_ = r.Init()
	return r.(*driver.RegistryMemory)
}

func NewRegistrySQL(c *configuration.ViperProvider, db *sqlx.DB) *driver.RegistrySQL {
	viper.Set("LOG_LEVEL", "debug")
	r := driver.NewRegistrySQL().WithConfig(c).(*driver.RegistrySQL).WithDB(db)
	_ = r.Init()
	return r.(*driver.RegistrySQL)
}

func MustEnsureRegistryKeys(r driver.Registry, key string) {
	if err := jwk.EnsureAsymmetricKeypairExists(context.Background(), r, new(veryInsecureRS256Generator), key); err != nil {
		panic(err)
	}
}
