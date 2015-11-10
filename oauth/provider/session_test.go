package provider_test

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	. "github.com/ory-am/hydra/oauth/provider"
	"testing"
)

func TestDefaultSession(t *testing.T) {
	token := &oauth2.Token{}
	s := DefaultSession{RemoteSubject: "subject", Extra: "extra", Token: token}
	assert.Equal(t, "subject", s.GetRemoteSubject())
	assert.Equal(t, "extra", s.GetExtra().(string))
	assert.Equal(t, token, s.GetToken())
}
