package http_test

import (
	"encoding/json"
	"fmt"
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
		if req.Header.Get("Authorization") != "Bearer fetch-token-ok" {
			http.Error(rw, "", http.StatusUnauthorized)
			return
		}
		pkg.WriteJSON(rw, struct {
			Allowed bool `json:"allowed"`
		}{Allowed: true})
	}).Methods("POST")
	router.HandleFunc("/oauth2/token", tokenHandler).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()

	c := New(ts.URL, "irrelevant", "irrelephant")
	c.SetClientToken(&oauth2.Token{TokenType: "bearer", AccessToken: "foobar"})
	allowed, err := c.IsRequestAllowed(&http.Request{Header: http.Header{"Authorization": []string{"Bearer token"}}}, "", "", "")
	assert.Nil(t, err)
	assert.True(t, allowed)
}

func TestIsAllowedRetriesOnlyOnceWhenTokenIsInvalid(t *testing.T) {
	var count int
	router := mux.NewRouter()
	router.HandleFunc("/oauth2/token", tokenHandler).Methods("POST")
	router.HandleFunc("/guard/allowed", func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, fmt.Sprintf("token invalid try ", count), http.StatusUnauthorized)
		count++
	}).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()
	c := New(ts.URL, "irrelevant", "irrelephant")
	c.SetClientToken(&oauth2.Token{TokenType: "bearer", AccessToken: "foobar"})

	allowed, err := c.IsAllowed(&AuthorizeRequest{
		Permission: "foo",
		Token:      "bar",
		Resource:   "res",
		Context: &operator.Context{
			Owner: "foo",
		},
	})
	assert.NotNil(t, err)
	assert.False(t, allowed)
	assert.Equal(t, 2, count)
}

func TestIsAllowedRetriesWhenTokenIsExpired(t *testing.T) {
	var try int
	router := mux.NewRouter()
	router.HandleFunc("/guard/allowed", func(rw http.ResponseWriter, req *http.Request) {
		try++
		if try == 1 {
			t.Logf("token invalid try ", try)
			http.Error(rw, fmt.Sprintf("token invalid try ", try), http.StatusUnauthorized)
			return
		}

		pkg.WriteJSON(rw, struct {
			Allowed bool `json:"allowed"`
		}{Allowed: true})
	}).Methods("POST")
	router.HandleFunc("/oauth2/token", tokenHandler).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()

	c := New(ts.URL, "irrelevant", "irrelephant")
	c.SetClientToken(&oauth2.Token{TokenType: "bearer", AccessToken: "foobar"})

	allowed, err := c.IsAllowed(&AuthorizeRequest{
		Permission: "foo",
		Token:      "bar",
		Resource:   "res",
		Context: &operator.Context{
			Owner: "foo",
		},
	})
	assert.Nil(t, err, "%s", err)
	assert.True(t, allowed)
	assert.Equal(t, 2, try)
}

func TestIsAuthenticatedRetriesWhenTokenIsExpired(t *testing.T) {
	var try int
	router := mux.NewRouter()
	router.HandleFunc("/oauth2/introspect", func(rw http.ResponseWriter, req *http.Request) {
		try++
		if try == 1 {
			http.Error(rw, fmt.Sprintf("token invalid try ", try), http.StatusUnauthorized)
			return
		}

		pkg.WriteJSON(rw, struct {
			Active bool `json:"active"`
		}{Active: true})
	}).Methods("POST")
	router.HandleFunc("/oauth2/token", tokenHandler).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()

	c := New(ts.URL, "irrelevant", "irrelephant")
	c.SetClientToken(&oauth2.Token{TokenType: "bearer", AccessToken: "client"})
	active, err := c.IsAuthenticated("federated.token")
	assert.Nil(t, err, "%s", err)
	assert.True(t, active)
	assert.Equal(t, 2, try)
}

func TestIsAuthenticatedRetriesOnlyOnceWhenTokenIsExpired(t *testing.T) {
	var count int
	router := mux.NewRouter()
	router.HandleFunc("/oauth2/introspect", func(rw http.ResponseWriter, req *http.Request) {
		count++
		http.Error(rw, fmt.Sprintf("token invalid try ", count), http.StatusUnauthorized)
	}).Methods("POST")
	router.HandleFunc("/oauth2/token", tokenHandler).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()

	c := New(ts.URL, "irrelevant", "irrelephant")
	c.SetClientToken(&oauth2.Token{TokenType: "bearer", AccessToken: "client"})
	active, err := c.IsAuthenticated("federated.token")
	assert.NotNil(t, err)
	assert.False(t, active)
	assert.Equal(t, 2, count)
}

func TestIsAllowed(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/oauth2/token", tokenHandler).Methods("POST")
	router.HandleFunc("/guard/allowed", func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "Bearer fetch-token-ok" {
			http.Error(rw, "", http.StatusUnauthorized)
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

	c := New(ts.URL, "irrelevant", "irrelephant")
	allowed, err := c.IsAllowed(&AuthorizeRequest{Permission: "foo", Token: "bar", Resource: "res", Context: &operator.Context{Owner: "foo"}})
	assert.Nil(t, err)
	assert.True(t, allowed)
}

func TestIsAuthenticated(t *testing.T) {
	router := mux.NewRouter()
	called := false
	router.HandleFunc("/oauth2/token", tokenHandler).Methods("POST")
	router.HandleFunc("/oauth2/introspect", func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "Bearer fetch-token-ok" {
			http.Error(rw, "", http.StatusUnauthorized)
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

	c := New(ts.URL, "irrelevant", "irrelephant")
	active, err := c.IsAuthenticated("federated.token")
	assert.Nil(t, err)
	assert.True(t, active)
	assert.True(t, called)
}
