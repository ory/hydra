// Package policy offers management capabilities for access control policies.
// To read up on policies, go to:
//
// - https://github.com/ory/ladon
//
// - https://ory-am.gitbooks.io/hydra/content/policy.html
//
// Contains source files:
//
// - handler.go: A HTTP handler capable of managing policies.
//
// - warden_http.go: A Go API using HTTP to validate  managing policies.
//
// - warden_test.go: Functional tests all of the above.
package policy

// swagger:parameters listPolicies
type swaggerListPolicyParameters struct {
	// The offset from where to start looking.
	// in: query
	Offset int `json:"offset"`

	// The maximum amount of policies returned.
	// in: query
	Limit int `json:"limit"`
}

// swagger:parameters getPolicy deletePolicy
type swaggerGetPolicyParameters struct {
	// The id of the policy.
	// in: path
	ID int `json:"id"`
}

// swagger:parameters updatePolicy
type swaggerUpdatePolicyParameters struct {
	// The id of the policy.
	// in: path
	ID int `json:"id"`

	// in: body
	Body swaggerPolicy
}

// swagger:parameters createPolicy
type swaggerCreatePolicyParameters struct {
	// in: body
	Body swaggerPolicy
}

// A policy
// swagger:response listPolicyResponse
type swaggerListPolicyResponse struct {
	// in: body
	Body swaggerPolicy
}

// swagger:model policy
type swaggerPolicy struct {
	// ID of the policy.
	ID string `json:"id" gorethink:"id"`

	// Description of the policy.
	Description string `json:"description" gorethink:"description"`

	// Subjects impacted by the policy.
	Subjects []string `json:"subjects" gorethink:"subjects"`
	// Effect of the policy
	Effect string `json:"effect" gorethink:"effect"`

	// Resources impacted by the policy.
	Resources []string `json:"resources" gorethink:"resources"`

	// Actions impacted by the policy.
	Actions []string `json:"actions" gorethink:"actions"`

	// Conditions under which the policy is active.
	Conditions map[string]struct {
		Type    string      `json:"type"`
		Options interface{} `json:"options"`
	}
}
