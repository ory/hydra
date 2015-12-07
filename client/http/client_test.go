package http_test

import (
	"encoding/json"
	"github.com/RangelReale/osin"
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

func TestIsRequestAllowed(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/guard/allowed", func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "Bearer foobar", req.Header.Get("Authorization"))
		pkg.WriteJSON(rw, struct {
			Allowed bool `json:"allowed"`
		}{Allowed: true})
	}).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()

	c := New(ts.URL, &oauth2.Token{TokenType: "bearer", AccessToken: "foobar"})
	allowed, err := c.IsRequestAllowed(&http.Request{Header: http.Header{"Authorization": []string{"Bearer token"}}}, "", "", "")
	assert.Nil(t, err)
	assert.True(t, allowed)
}

func TestIsAllowed(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/guard/allowed", func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "Bearer foobar", req.Header.Get("Authorization"))
		var p handler.GrantedPayload
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&p); err != nil {
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

	c := New(ts.URL, &oauth2.Token{TokenType: "bearer", AccessToken: "foobar"})
	allowed, err := c.IsAllowed(&AuthorizeRequest{Permission: "foo", Token: "bar", Resource: "res", Context: &operator.Context{Owner: "foo"}})
	assert.Nil(t, err)
	assert.True(t, allowed)
}

func TestIsAuthenticated(t *testing.T) {
	router := mux.NewRouter()
	called := false
	router.HandleFunc("/oauth2/introspect", func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		bearer := osin.CheckBearerAuth(req)
		assert.NotNil(t, bearer)
		assert.NotEmpty(t, bearer.Code)
		assert.NotEmpty(t, req.Form.Get("token"))
		pkg.WriteJSON(rw, struct {
			Active bool `json:"active"`
		}{Active: true})
		called = true
	}).Methods("POST")
	ts := httptest.NewServer(router)
	defer ts.Close()

	c := New(ts.URL, &oauth2.Token{TokenType: "bearer", AccessToken: "client"})
	active, err := c.IsAuthenticated("federated.token")
	assert.Nil(t, err)
	assert.True(t, active)
	assert.True(t, called)
}
