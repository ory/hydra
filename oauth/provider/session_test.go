package provider_test

import (
	. "github.com/ory-am/hydra/oauth/provider"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

func TestDefaultSession(t *testing.T) {
	token := &oauth2.Token{}
	s := DefaultSession{RemoteSubject: "subject", Extra: "extra", Token: token}
	assert.Equal(t, "subject", s.GetRemoteSubject())
	assert.Equal(t, "extra", s.GetExtra().(string))
	assert.Equal(t, token, s.GetToken())
}
