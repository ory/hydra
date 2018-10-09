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
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type jsonError struct {
	Error *DefaultError `json:"error"`
}

type reporter func(logger logrus.FieldLogger, args ...interface{}) func(w http.ResponseWriter, r *http.Request, code int, err error)

// json outputs JSON.
type JSONWriter struct {
	logger        logrus.FieldLogger
	Reporter      reporter
	ErrorEnhancer func(r *http.Request, err error) interface{}
}

// NewJSONWriter returns a json
func NewJSONWriter(logger logrus.FieldLogger) *JSONWriter {
	writer := &JSONWriter{logger: logger}

	writer.Reporter = defaultReporter
	writer.ErrorEnhancer = defaultJSONErrorEnhancer
	return writer
}

func defaultJSONErrorEnhancer(r *http.Request, err error) interface{} {
	return &jsonError{Error: toDefaultError(err, r.Header.Get("X-Request-ID"))}
}

// Write a response object to the ResponseWriter with status code 200.
func (h *JSONWriter) Write(w http.ResponseWriter, r *http.Request, e interface{}) {
	h.WriteCode(w, r, http.StatusOK, e)
}

// WriteCode writes a response object to the ResponseWriter and sets a response code.
func (h *JSONWriter) WriteCode(w http.ResponseWriter, r *http.Request, code int, e interface{}) {
	js, err := json.Marshal(e)
	if err != nil {
		h.WriteError(w, r, errors.WithStack(err))
		return
	}

	if code == 0 {
		code = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}

// WriteCreated writes a response object to the ResponseWriter with status code 201 and
// the Location header set to location.
func (h *JSONWriter) WriteCreated(w http.ResponseWriter, r *http.Request, location string, e interface{}) {
	w.Header().Set("Location", location)
	h.WriteCode(w, r, http.StatusCreated, e)
}

// WriteError writes an error to ResponseWriter and tries to extract the error's status code by
// asserting StatusCodeCarrier. If the error does not implement StatusCodeCarrier, the status code
// is set to 500.
func (h *JSONWriter) WriteError(w http.ResponseWriter, r *http.Request, err interface{}) {
	if s, ok := errors.Cause(toError(err)).(StatusCodeCarrier); ok {
		h.WriteErrorCode(w, r, s.StatusCode(), err)
		return
	}

	h.WriteErrorCode(w, r, http.StatusInternalServerError, err)
	return
}

// WriteErrorCode writes an error to ResponseWriter and forces an error code.
func (h *JSONWriter) WriteErrorCode(w http.ResponseWriter, r *http.Request, code int, err interface{}) {
	if err == nil {
		err = toError(err)
	}

	if code == 0 {
		code = http.StatusInternalServerError
	}

	// All errors land here, so it's a really good idea to do the logging here as well!
	h.Reporter(h.logger, "An error occurred while handling a request")(w, r, code, toError(err))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Enhancing must happen after logging or context will be lost.
	if h.ErrorEnhancer != nil {
		err = h.ErrorEnhancer(r, toError(err))
	}

	if err := json.NewEncoder(w).Encode(err); err != nil {
		// There was an error, but there's actually not a lot we can do except log that this happened.
		h.Reporter(h.logger, "Could not write jsonError to response writer")(w, r, code, errors.WithStack(err))
	}
}
