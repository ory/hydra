package sql_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/pborman/uuid"

	"github.com/ory/hydra/oauth2/trust"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/internal/testhelpers"

	"github.com/ory/hydra/jwk"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/internal"
)

func TestManagers(t *testing.T) {
	registries := map[string]driver.Registry{
		"memory": internal.NewRegistryMemory(t, internal.NewConfigurationWithDefaults()),
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

		t.Run("package=consent/janitor="+k, testhelpers.JanitorTests(m.Config(), m.ConsentManager(), m.ClientManager(), m.OAuth2Storage()))

		t.Run("package=jwk/manager="+k, func(t *testing.T) {
			keyGenerators := new(driver.RegistryBase).KeyGenerators()
			assert.Equalf(t, 6, len(keyGenerators), "Test for key generator is not implemented")

			for _, tc := range []struct {
				keyGenerator jwk.KeyGenerator
				alg          string
				skip         bool
			}{
				{keyGenerator: keyGenerators["RS256"], alg: "RS256", skip: false},
				{keyGenerator: keyGenerators["ES256"], alg: "ES256", skip: false},
				{keyGenerator: keyGenerators["ES512"], alg: "ES512", skip: false},
				{keyGenerator: keyGenerators["HS256"], alg: "HS256", skip: true},
				{keyGenerator: keyGenerators["HS512"], alg: "HS512", skip: true},
				{keyGenerator: keyGenerators["EdDSA"], alg: "EdDSA", skip: m.Config().HsmEnabled()},
			} {
				t.Run("key_generator="+tc.alg, func(t *testing.T) {
					if tc.skip {
						t.Skipf("Skipping test. Not applicable for alg: %s", tc.alg)
					}
					if m.Config().HsmEnabled() {
						t.Run("TestManagerGenerateAndPersistKeySet", jwk.TestHelperManagerGenerateAndPersistKeySet(m.KeyManager(), tc.alg))
					} else {
						kid := uuid.New()
						ks, err := tc.keyGenerator.Generate(kid, "sig")
						require.NoError(t, err)
						t.Run("TestManagerKey", jwk.TestHelperManagerKey(m.KeyManager(), tc.alg, ks, kid))
						t.Run("TestManagerKeySet", jwk.TestHelperManagerKeySet(m.KeyManager(), tc.alg, ks, kid))
						t.Run("TestManagerGenerateAndPersistKeySet", jwk.TestHelperManagerGenerateAndPersistKeySet(m.KeyManager(), tc.alg))
					}
				})
			}

			t.Run("TestManagerGenerateAndPersistKeySetWithUnsupportedKeyGenerator", func(t *testing.T) {
				_, err := m.KeyManager().GenerateAndPersistKeySet(context.TODO(), "foo", "bar", "UNKNOWN", "sig")
				require.Error(t, err)
				assert.IsType(t, errors.WithStack(jwk.ErrUnsupportedKeyAlgorithm), err)
			})
		})

		t.Run("package=grant/trust/manager="+k, func(t *testing.T) {
			t.Run("case=create-get-delete", trust.TestHelperGrantManagerCreateGetDeleteGrant(m.GrantManager()))
			t.Run("case=errors", trust.TestHelperGrantManagerErrors(m.GrantManager()))
		})
	}
}
