package cli

import (
	"context"
	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/x"
	"github.com/ory/x/josex"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_toSDKFriendlyJSONWebKey(t *testing.T) {
	publicJWK := []byte(`{
		"kty": "RSA",
		"e": "AQAB",
		"use": "sig",
		"kid": "7a5ff76a-6766-11ea-bc55-0242ac130003",
		"alg": "RS256",
		"n": "l80jJJqcc1PpefIGVIjuPvA1D7NscnuF9aQqLa7I9rDUK4IaSOO3kL_EF13k-jTzcA5q4OZn5dR0kmqIMZT2gQ"
	}`)

	publicPEM := []byte(`
		-----BEGIN PUBLIC KEY-----
		MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAPf64dykufSkwnvUiBAwd5Si0K6t4m5i
		qJD8TmLJCmFjKaOUa6nszcFt/FkAuORfdlrD9mEZLPrPx74RSluyTBMCAwEAAQ==
		-----END PUBLIC KEY-----
	`)

	type args struct {
		key []byte
		kid string
		use string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "JWK with algorithm",
			args: args{
				key: publicJWK,
				kid: "public:7a5ff76a-6766-11ea-bc55-0242ac130003",
				use: "sig",
			},
			want: "RS256",
		},
		{
			name: "PEM key without algorithm",
			args: args{
				key: publicPEM,
				kid: "public:7a5ff76a-6766-11ea-bc55-0242ac130003",
				use: "sig",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, _ := josex.LoadPublicKey(tt.args.key)
			if got := toSDKFriendlyJSONWebKey(key, tt.args.kid, tt.args.use); got.Algorithm != tt.want {
				t.Errorf("toSDKFriendlyJSONWebKey() = %v, want %v", got.Algorithm, tt.want)
			}
		})
	}

	conf := internal.NewConfigurationWithDefaults()
	reg := internal.NewRegistryMemory(t, conf)
	router := x.NewRouterPublic()

	h := reg.KeyHandler()
	m := reg.KeyManager()
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
		NewHandler().Keys.ImportKeys(&cmd, []string{"setName", "../test/private_key.json", "../test/public_key.json"})
		//running again to make sure the row in storage is not deleted issue: #2436
		NewHandler().Keys.ImportKeys(&cmd, []string{"setName", "../test/private_key.json", "../test/public_key.json"})
		v, _ := m.GetKeySet(context.TODO(), "setName")
		assert.NotEmpty(t, v.Keys[0])
	})

}
