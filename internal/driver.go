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
	viper.Reset()

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
