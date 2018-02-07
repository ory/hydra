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

package main

// The standard error format
// swagger:response genericError
type genericError struct {
	// in: body
	Body struct {
		// Code is the HTTP status code
		//
		// example: 404
		Code int `json:"code,omitempty"`

		// Status is the HTTP status
		//
		// example: "Not found"
		Status string `json:"status,omitempty"`

		// Request is the value from the X-Request-Id HTTP header, if set
		//
		// Example: de8760cf-59c7-4b9a-87da-59b8c748d9db
		Request string `json:"request,omitempty"`

		// Reason is the reason why the request failed
		//
		// example: "The element could not be found in the database"
		Reason string `json:"reason,omitempty"`

		// Details contains detailed information aimed at debugging the error
		Details []map[string]interface{} `json:"details,omitempty"`

		// Message is a human readable error message
		//
		// example: "The request failed because the resource could not be found"
		Message string `json:"message"`
	}
}

// An empty response
// swagger:response emptyResponse
type emptyResponse struct{}
