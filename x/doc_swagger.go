/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package x

// Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
// typically 201.
//
// swagger:response emptyResponse
type emptyResponse struct{}

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
type errorOAuth2Default struct {
	// in: body
	Body errorOAuth2
}

// Bad Request Error Response
//
// swagger:response errorOAuth2BadRequest
type errorOAuth2BadRequest struct {
	// in: body
	Body errorOAuth2
}

// Not Found Error Response
//
// swagger:response errorOAuth2NotFound
type errorOAuth2NotFound struct {
	// in: body
	Body errorOAuth2
}
