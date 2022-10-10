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

// OAuth2 API Error
//
// An API error caused by Ory's OAuth 2.0 APIs.
//
// swagger:model oAuth2ApiError
type oAuth2ApiError struct {
	// Name is the error name.
	//
	// example: The requested resource could not be found
	Name string `json:"error"`

	// Description contains further information on the nature of the error.
	//
	// example: Object with ID 12345 does not exist
	Description string `json:"error_description"`

	// Code represents the error status code (404, 403, 401, ...).
	//
	// example: 404
	Code int `json:"status_code"`

	// Debug contains debug information. This is usually not available and has to be enabled.
	//
	// example: The database adapter was unable to find the element
	Debug string `json:"error_debug"`
}

// Empty responses are sent when, for example, resources are deleted. The HTTP status code for empty responses is
// typically 201.
//
// swagger:response emptyResponse
type emptyResponse struct{}

// Error
//
// swagger:model error
type errorModel struct {
	// Error ID
	//
	// Useful when trying to identify various errors in application logic.
	ID string `json:"id,omitempty"`

	// HTTP Status Code
	//
	// example: 404
	StatusCode int `json:"code,omitempty"`

	// HTTP Status Description
	//
	// example: Not Found
	Status string `json:"status,omitempty"`

	// HTTP Request ID
	//
	// The request ID is often exposed internally in order to trace
	// errors across service architectures. This is often a UUID.
	//
	// example: d7ef54b1-ec15-46e6-bccb-524b82c035e6
	RIDField string `json:"request,omitempty"`

	// Error Reason
	//
	// example: User with ID 1234 does not exist.
	ReasonField string `json:"reason,omitempty"`

	// Debug Details
	//
	// This field is often not exposed to protect against leaking
	// sensitive information.
	//
	// example: SQL field "foo" is not a bool.
	DebugField string `json:"debug,omitempty"`

	// Additional Error Details
	//
	// Further error details
	DetailsField map[string]interface{} `json:"details,omitempty"`

	// Error Message
	//
	// The error's message.
	//
	// example: The resource could not be found
	// required: true
	ErrorField string `json:"message"`
}

// swagger:model errorBody
type errorBody struct {
	errorModel
}

// Default Error Response
//
// swagger:response errorDefault
type errorDefault struct {
	// in: body
	Body errorBody
}

// Bad Request Error Response
//
// swagger:response errorBadRequest
type errorBadRequest struct {
	// in: body
	Body errorBody
}

// Not Found Error Response
//
// swagger:response errorNotFound
type errorNotFound struct {
	// in: body
	Body errorBody
}
