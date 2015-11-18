package http

import . "github.com/ory-am/hydra/client"

type client struct {
	ep string
}

func New(endpoint string) Client {
	return &client{
		ep: endpoint,
	}
}

func (c *client) IsAllowed(ar *AuthorizeRequest) (bool, error) {
	return true, nil
}

func (c *client) IsAuthenticated(token string) (bool, error) {
	return true, nil
}
