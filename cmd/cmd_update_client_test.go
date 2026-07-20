// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/cmd"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/configx"
	"github.com/ory/x/snapshotx"
	"github.com/ory/x/sqlxx"
)

// newUpdateClientCmd returns a fresh update command wired to the given admin
// endpoint. Each test case needs its own command because cobra remembers which
// flags were set across executions of the same command instance.
func newUpdateClientCmd(t *testing.T, adminURL string) *cobra.Command {
	c := cmd.NewUpdateClientCmd()
	cmdx.RegisterHTTPClientFlags(c.Flags())
	cmdx.RegisterFormatFlags(c.Flags())
	require.NoError(t, c.Flags().Set(cmdx.FlagEndpoint, adminURL))
	require.NoError(t, c.Flags().Set(cmdx.FlagFormat, string(cmdx.FormatJSON)))
	return c
}

// newTestClient returns a client with every flag-settable field set to a
// non-zero value, so that any field unintentionally dropped by an update is
// caught by comparing the persisted client before and after.
// sector_identifier_uri stays empty because a non-empty value makes the server
// fetch and validate the URL on every update.
func newTestClient() *client.Client {
	return &client.Client{
		AccessTokenStrategy:               "opaque",
		AllowedCORSOrigins:                sqlxx.StringSliceJSONFormat{"https://cors.example.com"},
		Audience:                          sqlxx.StringSliceJSONFormat{"https://api.example.com"},
		BackChannelLogoutSessionRequired:  true,
		BackChannelLogoutURI:              "https://example.com/bc-logout",
		ClientURI:                         "https://example.com",
		Contacts:                          sqlxx.StringSliceJSONFormat{"admin@example.com"},
		FrontChannelLogoutSessionRequired: true,
		FrontChannelLogoutURI:             "https://example.com/fc-logout",
		GrantTypes:                        sqlxx.StringSliceJSONFormat{"authorization_code", "refresh_token"},
		JSONWebKeysURI:                    "https://example.com/jwks.json",
		LogoURI:                           "https://example.com/logo.png",
		Metadata:                          sqlxx.JSONRawMessage(`{"foo":"bar"}`),
		Name:                              "original name",
		Owner:                             "original owner",
		PolicyURI:                         "https://example.com/policy",
		PostLogoutRedirectURIs:            sqlxx.StringSliceJSONFormat{"https://example.com/logged-out"},
		RedirectURIs:                      sqlxx.StringSliceJSONFormat{"https://example.com/callback"},
		RequestObjectSigningAlgorithm:     "RS256",
		RequestURIs:                       sqlxx.StringSliceJSONFormat{"https://example.com/request-object"},
		ResponseTypes:                     sqlxx.StringSliceJSONFormat{"code"},
		Scope:                             "openid offline_access",
		Secret:                            "original-secret",
		SkipConsent:                       true,
		SkipLogoutConsent:                 sqlxx.NullBool{Bool: true, Valid: true},
		SubjectType:                       "public",
		TokenEndpointAuthMethod:           "client_secret_basic",
		TermsOfServiceURI:                 "https://example.com/tos",
		// The validator normalizes an empty value to "none" on every update, so
		// set it here to keep before/after comparisons exact.
		UserinfoSignedResponseAlg: "none",
	}
}

// testClientResponseMatches asserts that the command output reflects the
// persisted client. The plaintext or hashed secret and the timestamps are
// excluded because they legitimately differ between response and storage.
func testClientResponseMatches(t *testing.T, stdout string, persisted *client.Client) {
	t.Helper()
	raw, err := json.Marshal(persisted)
	require.NoError(t, err)
	expected, actual := string(raw), stdout
	for _, volatile := range []string{"client_secret", "created_at", "updated_at"} {
		expected, err = sjson.Delete(expected, volatile)
		require.NoError(t, err)
		actual, err = sjson.Delete(actual, volatile)
		require.NoError(t, err)
	}
	// The persisted client serializes unset lifespans as null while the SDK
	// omits them; both mean "not set".
	for _, doc := range []*string{&expected, &actual} {
		for key, value := range gjson.Parse(*doc).Map() {
			if value.Type == gjson.Null {
				*doc, err = sjson.Delete(*doc, key)
				require.NoError(t, err)
			}
		}
	}
	assert.JSONEq(t, expected, actual)
}

