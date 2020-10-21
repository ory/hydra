package clients

import (
	"context"
	"encoding/json"
	"github.com/ory/hydra/cmd/cli"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
	"github.com/ory/x/cmdx"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"net/http"
	"net/http/httptest"
	"testing"
)

func noopCors(handler http.Handler) http.Handler {
	return handler
}

func setupTestCmd(t *testing.T, cmd *cobra.Command) (driver.Registry, func()) {
	c := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(c)
	ra, rp := x.NewRouterAdmin(), x.NewRouterPublic()

	reg.ClientHandler().SetRoutes(ra)
	reg.ConsentHandler().SetRoutes(ra)
	reg.KeyHandler().SetRoutes(ra, rp, noopCors)
	reg.OAuth2Handler().SetRoutes(ra, rp, noopCors)

	s := httptest.NewServer(ra)

	require.NoError(t, cmd.Flags().Set(cli.FlagAdminEndpoint, s.URL))

	return reg, s.Close
}

func requireMarshaledJSON(t *testing.T, v interface{}) string {
	r, err := json.Marshal(v)
	require.NoError(t, err, "%+v", v)
	return string(r)
}

func assertContainedInJSON(t *testing.T, obj map[string]interface{}, j string) {
	for k, v := range obj {
		assert.Equal(t, requireMarshaledJSON(t, v), gjson.Get(j, k).Raw)
	}
}

func assertPartialClient(t *testing.T, inputClient map[string]interface{}, stdOut string, reg driver.Registry) {
	newID := gjson.Get(stdOut, "client_id").Str

	// check if client got created
	c, err := reg.ClientManager().GetClient(context.Background(), newID)
	assert.NoError(t, err, stdOut)

	assertContainedInJSON(t, inputClient, stdOut)

	// secret is not available from the store as it is bcrypted
	delete(inputClient, "client_secret")
	assertContainedInJSON(t, inputClient, requireMarshaledJSON(t, c))
}

func newTestSuite(t *testing.T, newCmd func() *cobra.Command) func(name string, runner func(t *testing.T, cmd *cobra.Command, reg driver.Registry)) {
	return func(name string, runner func(t *testing.T, cmd *cobra.Command, reg driver.Registry)) {
		cmd := newCmd()

		reg, cleanup := setupTestCmd(t, cmd)
		defer cleanup()

		require.NoError(t, cmd.Flags().Set(cmdx.FlagFormat, string(cmdx.FormatJSON)))

		t.Run(name, func(t *testing.T) {
			runner(t, cmd, reg)
		})
	}
}
