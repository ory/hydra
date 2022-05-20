package persistence

import (
	"context"

	"github.com/gobuffalo/pop/v6"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2/trust"
	"github.com/ory/hydra/x"
	"github.com/ory/x/popx"
)

type (
	Persister interface {
		consent.Manager
		client.Manager
		x.FositeStorer
		jwk.Manager
		trust.GrantManager

		MigrationStatus(ctx context.Context) (popx.MigrationStatuses, error)
		MigrateDown(context.Context, int) error
		MigrateUp(context.Context) error
		PrepareMigration(context.Context) error
		Connection(context.Context) *pop.Connection
		Ping() error
	}
	Provider interface {
		Persister() Persister
	}
)
