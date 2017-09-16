package oauth2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestClientCredentials(t *testing.T) {
	tok, err := oauthClientConfig.Token(oauth2.NoContext)
	require.NoError(t, err)
	assert.NotEmpty(t, tok.AccessToken)
}
