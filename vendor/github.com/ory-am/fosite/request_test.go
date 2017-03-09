package fosite_test

import (
	"net/url"
	"testing"
	"time"

	. "github.com/ory-am/fosite"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	r := &Request{
		RequestedAt:   time.Now(),
		Client:        &DefaultClient{},
		Scopes:        Arguments{},
		GrantedScopes: []string{},
		Form:          url.Values{"foo": []string{"bar"}},
		Session:       new(DefaultSession),
	}

	assert.Equal(t, r.RequestedAt, r.GetRequestedAt())
	assert.Equal(t, r.Client, r.GetClient())
	assert.Equal(t, r.GrantedScopes, r.GetGrantedScopes())
	assert.Equal(t, r.Scopes, r.GetRequestedScopes())
	assert.Equal(t, r.Form, r.GetRequestForm())
	assert.Equal(t, r.Session, r.GetSession())
}

func TestMergeRequest(t *testing.T) {
	a := &Request{
		RequestedAt:   time.Now(),
		Client:        &DefaultClient{ID: "123"},
		Scopes:        Arguments{"asdff"},
		GrantedScopes: []string{"asdf"},
		Form:          url.Values{"foo": []string{"fasdf"}},
		Session:       new(DefaultSession),
	}
	b := &Request{
		RequestedAt:   time.Now(),
		Client:        &DefaultClient{},
		Scopes:        Arguments{},
		GrantedScopes: []string{},
		Form:          url.Values{},
		Session:       new(DefaultSession),
	}

	b.Merge(a)
	assert.EqualValues(t, a.RequestedAt, b.RequestedAt)
	assert.EqualValues(t, a.Client, b.Client)
	assert.EqualValues(t, a.Scopes, b.Scopes)
	assert.EqualValues(t, a.GrantedScopes, b.GrantedScopes)
	assert.EqualValues(t, a.Form, b.Form)
	assert.EqualValues(t, a.Session, b.Session)
}
