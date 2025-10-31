// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"

	"github.com/ory/hydra/v2/fosite"
	keysetpagination "github.com/ory/x/pagination/keysetpagination_v2"
)

// swagger:ignore
type Filter struct {
	PageOpts []keysetpagination.Option
	Name     string
	Owner    string
	IDs      []string
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

	GetClients(ctx context.Context, filters Filter) ([]Client, *keysetpagination.Paginator, error)

	GetConcreteClient(ctx context.Context, id string) (*Client, error)
}

type ManagerProvider interface {
	ClientManager() Manager
}
