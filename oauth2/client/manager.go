package client

import (
	"github.com/ory-am/fosite"
 "github.com/ory-am/fosite/client"
)

type ClientManager interface {
	fosite.Storage

	CreateClient(client.Client) (error)

	RemoveClient(id string) (error)
}
