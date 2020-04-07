package migratest

import (
	"fmt"
	"github.com/bmizerany/assert"
	"github.com/gobuffalo/pop/v5"
	"github.com/ory/dockertest/v3"
	"github.com/ory/x/resilience"
	"github.com/ory/x/sqlcon"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"log"
	"sort"
	"testing"
	"time"
)

func ConnectToTestCockroachDB(t *testing.T) *pop.Connection {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "cockroachdb/cockroach",
		Tag:        "v19.2.0",
		Cmd:        []string{"start", "--insecure"},
	})
	require.NoError(t, err)

	var c *pop.Connection
	url := fmt.Sprintf("cockroach://root@localhost:%s/defaultdb?sslmode=disable", resource.GetPort("26257/tcp"))
	maxCons, maxIdle, maxLife := sqlcon.ParseConnectionOptions(logrus.New(), url)

	if err := resilience.Retry(logrus.New(), time.Second*5, time.Minute*5, func() error {
		c, err = pop.NewConnection(&pop.ConnectionDetails{
			URL:             url,
			ConnMaxLifetime: maxLife,
			Pool:            maxCons,
			IdlePool:        maxIdle,
		})
		if err != nil {
			return err
		}

		return c.Open()
	}); err != nil {
		if pErr := pool.Purge(resource); pErr != nil {
			log.Fatalf("Could not connect to docker and unable to remove image: %s - %s", err, pErr)
		}
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return c
}

func TestMigrations(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	cr := ConnectToTestCockroachDB(t)

	tm := NewTestMigrator(t, cr, "../migrations", 13)
	require.NoError(t, tm.Up())
	ms := tm.Migrations["up"]
	sort.Sort(ms)
	fmt.Printf("%+v\n", ms)

	t.Run("case=up to migration 13", func(t *testing.T) {
		client := &client13{
			PK:                                13,
			ClientID:                          "client-13",
			Name:                              "Client-13",
			Secret:                            "secret-13",
			RedirectURIs:                      []string{"http://redirect/13-1"},
			GrantTypes:                        []string{"grant-13_1"},
			ResponseTypes:                     []string{"response-13_1"},
			Scope:                             "scope-13",
			Audience:                          []string{"audience-13_1"},
			Owner:                             "owner-13",
			PolicyURI:                         "http://policy/13",
			AllowedCORSOrigins:                []string{"http://cors/13"},
			TermsOfServiceURI:                 "http://tos/13",
			ClientURI:                         "http://client/13",
			LogoURI:                           "http://logo/13",
			Contacts:                          []string{"contact-13_1"},
			SecretExpiresAt:                   0,
			SubjectType:                       "subject-13",
			SectorIdentifierURI:               "sector-13",
			JSONWebKeysURI:                    "http://jwk/13",
			JSONWebKeys:                       nil,
			TokenEndpointAuthMethod:           "token_auth_method-13",
			RequestURIs:                       []string{"http://request/13"},
			RequestObjectSigningAlgorithm:     "request_alg-13",
			UserinfoSignedResponseAlg:         "response_alg-13",
			CreatedAt:                         time.Now(),
			UpdatedAt:                         time.Now(),
			FrontChannelLogoutURI:             "http://logout/13",
			FrontChannelLogoutSessionRequired: true,
			PostLogoutRedirectURIs:            []string{"http://post_logout/13_1"},
			BackChannelLogoutURI:              "http://back_logout/13",
			BackChannelLogoutSessionRequired:  true,
		}

		require.NoError(t, cr.Create(client))

		cmp := &client13{}
		require.NoError(t, cr.Find(cmp, 13))

		assert.Equal(t, client, cmp)
	})
}
