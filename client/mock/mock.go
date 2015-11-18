package http

import . "github.com/ory-am/hydra/client"

type client struct {
	result bool
}

func NewAlwaysTrue() Client {
	return &client{
		result: true,
	}
}

func NewAlwaysFalse() Client {
	return &client{
		result: false,
	}
}

func (c *client) IsAllowed(ar *AuthorizeRequest) (bool, error) {
	return c.result, nil
}

func (c *client) IsAuthenticated(token string) (bool, error) {
	return c.result, nil
}
