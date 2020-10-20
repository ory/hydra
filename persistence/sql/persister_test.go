package sql_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/jwk"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/hydra/internal"
	"github.com/ory/viper"
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
		t.Run("package=client/manager="+k, func(t *testing.T) {
			t.Run("case=create-get-update-delete", client.TestHelperCreateGetUpdateDeleteClient(k, m.ClientManager()))

			t.Run("case=autogenerate-key", client.TestHelperClientAutoGenerateKey(k, m.ClientManager()))

			t.Run("case=auth-client", client.TestHelperClientAuthenticate(k, m.ClientManager()))
		})

		t.Run("package=consent/manager="+k, consent.ManagerTests(m.ConsentManager(), m.ClientManager(), m.OAuth2Storage()))

		t.Run("package=jwk/manager="+k, func(t *testing.T) {
			var testGenerator = &jwk.RS256Generator{}

			t.Run("TestManagerKey", func(t *testing.T) {
				ks, err := testGenerator.Generate("TestManagerKey", "sig")
				require.NoError(t, err)

				for name, r := range registries {
					t.Run(fmt.Sprintf("case=%s", name), jwk.TestHelperManagerKey(r.KeyManager(), ks, "TestManagerKey"))
				}
			})

			t.Run("TestManagerKeySet", func(t *testing.T) {
				ks, err := testGenerator.Generate("TestManagerKeySet", "sig")
				require.NoError(t, err)
				ks.Key("private")

				for name, r := range registries {
					t.Run(fmt.Sprintf("case=%s", name), jwk.TestHelperManagerKeySet(r.KeyManager(), ks, "TestManagerKeySet"))
				}
			})
		})
	}
}
