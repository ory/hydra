package consent_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ory/x/sqlxx"
	"github.com/pborman/uuid"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/internal/testhelpers"

	"github.com/ory/x/urlx"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/internal"
	hydra "github.com/ory/hydra/internal/httpclient/client"
	"github.com/ory/hydra/internal/httpclient/client/admin"
	"github.com/ory/hydra/internal/httpclient/models"
)

func TestStrategyDeviceConsentNext(t *testing.T) {
	reg := internal.NewMockedRegistry(t)
	reg.Config().MustSet(config.KeyAccessTokenStrategy, "opaque")
	reg.Config().MustSet(config.KeyConsentRequestMaxAge, time.Hour)
	reg.Config().MustSet(config.KeyConsentRequestMaxAge, time.Hour)
	reg.Config().MustSet(config.KeyScopeStrategy, "exact")
	reg.Config().MustSet(config.KeySubjectTypesSupported, []string{"pairwise", "public"})
	reg.Config().MustSet(config.KeySubjectIdentifierAlgorithmSalt, "76d5d2bf-747f-4592-9fbd-d2b895a54b3a")

	_, adminTS := testhelpers.NewOAuth2Server(t, reg)
	adminClient := hydra.NewHTTPClientWithConfig(nil, &hydra.TransportConfig{Schemes: []string{"http"}, Host: urlx.ParseOrPanic(adminTS.URL).Host})

	acceptLoginHandler := func(t *testing.T, subject string, payload *models.AcceptLoginRequest) http.HandlerFunc {
		return checkAndAcceptLoginHandler(t, adminClient.Admin, subject, func(*testing.T, *admin.GetLoginRequestOK, error) *models.AcceptLoginRequest {
			if payload == nil {
				return new(models.AcceptLoginRequest)
			}
			return payload
		})
	}

	acceptConsentHandler := func(t *testing.T, payload *models.AcceptConsentRequest) http.HandlerFunc {
		return checkAndAcceptConsentHandler(t, adminClient.Admin, func(*testing.T, *admin.GetConsentRequestOK, error) *models.AcceptConsentRequest {
			if payload == nil {
				return new(models.AcceptConsentRequest)
			}
			return payload
		})
	}

	createClientWithRedir := func(t *testing.T, redir string) *client.Client {
		c := &client.Client{RedirectURIs: []string{redir}}
		return createClient(t, reg, c)
	}

	createDefaultClient := func(t *testing.T) *client.Client {
		return createClientWithRedir(t, testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler))
	}

	makeRequestAndExpectError := func(t *testing.T, hc *http.Client, c *client.Client, values url.Values, errContains string) {
		_, res := makeOAuth2Request(t, reg, hc, c, values)
		assert.EqualValues(t, http.StatusNotImplemented, res.StatusCode)
		assert.Empty(t, res.Request.URL.Query().Get("code"))
		assert.Contains(t, res.Request.URL.Query().Get("error_description"), errContains, "%v", res.Request.URL.Query())
	}

	t.Run("case=should fail as incorect grant type supplied", func(t *testing.T) {
		c := createDefaultClient(t)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, "aeneas-rekkas", nil),
			acceptConsentHandler(t, nil))

		makeRequestAndExpectError(t, nil, c, url.Values{"device_verifier": {"not_available"}}, "")
	})

	t.Run("case=should fail as unknown device_verifier supplied", func(t *testing.T) {
		c := &client.Client{RedirectURIs: []string{testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler)}}
		c.GrantTypes = sqlxx.StringSlicePipeDelimiter{"urn:ietf:params:oauth:grant-type:device_code"}
		createClient(t, reg, c)
		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, "aeneas-rekkas", nil),
			acceptConsentHandler(t, nil))

		makeRequestAndExpectError(t, nil, c, url.Values{"device_verifier": {"not_available"}}, "")
	})

	t.Run("case=should fail as unknown no csrf available", func(t *testing.T) {
		c := &client.Client{RedirectURIs: []string{testhelpers.NewCallbackURL(t, "callback", testhelpers.HTTPServerNotImplementedHandler)}}
		c.GrantTypes = sqlxx.StringSlicePipeDelimiter{"urn:ietf:params:oauth:grant-type:device_code"}
		createClient(t, reg, c)

		verifier := strings.Replace(uuid.New(), "-", "", -1)
		csrf := strings.Replace(uuid.New(), "-", "", -1)
		challange := strings.Replace(uuid.New(), "-", "", -1)

		reg.ConsentManager().CreateDeviceGrantRequest(context.TODO(), &consent.DeviceGrantRequest{
			ID:       challange,
			Verifier: verifier,
			CSRF:     csrf,
		})

		testhelpers.NewLoginConsentUI(t, reg.Config(),
			acceptLoginHandler(t, "aeneas-rekkas", nil),
			acceptConsentHandler(t, nil))

		makeRequestAndExpectError(t, nil, c, url.Values{"device_verifier": {verifier}}, "")
	})
}