func TestUpdateClientOnlyUpdatesProvidedFlags(t *testing.T) {
	reg := testhelpers.NewRegistryMemory(t, driver.WithConfigOptions(configx.WithValue(
		config.KeySubjectTypesSupported, []string{"public", "pairwise"},
	)))
	_, admin := testhelpers.NewOAuth2Server(t.Context(), t, reg)

	for _, tc := range []struct {
		flag   string
		args   []string
		expect func(c *client.Client)
	}{
		{
			flag:   "access-token-strategy",
			args:   []string{"--access-token-strategy", "jwt"},
			expect: func(c *client.Client) { c.AccessTokenStrategy = "jwt" },
		},
		{
			flag: "allowed-cors-origin",
			args: []string{"--allowed-cors-origin", "https://cors-2.example.com"},
			expect: func(c *client.Client) {
				c.AllowedCORSOrigins = sqlxx.StringSliceJSONFormat{"https://cors-2.example.com"}
			},
		},
		{
			flag:   "audience",
			args:   []string{"--audience", "https://api-2.example.com"},
			expect: func(c *client.Client) { c.Audience = sqlxx.StringSliceJSONFormat{"https://api-2.example.com"} },
		},
		{
			flag:   "backchannel-logout-callback",
			args:   []string{"--backchannel-logout-callback", "https://example.com/bc-logout-2"},
			expect: func(c *client.Client) { c.BackChannelLogoutURI = "https://example.com/bc-logout-2" },
		},
		{
			flag:   "backchannel-logout-session-required",
			args:   []string{"--backchannel-logout-session-required=false"},
			expect: func(c *client.Client) { c.BackChannelLogoutSessionRequired = false },
		},
		{
			flag:   "client-uri",
			args:   []string{"--client-uri", "https://example.com/about"},
			expect: func(c *client.Client) { c.ClientURI = "https://example.com/about" },
		},
		{
			flag:   "contact",
			args:   []string{"--contact", "ops@example.com,dev@example.com"},
			expect: func(c *client.Client) { c.Contacts = sqlxx.StringSliceJSONFormat{"ops@example.com", "dev@example.com"} },
		},
		{
			flag:   "frontchannel-logout-callback",
			args:   []string{"--frontchannel-logout-callback", "https://example.com/fc-logout-2"},
			expect: func(c *client.Client) { c.FrontChannelLogoutURI = "https://example.com/fc-logout-2" },
		},
		{
			flag:   "frontchannel-logout-session-required",
			args:   []string{"--frontchannel-logout-session-required=false"},
			expect: func(c *client.Client) { c.FrontChannelLogoutSessionRequired = false },
		},
		{
			flag: "grant-type",
			args: []string{"--grant-type", "authorization_code,client_credentials"},
			expect: func(c *client.Client) {
				c.GrantTypes = sqlxx.StringSliceJSONFormat{"authorization_code", "client_credentials"}
			},
		},
		{
			flag:   "jwks-uri",
			args:   []string{"--jwks-uri", "https://example.com/jwks-2.json"},
			expect: func(c *client.Client) { c.JSONWebKeysURI = "https://example.com/jwks-2.json" },
		},
		{
			flag:   "logo-uri",
			args:   []string{"--logo-uri", "https://example.com/logo-2.png"},
			expect: func(c *client.Client) { c.LogoURI = "https://example.com/logo-2.png" },
		},
		{
			flag:   "metadata",
			args:   []string{"--metadata", `{"widget":"sprocket"}`},
			expect: func(c *client.Client) { c.Metadata = sqlxx.JSONRawMessage(`{"widget":"sprocket"}`) },
		},
		{
			flag:   "name",
			args:   []string{"--name", "updated name"},
			expect: func(c *client.Client) { c.Name = "updated name" },
		},
		{
			flag:   "owner",
			args:   []string{"--owner", "updated owner"},
			expect: func(c *client.Client) { c.Owner = "updated owner" },
		},
		{
			flag:   "policy-uri",
			args:   []string{"--policy-uri", "https://example.com/policy-2"},
			expect: func(c *client.Client) { c.PolicyURI = "https://example.com/policy-2" },
		},
		{
			flag: "post-logout-callback",
			args: []string{"--post-logout-callback", "https://example.com/logged-out-2"},
			expect: func(c *client.Client) {
				c.PostLogoutRedirectURIs = sqlxx.StringSliceJSONFormat{"https://example.com/logged-out-2"}
			},
		},
		{
			flag: "redirect-uri",
			args: []string{"--redirect-uri", "https://example.com/callback-2,https://example.com/callback-3"},
			expect: func(c *client.Client) {
				c.RedirectURIs = sqlxx.StringSliceJSONFormat{"https://example.com/callback-2", "https://example.com/callback-3"}
			},
		},
		{
			flag:   "request-object-signing-alg",
			args:   []string{"--request-object-signing-alg", "ES256"},
			expect: func(c *client.Client) { c.RequestObjectSigningAlgorithm = "ES256" },
		},
		{
			flag: "request-uri",
			args: []string{"--request-uri", "https://example.com/request-object-2"},
			expect: func(c *client.Client) {
				c.RequestURIs = sqlxx.StringSliceJSONFormat{"https://example.com/request-object-2"}
			},
		},
		{
			flag:   "response-type",
			args:   []string{"--response-type", "code,id_token"},
			expect: func(c *client.Client) { c.ResponseTypes = sqlxx.StringSliceJSONFormat{"code", "id_token"} },
		},
		{
			flag:   "scope",
			args:   []string{"--scope", "read,write"},
			expect: func(c *client.Client) { c.Scope = "read write" },
		},
		{
			flag:   "skip-consent",
			args:   []string{"--skip-consent=false"},
			expect: func(c *client.Client) { c.SkipConsent = false },
		},
		{
			flag:   "skip-logout-consent",
			args:   []string{"--skip-logout-consent=false"},
			expect: func(c *client.Client) { c.SkipLogoutConsent = sqlxx.NullBool{Bool: false, Valid: true} },
		},
		{
			flag:   "subject-type",
			args:   []string{"--subject-type", "pairwise"},
			expect: func(c *client.Client) { c.SubjectType = "pairwise" },
		},
		{
			flag:   "token-endpoint-auth-method",
			args:   []string{"--token-endpoint-auth-method", "client_secret_post"},
			expect: func(c *client.Client) { c.TokenEndpointAuthMethod = "client_secret_post" },
		},
		{
			flag:   "tos-uri",
			args:   []string{"--tos-uri", "https://example.com/tos-2"},
			expect: func(c *client.Client) { c.TermsOfServiceURI = "https://example.com/tos-2" },
		},
	} {
		t.Run("flag="+tc.flag, func(t *testing.T) {
			original := createClient(t, reg, newTestClient())
			before, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
			require.NoError(t, err)

			stdout := cmdx.ExecNoErr(t, newUpdateClientCmd(t, admin.URL), append(append([]string{}, tc.args...), original.GetID())...)

			after, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
			require.NoError(t, err)

			expected := *before
			tc.expect(&expected)
			expected.UpdatedAt = after.UpdatedAt
			assert.Equal(t, &expected, after)

			testClientResponseMatches(t, stdout, after)
		})
	}

	t.Run("flag=sector-identifier-uri", func(t *testing.T) {
		// The server fetches and validates the sector identifier document, which
		// fails here (the URL is not HTTPS). The error proves the flag reached the
		// server under the right field name; the failed update must not change
		// the stored client.
		original := createClient(t, reg, newTestClient())
		before, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)

		_ = cmdx.ExecExpectedErr(t, newUpdateClientCmd(t, admin.URL), "--sector-identifier-uri", "http://example.com/sectors.json", original.GetID())

		after, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)
		assert.Equal(t, before, after)
	})

	t.Run("case=no flags leaves the client unchanged", func(t *testing.T) {
		original := createClient(t, reg, newTestClient())
		before, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)

		stdout := cmdx.ExecNoErr(t, newUpdateClientCmd(t, admin.URL), original.GetID())

		after, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)

		expected := *before
		expected.UpdatedAt = after.UpdatedAt
		assert.Equal(t, &expected, after)

		testClientResponseMatches(t, stdout, after)
	})

	t.Run("case=fails for an unknown client id", func(t *testing.T) {
		_ = cmdx.ExecExpectedErr(t, newUpdateClientCmd(t, admin.URL), "--name", "whatever", uuid.Must(uuid.NewV4()).String())
	})
}

