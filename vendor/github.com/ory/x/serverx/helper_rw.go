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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package serverx

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/herodot"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type enhancedError struct {
	*fosite.RFC6749Error
	trace errors.StackTrace
	ID    string `json:"request_id"`
}

func (e *enhancedError) StackTrace() errors.StackTrace {
	return e.trace
}

func ErrorEnhancerRFC6749(r *http.Request, err error) interface{} {
	var trace []errors.Frame

	if e, ok := err.(stackTracer); ok {
		trace = e.StackTrace()
	}

	if e, ok := errors.Cause(err).(*herodot.DefaultError); ok {

		err := &enhancedError{
			RFC6749Error: &fosite.RFC6749Error{
				Name:        e.Error(),
				Description: e.Reason(),
				Code:        e.StatusCode(),
			},
			ID:    r.Header.Get("X-Request-Id"),
			trace: trace,
		}
		return err
	}

	return &enhancedError{
		RFC6749Error: fosite.ErrorToRFC6749Error(err),
		ID:           r.Header.Get("X-Request-Id"),
		trace:        trace,
	}
}
