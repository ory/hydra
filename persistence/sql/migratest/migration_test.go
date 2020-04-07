package migratest

import (
	"fmt"
	"github.com/gobuffalo/pop/v5"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/x"
	"github.com/ory/x/resilience"
	"github.com/ory/x/sqlcon/dockertest"
	"github.com/ory/x/sqlxx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func assertEqualClients(t *testing.T, expected, actual *client.Client) {
	now := time.Now()
	expected.CreatedAt = now
	expected.UpdatedAt = now
	actual.CreatedAt = now
	actual.UpdatedAt = now

	assert.Equal(t, expected, actual)
}

func generateExpectedClient(migration int) *client.Client {
	//return &client.Client{
	//	PK:                                int64(migration),
	//	ClientID:                          fmt.Sprintf("client-%04d", migration),
	//	Name:                              fmt.Sprintf("Client %04d", migration),
	//	Secret:                            fmt.Sprintf("secret-%04d", migration),
	//	RedirectURIs:                      []string{fmt.Sprintf("http://redirect/%04d_1", migration)},
	//	GrantTypes:                        []string{fmt.Sprintf("grant-%04d_1", migration)},
	//	ResponseTypes:                     []string{fmt.Sprintf("response-%04d_1", migration)},
	//	Scope:                             fmt.Sprintf("scope-%04d", migration),
	//	Audience:                          []string{fmt.Sprintf("autdience-%04d_1", migration)},
	//	Owner:                             fmt.Sprintf("owner-%04d", migration),
	//	PolicyURI:                         fmt.Sprintf("http://policy/%04d", migration),
	//	AllowedCORSOrigins:                []string{fmt.Sprintf("http://cors/%04d_1", migration)},
	//	TermsOfServiceURI:                 fmt.Sprintf("http://tos/%04d", migration),
	//	ClientURI:                         fmt.Sprintf("http://client/%04d", migration),
	//	LogoURI:                           fmt.Sprintf("http://logo/%04d", migration),
	//	Contacts:                          []string{fmt.Sprintf("contact-%04d_1", migration)},
	//	SecretExpiresAt:                   0,
	//	SubjectType:                       fmt.Sprintf("subject-%04d", migration),
	//	SectorIdentifierURI:               fmt.Sprintf("http://sector_id/%04d", migration),
	//	JSONWebKeysURI:                    fmt.Sprintf("http://jwks/%04d", migration),
	//	JSONWebKeys:                       &x.JoseJSONWebKeySet{},
	//	TokenEndpointAuthMethod:           fmt.Sprintf("token_auth-%04d", migration),
	//	RequestURIs:                       []string{fmt.Sprintf("http://request/%04d_1", migration)},
	//	RequestObjectSigningAlgorithm:     fmt.Sprintf("request_alg-%04d", migration),
	//	UserinfoSignedResponseAlg:         fmt.Sprintf("userinfo_alg-%04d", migration),
	//	FrontChannelLogoutURI:             fmt.Sprintf("http://front_logout/%04d", migration),
	//	FrontChannelLogoutSessionRequired: false,
	//	PostLogoutRedirectURIs:            []string{fmt.Sprintf("http://post_redirect/%04d_1", migration)},
	//	BackChannelLogoutURI:              fmt.Sprintf("http://back_logout/%04d", migration),
	//	BackChannelLogoutSessionRequired:  false,
	//	Metadata:                          sqlxx.JSONRawMessage("{}"),
	//}
	return &client.Client{
		PK:                      int64(migration),
		ClientID:                fmt.Sprintf("%d-data", migration),
		Name:                    "some-client",
		Secret:                  "abcdef",
		RedirectURIs:            sqlxx.StringSlicePipeDelimiter{"http://localhost", "http://google"},
		GrantTypes:              sqlxx.StringSlicePipeDelimiter{"authorize_code", "implicit"},
		ResponseTypes:           sqlxx.StringSlicePipeDelimiter{"token", "id_token"},
		Scope:                   "foo|bar",
		Owner:                   "aeneas",
		PolicyURI:               "http://policy",
		TermsOfServiceURI:       "http://tos",
		ClientURI:               "http://client",
		LogoURI:                 "http://logo",
		Contacts:                sqlxx.StringSlicePipeDelimiter{"aeneas", "foo"},
		JSONWebKeys:             &x.JoseJSONWebKeySet{},
		TokenEndpointAuthMethod: "none",
		Metadata:                sqlxx.JSONRawMessage("{}"),
		Audience:                sqlxx.StringSlicePipeDelimiter{},
		AllowedCORSOrigins:      sqlxx.StringSlicePipeDelimiter{},
		RequestURIs:             sqlxx.StringSlicePipeDelimiter{},
		PostLogoutRedirectURIs:  sqlxx.StringSlicePipeDelimiter{},
	}
}

func TestMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}
	pgURI := dockertest.RunTestPostgreSQL(t)
	pg, err := pop.NewConnection(&pop.ConnectionDetails{
		URL: pgURI,
	})
	require.NoError(t, err)

	require.NoError(t, resilience.Retry(logrus.New(), time.Second*5, time.Minute*5, func() error {
		if err := pg.Open(); err != nil {
			// an Open error probably means we have a problem with the connections config
			log.Printf("could not open pop connection: %+v", err)
			return err
		}
		return pg.RawQuery("select version()").Exec()
	}))

	x.CleanSQLPop(t, pg)

	tm := NewTestMigrator(t, pg, "../migrations", "./testdata")
	require.NoError(t, tm.Up())

	for i := 1; i <= 2; i++ {
		t.Run(fmt.Sprintf("case=migration %d", i), func(t *testing.T) {
			expected := generateExpectedClient(i)
			actual := &client.Client{}
			require.NoError(t, pg.Find(actual, expected.ClientID))
			assertEqualClients(t, expected, actual)
		})
	}

	x.CleanSQLPop(t, pg)
	require.NoError(t, pg.Close())
}
