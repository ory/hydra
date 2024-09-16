// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

// Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
// typically 204.
//
// swagger:response emptyResponse
type _ struct{}

// Error
//
// swagger:model errorOAuth2
type errorOAuth2 struct {
	// Error
	Name string `json:"error"`

	// Error Description
	Description string `json:"error_description"`

	// Error Hint
	//
	// Helps the user identify the error cause.
	//
	// Example: The redirect URL is not allowed.
	Hint string `json:"error_hint"`

	// HTTP Status Code
	//
	// Example: 401
	Code int `json:"status_code"`

	// Error Debug Information
	//
	// Only available in dev mode.
	Debug string `json:"error_debug,omitempty"`
}

// Default Error Response
//
// swagger:response errorOAuth2Default
type _ struct {
	// in: body
	Body errorOAuth2
}

// Bad Request Error Response
//
// swagger:response errorOAuth2BadRequest
type _ struct {
	// in: body
	Body errorOAuth2
}

// Not Found Error Response
//
// swagger:response errorOAuth2NotFound
type _ struct {
	// in: body
	Body errorOAuth2
}