func TestUpdateClientSecretHandling(t *testing.T) {
	reg := testhelpers.NewRegistryMemory(t)
	_, admin := testhelpers.NewOAuth2Server(t.Context(), t, reg)

	t.Run("case=keeps the secret when it is not provided", func(t *testing.T) {
		original := createClient(t, reg, newTestClient())
		before, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)

		stdout := cmdx.ExecNoErr(t, newUpdateClientCmd(t, admin.URL), "--name", "updated name", original.GetID())
		assert.False(t, gjson.Get(stdout, "client_secret").Exists(), "the secret must not be echoed when it did not change: %s", stdout)

		after, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)
		assert.Equal(t, before.Secret, after.Secret)
		require.NoError(t, reg.ClientHasher().Compare(t.Context(), []byte(after.Secret), []byte("original-secret")))
	})

	t.Run("case=updates the secret when provided", func(t *testing.T) {
		original := createClient(t, reg, newTestClient())
		before, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)

		stdout := cmdx.ExecNoErr(t, newUpdateClientCmd(t, admin.URL), "--secret", "new-secret", original.GetID())
		assert.Equal(t, "new-secret", gjson.Get(stdout, "client_secret").Str, "the new secret is echoed exactly once, in plaintext: %s", stdout)

		after, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)
		assert.NotEqual(t, before.Secret, after.Secret)
		require.NoError(t, reg.ClientHasher().Compare(t.Context(), []byte(after.Secret), []byte("new-secret")))

		// Everything except the secret is preserved.
		expected := *before
		expected.Secret = after.Secret
		expected.UpdatedAt = after.UpdatedAt
		assert.Equal(t, &expected, after)
	})

	t.Run("case=encrypts the returned secret when a pgp key is provided", func(t *testing.T) {
		original := createClient(t, reg, newTestClient())

		stdout := cmdx.ExecNoErr(t, newUpdateClientCmd(t, admin.URL),
			"--secret", "sooper-secret",
			"--pgp-key", base64EncodedPGPPublicKey(t),
			original.GetID(),
		)
		echoed := gjson.Get(stdout, "client_secret").Str
		assert.NotEmpty(t, echoed)
		assert.NotEqual(t, "sooper-secret", echoed, "the echoed secret must be encrypted")

		after, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)
		require.NoError(t, reg.ClientHasher().Compare(t.Context(), []byte(after.Secret), []byte("sooper-secret")),
			"the stored secret is the plaintext one, encryption only applies to the display")
	})

	t.Run("case=does not invent a secret when only a pgp key is provided", func(t *testing.T) {
		original := createClient(t, reg, newTestClient())
		before, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)

		stdout := cmdx.ExecNoErr(t, newUpdateClientCmd(t, admin.URL),
			"--name", "updated name",
			"--pgp-key", base64EncodedPGPPublicKey(t),
			original.GetID(),
		)
		assert.False(t, gjson.Get(stdout, "client_secret").Exists(), "no secret must be echoed when it did not change: %s", stdout)

		after, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)
		assert.Equal(t, before.Secret, after.Secret)
	})

	t.Run("case=does not echo an ignored secret flag in file mode", func(t *testing.T) {
		original := createClient(t, reg, newTestClient())

		// In file mode all other client flags are ignored, including --secret.
		// A secret that was never sent must not be echoed back.
		raw, err := json.Marshal(map[string]any{"client_name": "updated from file"})
		require.NoError(t, err)

		stdout, stderr, err := cmdx.Exec(t, newUpdateClientCmd(t, admin.URL), bytes.NewReader(raw),
			original.GetID(), "--file", "-", "--secret", "ignored-flag-secret")
		require.NoError(t, err, stderr)
		assert.False(t, gjson.Get(stdout, "client_secret").Exists(),
			"a secret that was not part of the request must not be echoed: %s", stdout)

		after, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)
		require.NoError(t, reg.ClientHasher().Compare(t.Context(), []byte(after.Secret), []byte("original-secret")),
			"the stored secret must be unchanged")
	})

	t.Run("case=rejects a secret that is too short", func(t *testing.T) {
		original := createClient(t, reg, newTestClient())
		before, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)

		_ = cmdx.ExecExpectedErr(t, newUpdateClientCmd(t, admin.URL), "--secret", "short", original.GetID())

		after, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)
		assert.Equal(t, before, after)
	})
}

