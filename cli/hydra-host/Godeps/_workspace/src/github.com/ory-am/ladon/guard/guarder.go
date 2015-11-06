package guard

import (
	"github.com/ory-am/ladon/guard/operator"
	"github.com/ory-am/ladon/policy"
)

// Guarder is responsible for deciding if the subject s has permission p on resource r.
type Guarder interface {
	// IsGranted returns true, if subject s has permission p on resource r
	//  policies, _ := store.FindPoliciesForSubject("peter")
	//  granted, error := guard.IsGranted("article/1234", "update", "peter", nil)
	IsGranted(resource, permission, subject string, policies []policy.Policy, ctx *operator.Context) (bool, error)
}
