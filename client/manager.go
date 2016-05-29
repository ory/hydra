package client

import (
	"github.com/ory-am/fosite"
)

type Manager interface {
	Storage

	Authenticate(id string, secret []byte) (*fosite.DefaultClient, error)
}

type Storage interface {
	fosite.Storage

	CreateClient(c *fosite.DefaultClient) error

	DeleteClient(id string) error

	GetClients() (map[string]*fosite.DefaultClient, error)
}
