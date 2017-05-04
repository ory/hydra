// Package warden decides if access requests should be allowed or denied. In a scientific taxonomy, the warden
// is classified as a Policy Decision Point. THe warden's primary goal is to implement `github.com/ory-am/hydra/firewall.Firewall`.
// To read up on the warden, go to:
//
// - https://ory-am.gitbooks.io/hydra/content/policy.html
//
// - http://docs.hydra13.apiary.io/#reference/warden:-access-control-for-resource-providers
//
// Contains source files:
//
// - handler.go: A HTTP handler capable of validating access tokens.
//
// - warden_http.go: A Go API using HTTP to validate access tokens.
//
// - warden_local.go: A Go API using storage managers to validate access tokens.
//
// - warden_test.go: Functional tests all of the above.
package warden

import "github.com/ory/hydra/firewall"

// The allowed response
// swagger:response wardenAllowedResponse
type swaggerWardenAllowedResponse struct {
	// in: body
	Body struct {
		// Allowed is true if the request is allowed or false otherwise
		Allowed bool `json:"allowed"`
	}
}

// swagger:parameters wardenAllowed
type swaggerWardenAllowedParameters struct {
	// in: body
	Body firewall.AccessRequest
}

// swagger:parameters wardenTokenAllowed
type swaggerWardenTokenAllowedParameters struct {
	// in: body
	Body wardenAccessRequest
}

// The token allowed response
// swagger:response wardenTokenAllowedResponse
type swaggerWardenTokenAllowedResponse struct {
	// in: body
	Body struct {
		*firewall.Context

		// Allowed is true if the request is allowed or false otherwise
		Allowed bool `json:"allowed"`
	}
}
