package client

import (
	"github.com/ory-am/fosite"
)

type Manager interface {
	Storage

	Authenticate(id string, secret []byte) (*Client, error)
}

type Storage interface {
	fosite.Storage

	CreateClient(c *Client) error

	UpdateClient(c *Client) error

	DeleteClient(id string) error

	GetClients() (map[string]Client, error)

	GetConcreteClient(id string) (*Client, error)
}
