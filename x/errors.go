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

import (
	"net/http"

	"github.com/ory/fosite"
	"github.com/ory/x/logrusx"

	"go.opentelemetry.io/otel/plugin/httptrace"
)

var (
	ErrNotFound = &fosite.RFC6749Error{
		Code:        http.StatusNotFound,
		Name:        http.StatusText(http.StatusNotFound),
		Description: "Unable to located the requested resource",
	}
	ErrConflict = &fosite.RFC6749Error{
		Code:        http.StatusConflict,
		Name:        http.StatusText(http.StatusConflict),
		Description: "Unable to process the requested resource because of conflict in the current state",
	}
)

func LogError(r *http.Request, err error, logger *logrusx.Logger) {
	if logger == nil {
		logger = logrusx.New("", "")
	}

	_, _, spanCtx := httptrace.Extract(r.Context(), r)

	if spanCtx.HasTraceID() {
		logger = logger.WithField("trace_id", spanCtx.TraceIDString())
	}
	if spanCtx.HasSpanID() {
		logger = logger.WithField("span_id", spanCtx.SpanIDString())
	}

	logger.WithRequest(r).
		WithError(err).Errorln("An error occurred")
}
