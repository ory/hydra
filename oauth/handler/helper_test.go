package handler_test

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ory-am/hydra/oauth/provider"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type prov struct{}

type userAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *prov) GetAuthenticationURL(state string) string {
	return fmt.Sprintf("/remote/oauth2/auth?response_type=code&client_id=someclient&state=%s&redirect_uri=%s", state, `%2Foauth2%2Fauth`)
}

func (p *prov) FetchSession(code string) (provider.Session, error) {
	if code != "code" {
		return nil, errors.New("Code not 'code'")
	}
	return &provider.DefaultSession{
		RemoteSubject: "remote-id",
		Extra:         map[string]interface{}{},
	}, nil
}

func (p *prov) GetID() string {
	return "MockProvider"
}

func authHandlerMock(t *testing.T, ts *httptest.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Logf("/remote/oauth2/auth got: %s", r.URL)

		redirect, _ := url.QueryUnescape(r.URL.Query().Get("redirect_uri"))
		parsed, _ := url.Parse(redirect)

		q := url.Values{}
		q.Set("state", r.URL.Query().Get("state"))
		q.Set("code", "code")
		parsed.RawQuery = q.Encode()

		t.Logf("Redirecting to: %s", ts.URL+parsed.String())
		http.Redirect(w, r, ts.URL+parsed.String(), http.StatusFound)
	}
}
