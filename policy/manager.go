package policy

import (
	"github.com/ory/ladon"
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
