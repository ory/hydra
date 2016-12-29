package policy

import "github.com/ory-am/ladon"

// Manager is responsible for managing and persisting policies.
type Manager interface {
	ladon.Manager

	// Update a policy.
	Update(policy ladon.Policy) error
}
