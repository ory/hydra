package mock

import (
	"errors"
	. "github.com/ory-am/hydra/client"
	"net/http"
)

type client struct {
	result bool
	err    error
}

func NewAlwaysTrue() Client {
	return &client{
		result: true,
		err:    nil,
	}
}

func NewAlwaysFalse() Client {
	return &client{
		result: false,
		err:    errors.New("auth denied"),
	}
}

func (c *client) IsRequestAllowed(req *http.Request, resource, permission, owner string) (bool, error) {
	return c.result, c.err
}

func (c *client) IsAllowed(ar *Action) (bool, error) {
	return c.result, c.err
}

func (c *client) IsAuthenticated(token string) (bool, error) {
	return c.result, c.err
}
