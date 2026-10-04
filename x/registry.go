// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/gorilla/sessions"
)

type RegistryCookieStore interface {
	CookieStore(ctx context.Context) (sessions.Store, error)
}

type Networker interface {
	NetworkID(ctx context.Context) uuid.UUID
}

type NetworkProvider interface {
	Networker() Networker
}

type Transactor interface {
	Transaction(ctx context.Context, f func(ctx context.Context) error) error
}
