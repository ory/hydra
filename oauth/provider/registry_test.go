package provider_test

import "testing"

import (
	"github.com/RangelReale/osin"
	. "github.com/ory-am/hydra/oauth/provider"
	"github.com/stretchr/testify/assert"
)

type provider struct{}

func (p *provider) GetAuthCodeURL(ar *osin.AuthorizeRequest) string {
	return "auth"
}

func (p *provider) Exchange(code string) (Session, error) {
	return &DefaultSession{}, nil
}
func (p *provider) GetID() string {
	return "mock"
}

func TestRegistry(t *testing.T) {
	m := &provider{}
	r := NewRegistry(map[string]Provider{
		"foo":    m,
		"fooBar": m,
	})

	p, err := r.Find("foo")
	assert.Nil(t, err)
	assert.Equal(t, m, p)

	p, err = r.Find("fooBar")
	assert.Nil(t, err)
	assert.Equal(t, m, p)

	p, err = r.Find("foobar")
	assert.Nil(t, err)
	assert.Equal(t, m, p)

	_, err = r.Find("bar")
	assert.NotNil(t, err)
}
