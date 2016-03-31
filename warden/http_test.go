package warden_test

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/gorilla/mux"
	"github.com/ory-am/common/pkg"
	. "github.com/ory-am/hydra/client"
	. "github.com/ory-am/hydra/client/http"
	"github.com/ory-am/hydra/policy/handler"
	"github.com/ory-am/ladon/guard/operator"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func tokenHandler(rw http.ResponseWriter, req *http.Request) {
	pkg.WriteJSON(rw, struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}{
		AccessToken: "fetch-token-ok",
		TokenType:   "Bearer",
		ExpiresIn:   9600,
	})
}

func TestIsRequestAllowed(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/guard/allowed", func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "Basic YXBwOmtleQ==" {
			http.Error(rw, req.Header.Get("Authorization"), http.StatusUnauthorized)
			return
		}
		pkg.WriteJSON(rw, struct {
			Allowed bool `json:"allowed"`
		}{Allowed: true})
	}).Methods("POST")
	router.HandleFunc("/oauth2/token", tokenHandler).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()

	c := New(ts.URL, "app", "key")
	c.SetClientToken(&oauth2.Token{TokenType: "bearer", AccessToken: "foobar"})
	allowed, err := c.IsRequestAllowed(&http.Request{Header: http.Header{"Authorization": []string{"Bearer token"}}}, "", "", "")
	assert.Nil(t, err)
	assert.True(t, allowed)
}

func TestIsAllowed(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/oauth2/token", tokenHandler).Methods("POST")
	router.HandleFunc("/guard/allowed", func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "Basic YXBwOmtleQ==" {
			http.Error(rw, req.Header.Get("Authorization"), http.StatusUnauthorized)
			return
		}
		var p handler.GrantedPayload
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&p); err != nil {
			t.Logf("Could not decode body %s", err)
			pkg.HttpError(rw, errors.New(err), http.StatusBadRequest)
			return
		}

		assert.Equal(t, "foo", p.Permission)
		assert.Equal(t, "bar", p.Token)
		assert.Equal(t, "res", p.Resource)
		assert.Equal(t, "foo", p.Context.Owner)
		pkg.WriteJSON(rw, struct {
			Allowed bool `json:"allowed"`
		}{Allowed: true})
	}).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()

	c := New(ts.URL, "app", "key")
	allowed, err := c.IsAllowed(&Action{Permission: "foo", Token: "bar", Resource: "res", Context: &operator.Context{Owner: "foo"}})
	assert.Nil(t, err, "%s", err)
	assert.True(t, allowed)
}

func TestIsAuthenticated(t *testing.T) {
	router := mux.NewRouter()
	called := false
	router.HandleFunc("/oauth2/token", tokenHandler).Methods("POST")
	router.HandleFunc("/oauth2/introspect", func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "Basic YXBwOmtleQ==" {
			http.Error(rw, req.Header.Get("Authorization"), http.StatusUnauthorized)
			return
		}
		req.ParseForm()
		assert.NotEmpty(t, req.Form.Get("token"))
		pkg.WriteJSON(rw, struct {
			Active bool `json:"active"`
		}{Active: true})
		called = true
	}).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()

	c := New(ts.URL, "app", "key")
	active, err := c.IsAuthenticated("some.token")
	assert.Nil(t, err, "%s", err)
	assert.True(t, active)
	assert.True(t, called)
}
