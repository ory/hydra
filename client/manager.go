/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package client

import (
	"context"

	"github.com/ory/fosite"
)

type Manager interface {
	Storage

	Authenticate(ctx context.Context, id string, secret []byte) (*Client, error)
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
