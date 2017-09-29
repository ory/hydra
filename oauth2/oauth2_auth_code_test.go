package oauth2_test

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestAuthCode(t *testing.T) {
	var code string
	var validConsent bool

	router.GET("/consent", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		cr, response, err := consentClient.GetOAuth2ConsentRequest(r.URL.Query().Get("consent"))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		assert.EqualValues(t, []string{"hydra.*", "offline"}, cr.RequestedScopes)
		assert.Equal(t, r.URL.Query().Get("consent"), cr.Id)
		assert.True(t, strings.Contains(cr.RedirectUrl, "oauth2/auth?client_id=app-client"))

		response, err = consentClient.AcceptOAuth2ConsentRequest(r.URL.Query().Get("consent"), hydra.ConsentRequestAcceptance{
			Subject:     "foo",
			GrantScopes: []string{"hydra.*", "offline"},
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, response.StatusCode)

		http.Redirect(w, r, cr.RedirectUrl, http.StatusFound)
		validConsent = true
	})

	router.GET("/callback", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		code = r.URL.Query().Get("code")
		w.Write([]byte(r.URL.Query().Get("code")))
	})

	cookieJar, _ := cookiejar.New(nil)
	req, err := http.NewRequest("GET", oauthConfig.AuthCodeURL("some-foo-state"), nil)
	require.NoError(t, err)

	resp, err := (&http.Client{Jar: cookieJar}).Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	require.True(t, validConsent)
	require.NotEmpty(t, code)

	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	require.NoError(t, err, code)

	time.Sleep(time.Second * 5)

	res, err := testRefresh(t, token)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func testRefresh(t *testing.T, token *oauth2.Token) (*http.Response, error) {
	req, err := http.NewRequest("POST", oauthClientConfig.TokenURL, strings.NewReader(url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{token.RefreshToken},
	}.Encode()))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(oauthClientConfig.ClientID, oauthClientConfig.ClientSecret)

	return http.DefaultClient.Do(req)
}
