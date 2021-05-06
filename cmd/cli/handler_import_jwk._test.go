package cli

import (
	"context"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
)

func Test_toImportJSONWebKey(t *testing.T) {
	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, conf)
	router := x.NewRouterPublic()

	h := reg.KeyHandler()
	m := reg.KeyManager()
	if os.Getenv("IMPORT_ERROR") == "1" {

		h.SetRoutes(router.RouterAdmin(), router, func(h http.Handler) http.Handler {
			return h
		})
		testServer := httptest.NewServer(router)

		cmd := cobra.Command{
			Use: "key",
		}
		cmd.Flags().String("use", "sig", "Sets the \"use\" value of the JSON Web Key if not \"use\" value was defined by the key itself")
		cmd.Flags().Bool("fake-tls-termination", false, "Sets the \"use\" value of the JSON Web Key if not \"use\" value was defined by the key itself")
		cmd.Flags().String("access-token", "", "Set an access token to be used in the Authorization header, defaults to environment variable OAUTH2_ACCESS_TOKEN")
		cmd.Flags().String("endpoint", "", "Set the URL where ORY Hydra is hosted, defaults to environment variable HYDRA_ADMIN_URL. A unix socket can be set in the form unix:///path/to/socket")
		cmd.Flags().Bool("skip-tls-verify", true, "Foolishly accept TLS certificates signed by unknown certificate authorities")
		os.Setenv("HYDRA_URL", testServer.URL)

		t.Run("Test_ImportKeys/Run_multiple_time_With_same_Values", func(t *testing.T) {
			NewHandler().Keys.ImportKeys(&cmd, []string{"import-1", "../test/private_key.json", "../test/public_key.json"})
			//running again to make sure the row in storage is not deleted issue: #2436
			NewHandler().Keys.ImportKeys(&cmd, []string{"import-1", "../test/private_key.json", "../test/public_key.json"})

		})
		return
	}
	//code was added to catch the os.Exit(1) in ImportKeys()
	testCmd := exec.Command(os.Args[0], "-test.run=Test_toImportJSONWebKey")
	testCmd.Env = append(os.Environ(), "IMPORT_ERROR=1")
	err := testCmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		v, _ := m.GetKeySet(context.TODO(), "import-1")
		assert.NotEmpty(t, v)
	}
}