func TestUpdateClient(t *testing.T) {
	t.Parallel()

	reg := testhelpers.NewRegistryMemory(t)
	_, admin := testhelpers.NewOAuth2Server(t.Context(), t, reg)

	original := createClient(t, reg, nil)
	t.Run("case=updates only the provided fields", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, newUpdateClientCmd(t, admin.URL), "--grant-type", "implicit", original.GetID()))
		expected, err := reg.ClientManager().GetClient(t.Context(), actual.Get("client_id").Str)
		require.NoError(t, err)

		assert.Equal(t, expected.GetID(), actual.Get("client_id").Str)
		assert.Equal(t, "implicit", actual.Get("grant_types").Array()[0].Str)
		assert.Equal(t, "client_secret_post", actual.Get("token_endpoint_auth_method").Str,
			"fields that were not passed as flags must keep their values")
		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})

	t.Run("case=supports encryption", func(t *testing.T) {
		actual := gjson.Parse(cmdx.ExecNoErr(t, newUpdateClientCmd(t, admin.URL),
			original.GetID(),
			"--secret", "some-userset-secret",
			"--pgp-key", base64EncodedPGPPublicKey(t),
		))
		assert.Equal(t, original.ID, actual.Get("client_id").Str)
		assert.NotEmpty(t, actual.Get("client_secret").Str)
		assert.NotEqual(t, original.Secret, actual.Get("client_secret").Str)

		snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
	})

	t.Run("case=updates from file", func(t *testing.T) {
		original, err := reg.ClientManager().GetConcreteClient(t.Context(), original.GetID())
		require.NoError(t, err)

		raw, err := json.Marshal(original)
		require.NoError(t, err)

		t.Run("file=stdin", func(t *testing.T) {
			raw, err = sjson.SetBytes(raw, "client_name", "updated through file stdin")
			require.NoError(t, err)

			stdout, stderr, err := cmdx.Exec(t, newUpdateClientCmd(t, admin.URL), bytes.NewReader(raw), original.GetID(), "--file", "-")
			require.NoError(t, err, stderr)

			actual := gjson.Parse(stdout)
			assert.Equal(t, original.ID, actual.Get("client_id").Str)
			assert.Equal(t, "updated through file stdin", actual.Get("client_name").Str)

			snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
		})

		t.Run("file=from disk", func(t *testing.T) {
			raw, err = sjson.SetBytes(raw, "client_name", "updated through file from disk")
			require.NoError(t, err)

			fn := writeTempFile(t, json.RawMessage(raw))

			stdout, stderr, err := cmdx.Exec(t, newUpdateClientCmd(t, admin.URL), nil, original.GetID(), "--file", fn)
			require.NoError(t, err, stderr)

			actual := gjson.Parse(stdout)
			assert.Equal(t, original.ID, actual.Get("client_id").Str)
			assert.Equal(t, "updated through file from disk", actual.Get("client_name").Str)

			snapshotx.SnapshotT(t, json.RawMessage(actual.Raw), snapshotExcludedClientFields...)
		})
	})
}
