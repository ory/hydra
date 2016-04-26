package oauth2

import (
	"github.com/ory-am/fosite"
)

type ClientManager interface {
	ClientStorage

	Authenticate(id string, secret []byte) (*OAuth2Client, error)
}

type ClientStorage interface {
	fosite.Storage

	CreateClient(c *OAuth2Client) error

	DeleteClient(id string) error
}
