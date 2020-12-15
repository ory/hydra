package sql_test

import (
	"testing"
	"time"

	"github.com/pborman/uuid"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/jwk"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
)

func TestManagers(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	conf.MustSet(config.KeyAccessTokenLifespan, time.Hour)
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

			t.Run("case=update-two-clients", client.TestHelperUpdateTwoClients(k, m.ClientManager()))
		})

		t.Run("package=consent/manager="+k, consent.ManagerTests(m.ConsentManager(), m.ClientManager(), m.OAuth2Storage()))

		t.Run("package=jwk/manager="+k, func(t *testing.T) {
			var testGenerator = &jwk.RS256Generator{}

			t.Run("TestManagerKey", func(t *testing.T) {
				ks, err := testGenerator.Generate("TestManagerKey", "sig")
				require.NoError(t, err)

				jwk.TestHelperManagerKey(m.KeyManager(), ks, uuid.New())
			})

			t.Run("TestManagerKeySet", func(t *testing.T) {
				ks, err := testGenerator.Generate("TestManagerKeySet", "sig")
				require.NoError(t, err)
				ks.Key("private")

				jwk.TestHelperManagerKeySet(m.KeyManager(), ks, uuid.New())
			})
		})
	}
}
