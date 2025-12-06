// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/ory/hydra/v2/fosite"
)

func TestRequest(t *testing.T) {
	r := &Request{
		RequestedAt:       time.Now().UTC(),
		Client:            &DefaultClient{},
		RequestedScope:    Arguments{"scope"},
		GrantedScope:      Arguments{"scope"},
		RequestedAudience: Arguments{"scope"},
		GrantedAudience:   Arguments{"scope"},
		Form:              url.Values{"foo": []string{"bar"}},
		Session:           new(DefaultSession),
	}

	assert.Equal(t, r.RequestedAt, r.GetRequestedAt())
	assert.Equal(t, r.Client, r.GetClient())
	assert.Equal(t, r.GrantedScope, r.GetGrantedScopes())
	assert.Equal(t, r.RequestedScope, r.GetRequestedScopes())
	assert.Equal(t, r.Form, r.GetRequestForm())
	assert.Equal(t, r.Session, r.GetSession())
}

func TestMergeRequest(t *testing.T) {
	a := &Request{
		ID:                "123",
		RequestedAt:       time.Now().UTC(),
		Client:            &DefaultClient{ID: "123"},
		RequestedScope:    Arguments{"scope-3", "scope-4"},
		RequestedAudience: Arguments{"aud-3", "aud-4"},
		GrantedScope:      []string{"scope-1", "scope-2"},
		GrantedAudience:   []string{"aud-1", "aud-2"},
		Form:              url.Values{"foo": []string{"fasdf"}},
		Session:           new(DefaultSession),
	}
	b := &Request{
		RequestedAt:    time.Now().UTC(),
		Client:         &DefaultClient{},
		RequestedScope: Arguments{},
		GrantedScope:   []string{},
		Form:           url.Values{},
		Session:        new(DefaultSession),
	}

	b.Merge(a)
	assert.EqualValues(t, a.RequestedAt, b.RequestedAt)
	assert.EqualValues(t, a.Client, b.Client)
	assert.EqualValues(t, a.RequestedScope, b.RequestedScope)
	assert.EqualValues(t, a.RequestedAudience, b.RequestedAudience)
	assert.EqualValues(t, a.GrantedScope, b.GrantedScope)
	assert.EqualValues(t, a.GrantedAudience, b.GrantedAudience)
	assert.EqualValues(t, a.Form, b.Form)
	assert.EqualValues(t, a.Session, b.Session)
	assert.EqualValues(t, a.ID, b.ID)
}

func TestSanitizeRequest(t *testing.T) {
	a := &Request{
		RequestedAt:    time.Now().UTC(),
		Client:         &DefaultClient{ID: "123"},
		RequestedScope: Arguments{"asdff"},
		GrantedScope:   []string{"asdf"},
		Form: url.Values{
			"foo":           []string{"fasdf"},
			"bar":           []string{"fasdf", "faaaa"},
			"baz":           []string{"fasdf"},
			"grant_type":    []string{"code"},
			"response_type": []string{"id_token"},
			"client_id":     []string{"1234"},
			"scope":         []string{"read"},
		},
		Session: new(DefaultSession),
	}

	b := a.Sanitize([]string{"bar", "baz"})
	assert.NotEqual(t, a.Form.Encode(), b.GetRequestForm().Encode())

	assert.Empty(t, b.GetRequestForm().Get("foo"))
	assert.Equal(t, "fasdf", b.GetRequestForm().Get("bar"))
	assert.Equal(t, []string{"fasdf", "faaaa"}, b.GetRequestForm()["bar"])
	assert.Equal(t, "fasdf", b.GetRequestForm().Get("baz"))

	assert.Equal(t, "fasdf", a.GetRequestForm().Get("foo"))
	assert.Equal(t, "fasdf", a.GetRequestForm().Get("bar"))
	assert.Equal(t, []string{"fasdf", "faaaa"}, a.GetRequestForm()["bar"])
	assert.Equal(t, "fasdf", a.GetRequestForm().Get("baz"))
	assert.Equal(t, "code", a.GetRequestForm().Get("grant_type"))
	assert.Equal(t, "id_token", a.GetRequestForm().Get("response_type"))
	assert.Equal(t, "1234", a.GetRequestForm().Get("client_id"))
	assert.Equal(t, "read", a.GetRequestForm().Get("scope"))
}

func TestIdentifyRequest(t *testing.T) {
	a := &Request{
		RequestedAt:    time.Now().UTC(),
		Client:         &DefaultClient{},
		RequestedScope: Arguments{},
		GrantedScope:   []string{},
		Form:           url.Values{"foo": []string{"bar"}},
		Session:        new(DefaultSession),
	}

	b := a.Sanitize([]string{})
	b.GetID()
	assert.Equal(t, a.ID, b.GetID())
}
