package oauth2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"github.com/stretchr/testify/require"
)

func TestClientCredentials(t *testing.T) {
	tok, err := oauthClientConfig.Token(oauth2.NoContext)
	require.NoError(t, err)
	assert.NotEmpty(t, tok.AccessToken)
}
