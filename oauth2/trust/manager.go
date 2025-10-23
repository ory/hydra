// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package trust

import (
	"context"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"

	keysetpagination "github.com/ory/x/pagination/keysetpagination_v2"
)

type GrantManager interface {
	CreateGrant(ctx context.Context, g Grant, publicKey jose.JSONWebKey) error
	GetConcreteGrant(ctx context.Context, id uuid.UUID) (Grant, error)
	DeleteGrant(ctx context.Context, id uuid.UUID) error
	GetGrants(ctx context.Context, optionalIssuer string, pageOpts ...keysetpagination.Option) ([]Grant, *keysetpagination.Paginator, error)
	FlushInactiveGrants(ctx context.Context, notAfter time.Time, limit int, batchSize int) error
}
