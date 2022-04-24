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

// Generic Error Response
//
// Error responses are sent when an error (e.g. unauthorized, bad request, ...) occurred.
//
// swagger:model jsonError
type jsonError struct {
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
