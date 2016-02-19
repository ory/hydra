package policy

const AllowAccess = "allow"
const DenyAccess = "deny"

// Storage is responsible for managing the policy storage backend.
type Storage interface {

	// Create persists a new policy in the storage backend.
	// The policies subjects, permissions and resources can be regular expressions like "create|delete". They will always have a ^ pre- and $ appended
	// to enforce proper matching. So "create|delete" becomes "^create|delete$".
	Create(policy Policy) error

	// Get retrieves a policy from the storage backend.
	Get(id string) (Policy, error)

	// Delete removes a policy from the storage backend.
	Delete(id string) error

	// Finds all policies associated with subject.
	FindPoliciesForSubject(subject string) ([]Policy, error)
}
