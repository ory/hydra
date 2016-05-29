package oauth2_test

import (
	"testing"

	"github.com/ory-am/hydra/pkg"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestClientCredentials(t *testing.T) {
	tok, err := oauthClientConfig.Token(oauth2.NoContext)
	pkg.RequireError(t, false, err)
	assert.NotEmpty(t, tok.AccessToken)
}
