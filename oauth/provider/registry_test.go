package provider_test

import "testing"

import (
	. "github.com/ory-am/hydra/oauth/provider"
	"github.com/stretchr/testify/assert"
)

type provider struct{}

func (p *provider) GetAuthenticationURL(state string) string {
	return "auth"
}

func (p *provider) FetchSession(code string) (Session, error) {
	return &DefaultSession{}, nil
}
func (p *provider) GetID() string {
	return "fooBar"
}

func TestRegistry(t *testing.T) {
	m := &provider{}
	r := NewRegistry([]Provider{m})

	p, err := r.Find("fooBar")
	assert.Nil(t, err)
	assert.Equal(t, m, p)

	p, err = r.Find("foobar")
	assert.Nil(t, err)
	assert.Equal(t, m, p)

	_, err = r.Find("bar")
	assert.NotNil(t, err)
}
