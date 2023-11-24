// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"

	"github.com/ory/fosite"
)

// swagger:ignore
type Filter struct {
	// The maximum amount of clients to returned, upper bound is 500 clients.
	// in: query
	Limit int `json:"limit"`

	// The offset from where to start looking.
	// in: query
	Offset int `json:"offset"`

	// The name of the clients to filter by.
	// in: query
	Name string `json:"client_name"`

	// The owner of the clients to filter by.
	// in: query
	Owner string `json:"owner"`
}

type Manager interface {
	Storage

	AuthenticateClient(ctx context.Context, id string, secret []byte) (*Client, error)
}

type Storage interface {
	GetClient(ctx context.Context, id string) (fosite.Client, error)

	CreateClient(ctx context.Context, c *Client) error

	UpdateClient(ctx context.Context, c *Client) error

	DeleteClient(ctx context.Context, id string) error

	GetClients(ctx context.Context, filters Filter) ([]Client, error)

	CountClients(ctx context.Context) (int, error)

	GetConcreteClient(ctx context.Context, id string) (*Client, error)
}

type ManagerProvider interface {
	ClientManager() Manager
}
