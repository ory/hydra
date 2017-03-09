package fosite

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizeRequest(t *testing.T) {
	var urlparse = func(rawurl string) *url.URL {
		u, _ := url.Parse(rawurl)
		return u
	}

	for k, c := range []struct {
		ar           *AuthorizeRequest
		isRedirValid bool
	}{
		{
			ar:           NewAuthorizeRequest(),
			isRedirValid: false,
		},
		{
			ar: &AuthorizeRequest{
				RedirectURI: urlparse("https://foobar"),
			},
			isRedirValid: false,
		},
		{
			ar: &AuthorizeRequest{
				RedirectURI: urlparse("https://foobar"),
				Request: Request{
					Client: &DefaultClient{RedirectURIs: []string{""}},
				},
			},
			isRedirValid: false,
		},
		{
			ar: &AuthorizeRequest{
				Request: Request{
					Client: &DefaultClient{RedirectURIs: []string{""}},
				},
				RedirectURI: urlparse(""),
			},
			isRedirValid: false,
		},
		{
			ar: &AuthorizeRequest{
				Request: Request{
					Client: &DefaultClient{RedirectURIs: []string{""}},
				},
				RedirectURI: urlparse(""),
			},
			isRedirValid: false,
		},
		{
			ar: &AuthorizeRequest{
				RedirectURI: urlparse("https://foobar.com#123"),
				Request: Request{
					Client: &DefaultClient{RedirectURIs: []string{"https://foobar.com#123"}},
				},
			},
			isRedirValid: false,
		},
		{
			ar: &AuthorizeRequest{
				Request: Request{
					Client: &DefaultClient{RedirectURIs: []string{"https://foobar.com"}},
				},
				RedirectURI: urlparse("https://foobar.com#123"),
			},
			isRedirValid: false,
		},
		{
			ar: &AuthorizeRequest{
				Request: Request{
					Client:      &DefaultClient{RedirectURIs: []string{"https://foobar.com/cb"}},
					RequestedAt: time.Now(),
					Scopes:      []string{"foo", "bar"},
				},
				RedirectURI:   urlparse("https://foobar.com/cb"),
				ResponseTypes: []string{"foo", "bar"},
				State:         "foobar",
			},
			isRedirValid: true,
		},
	} {
		assert.Equal(t, c.ar.Client, c.ar.GetClient(), "%d", k)
		assert.Equal(t, c.ar.RedirectURI, c.ar.GetRedirectURI(), "%d", k)
		assert.Equal(t, c.ar.RequestedAt, c.ar.GetRequestedAt(), "%d", k)
		assert.Equal(t, c.ar.ResponseTypes, c.ar.GetResponseTypes(), "%d", k)
		assert.Equal(t, c.ar.Scopes, c.ar.GetRequestedScopes(), "%d", k)
		assert.Equal(t, c.ar.State, c.ar.GetState(), "%d", k)
		assert.Equal(t, c.isRedirValid, c.ar.IsRedirectURIValid(), "%d", k)

		c.ar.GrantScope("foo")
		c.ar.SetSession(&DefaultSession{})
		c.ar.SetRequestedScopes([]string{"foo"})
		assert.True(t, c.ar.GetGrantedScopes().Has("foo"))
		assert.True(t, c.ar.GetRequestedScopes().Has("foo"))
		assert.Equal(t, &DefaultSession{}, c.ar.GetSession())
	}
}
