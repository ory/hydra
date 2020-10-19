package sql_test

import (
	"fmt"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"
	"github.com/ory/viper"
	"testing"
	"time"
)

func TestManagers(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	viper.Set(configuration.ViperKeyAccessTokenLifespan, time.Hour)
	registries := map[string]driver.Registry{
		"memory": internal.NewRegistryMemory(t, conf),
	}

	if !testing.Short() {
		registries["postgres"], registries["mysql"], registries["cockroach"], _ = internal.ConnectDatabases(t)
	}

	for k, m := range registries {
		t.Run("package=client", func(t *testing.T) {
			t.Run("case=create-get-update-delete", func(t *testing.T) {
				t.Run(fmt.Sprintf("db=%s", k), client.TestHelperCreateGetUpdateDeleteClient(k, m.ClientManager()))
			})

			t.Run("case=autogenerate-key", func(t *testing.T) {
				t.Run(fmt.Sprintf("db=%s", k), client.TestHelperClientAutoGenerateKey(k, m.ClientManager()))
			})

			t.Run("case=auth-client", func(t *testing.T) {
				t.Run(fmt.Sprintf("db=%s", k), client.TestHelperClientAuthenticate(k, m.ClientManager()))
			})
		})

		t.Run("package=consent/manager="+k, consent.ManagerTests(m.ConsentManager(), m.ClientManager(), m.OAuth2Storage()))
	}
}
