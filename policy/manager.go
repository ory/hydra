package policy

import (
	"github.com/jmoiron/sqlx"
	"github.com/ory/ladon"
	"github.com/ory/ladon/manager/sql"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
)

// Manager is responsible for managing and persisting policies.
type Manager interface {
	// Create persists the policy.
	Create(policy ladon.Policy) error

	// Get retrieves a policy.
	Get(id string) (ladon.Policy, error)

	// Delete removes a policy.
	Delete(id string) error

	// List policies.
	List(limit, offset int64) (ladon.Policies, error)

	// Update a policy.
	Update(policy ladon.Policy) error
}

func NewSQLManager(db *sqlx.DB) (ladon.Manager, error) {
	m := sql.NewSQLManager(db, nil)
	migrate.SetTable("hydra_policy_migration")
	if err := m.CreateSchemas(); err != nil {
		return nil, errors.WithStack(err)
	}
	return m, nil
}
