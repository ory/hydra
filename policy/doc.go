// Package policy offers management capabilities for access control policies.
//
// Access Control Policies (ACP) are a concept similar to Role Based Access Control and Access Control Lists. ACPs
// however are more flexible and capable of handling complex and abstract access control scenarios. A ACP answers "**Who**
// is **able** to do **what** on **something** given a **context**."
//
//
// ACPs have five attributes:
//
// - Subject *(who)*: An arbitrary unique subject name, for example "ken" or "printer-service.mydomain.com".
// - Effect *(able)*: The effect which can be either "allow" or "deny".
// - Action *(what)*: An arbitrary action name, for example "delete", "create" or "scoped:action:something".
// - Resource *(something)*: An arbitrary unique resource name, for example "something", "resources.articles.1234" or some uniform resource name like "urn:isbn:3827370191".
// - Condition *(context)*: An optional condition that evaluates the context (e.g. IP Address, request datetime, resource owner name, department, ...). Different strategies are available to evaluate conditions:
//   - https://github.com/ory/ladon#cidr-condition
//   - https://github.com/ory/ladon#string-equal-condition
//   - https://github.com/ory/ladon#string-match-condition
//	 - https://github.com/ory/ladon#subject-condition
//   - https://github.com/ory/ladon#string-pairs-equal-condition
//
//
// You can find more information on ACPs here:
//
// - https://github.com/ory/ladon#usage for more information on policy usage.
//
// - https://github.com/ory/ladon#concepts
// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	ID string `json:"id"`
}

// swagger:parameters updatePolicy
type swaggerUpdatePolicyParameters struct {
	// The id of the policy.
	// in: path
	ID string `json:"id"`

	// in: body
	Body swaggerPolicy
}

// swagger:parameters createPolicy
type swaggerCreatePolicyParameters struct {
	// in: body
	Body swaggerPolicy
}

// A policy
// swagger:response policyList
type swaggerListPolicyResponse struct {
	// in: body
	// type: array
	Body []swaggerPolicy
}

// swagger:model policy
type swaggerPolicy struct {
	// ID of the policy.
	ID string `json:"id"`

	// Description of the policy.
	Description string `json:"description"`

	// Subjects impacted by the policy.
	Subjects []string `json:"subjects"`
	// Effect of the policy
	Effect string `json:"effect"`

	// Resources impacted by the policy.
	Resources []string `json:"resources"`

	// Actions impacted by the policy.
	Actions []string `json:"actions"`

	// Conditions under which the policy is active.
	Conditions map[string]struct {
		Type    string                 `json:"type"`
		Options map[string]interface{} `json:"options"`
	} `json:"conditions"`
}
