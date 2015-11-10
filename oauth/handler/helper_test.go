package handler_test

import (
	"bytes"
	"github.com/RangelReale/osin"
	"github.com/go-errors/errors"
	"golang.org/x/oauth2"
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

func (p *prov) GetAuthCodeURL(ar *osin.AuthorizeRequest) string {
	redirect, _ := url.Parse("/oauth2/auth")
	q := redirect.Query()
	q.Set(provider.ProviderQueryParam, p.GetID())
	q.Set(provider.RedirectQueryParam, ar.RedirectUri)
	q.Set(provider.ClientQueryParam, ar.Client.GetId())
	q.Set(provider.ScopeQueryParam, ar.Scope)
	q.Set(provider.StateQueryParam, ar.State)
	q.Set(provider.TypeQueryParam, string(ar.Type))
	redirect.RawQuery = q.Encode()

	var buf bytes.Buffer
	buf.WriteString("/remote/oauth2/auth")
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {"someclient"},
		"redirect_uri":  {redirect.String()},
		"scope":         {""},
		"state":         {ar.State},
	}
	buf.WriteByte('?')
	buf.WriteString(v.Encode())
	return buf.String()
}

func (p *prov) Exchange(code string) (provider.Session, error) {
	if code != "code" {
		return nil, errors.New("Code not 'code'")
	}
	return &provider.DefaultSession{
		RemoteSubject: "remote-id",
		Token:         &oauth2.Token{},
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

		q := parsed.Query()
		q2 := url.Values{}
		q2.Set("provider", q.Get(provider.ProviderQueryParam))
		q2.Set("redirect_uri", q.Get(provider.RedirectQueryParam))
		q2.Set("client_id", q.Get(provider.ClientQueryParam))
		q2.Set("scope", q.Get(provider.ScopeQueryParam))
		q2.Set("state", q.Get(provider.StateQueryParam))
		q2.Set("response_type", q.Get(provider.TypeQueryParam))
		q2.Set("access_code", "code")
		parsed.RawQuery = q2.Encode()

		t.Logf("Redirecting to: %s", ts.URL+parsed.String())
		http.Redirect(w, r, ts.URL+parsed.String(), http.StatusFound)
	}
}
