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

// OAuth2 Device Flow
//
// # Ory's OAuth 2.0 Device Authorization API
//
// swagger:model oAuth2ApiDeviceAuthorizationResponse
type oAuth2ApiDeviceAuthorizationResponse struct {
	// The device verification code.
	//
	// example: ory_dc_smldfksmdfkl.mslkmlkmlk
	DeviceCode string `json:"device_code"`

	// The end-user verification code.
	//
	// example: AAAAAA
	UserCode string `json:"user_code"`

	// The end-user verification URI on the authorization
	// server.  The URI should be short and easy to remember as end users
	// will be asked to manually type it into their user agent.
	//
	// example: https://auth.ory.sh/tv
	VerificationUri string `json:"verification_uri"`

	// A verification URI that includes the "user_code" (or
	// other information with the same function as the "user_code"),
	// which is designed for non-textual transmission.
	//
	// example: https://auth.ory.sh/tv?user_code=AAAAAA
	VerificationUriComplete string `json:"verification_uri_complete"`

	// The lifetime in seconds of the "device_code" and "user_code".
	//
	// example: 16830
	ExpiresIn int `json:"expires_in"`

	// The minimum amount of time in seconds that the client
	// SHOULD wait between polling requests to the token endpoint.  If no
	// value is provided, clients MUST use 5 as the default.
	//
	// example: 5
	Interval int `json:"interval"`
}
