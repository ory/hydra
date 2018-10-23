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

	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// json outputs JSON.
type TextWriter struct {
	logger      logrus.FieldLogger
	Reporter    reporter
	contentType string
}

// NewPlainWriter returns a json
func NewTextWriter(logger logrus.FieldLogger, contentType string) *TextWriter {
	if contentType == "" {
		contentType = "plain"
	}

	writer := &TextWriter{
		logger:      logger,
		contentType: "text/" + contentType,
	}

	writer.Reporter = defaultReporter
	return writer
}

// Write a response object to the ResponseWriter with status code 200.
func (h *TextWriter) Write(w http.ResponseWriter, r *http.Request, e interface{}) {
	h.WriteCode(w, r, http.StatusOK, e)
}

// WriteCode writes a response object to the ResponseWriter and sets a response code.
func (h *TextWriter) WriteCode(w http.ResponseWriter, r *http.Request, code int, e interface{}) {
	if code == 0 {
		code = http.StatusOK
	}

	w.Header().Set("Content-Type", h.contentType)
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", e)
}

// WriteCreated writes a response object to the ResponseWriter with status code 201 and
// the Location header set to location.
func (h *TextWriter) WriteCreated(w http.ResponseWriter, r *http.Request, location string, e interface{}) {
	w.Header().Set("Location", location)
	h.WriteCode(w, r, http.StatusCreated, e)
}

// WriteError writes an error to ResponseWriter and tries to extract the error's status code by
// asserting StatusCodeCarrier. If the error does not implement StatusCodeCarrier, the status code
// is set to 500.
func (h *TextWriter) WriteError(w http.ResponseWriter, r *http.Request, err interface{}) {
	if s, ok := errors.Cause(toError(err)).(StatusCodeCarrier); ok {
		h.WriteErrorCode(w, r, s.StatusCode(), err)
		return
	}

	h.WriteErrorCode(w, r, http.StatusInternalServerError, err)
	return
}

// WriteErrorCode writes an error to ResponseWriter and forces an error code.
func (h *TextWriter) WriteErrorCode(w http.ResponseWriter, r *http.Request, code int, err interface{}) {
	e := toError(err)
	if err == nil {
		err = e
	}

	if code == 0 {
		code = http.StatusInternalServerError
	}

	// All errors land here, so it's a really good idea to do the logging here as well!
	h.Reporter(h.logger, "An error occurred while handling a request")(w, r, code, e)

	w.Header().Set("Content-Type", h.contentType)
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", err)
}
