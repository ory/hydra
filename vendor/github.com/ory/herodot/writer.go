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
package herodot

import (
	"net/http"
)

// Writer is a helper to write arbitrary data to a ResponseWriter
type Writer interface {
	// Write a response object to the ResponseWriter with status code 200.
	Write(w http.ResponseWriter, r *http.Request, e interface{})

	// WriteCode writes a response object to the ResponseWriter and sets a response code.
	WriteCode(w http.ResponseWriter, r *http.Request, code int, e interface{})

	// WriteCreated writes a response object to the ResponseWriter with status code 201 and
	// the Location header set to location.
	WriteCreated(w http.ResponseWriter, r *http.Request, location string, e interface{})

	// WriteError writes an error to ResponseWriter and tries to extract the error's status code by
	// asserting StatusCodeCarrier. If the error does not implement StatusCodeCarrier, the status code
	// is set to 500.
	WriteError(w http.ResponseWriter, r *http.Request, err interface{})

	// WriteErrorCode writes an error to ResponseWriter and forces an error code.
	WriteErrorCode(w http.ResponseWriter, r *http.Request, code int, err interface{})
}
