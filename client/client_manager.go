package client

import (
	"github.com/ory-am/fosite"
)

type ClientManager interface {
	ClientStorage

	Authenticate(id string, secret []byte) (*Client, error)
}

type ClientStorage interface {
	fosite.Storage

	CreateClient(c *Client) error

	DeleteClient(id string) error
}
