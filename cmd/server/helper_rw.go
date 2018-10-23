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

package server

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/x/serverx"
)

func newJSONWriter(l logrus.FieldLogger) *herodot.JSONWriter {
	w := herodot.NewJSONWriter(l)
	w.ErrorEnhancer = serverx.ErrorEnhancerRFC6749
	return w
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type enhancedError struct {
	*fosite.RFC6749Error
	trace errors.StackTrace
}

func (e *enhancedError) StackTrace() errors.StackTrace {
	return e.trace
}

func writerErrorEnhancer(r *http.Request, err error) interface{} {
	if e, ok := errors.Cause(err).(*herodot.DefaultError); ok {
		var trace []errors.Frame

		if e, ok := err.(stackTracer); ok {
			trace = e.StackTrace()
		}

		err := &enhancedError{
			RFC6749Error: &fosite.RFC6749Error{
				Name:        e.Error(),
				Description: e.Reason(),
				Code:        e.StatusCode(),
			},
			trace: trace,
		}
		return err
	}
	return fosite.ErrorToRFC6749Error(err)
}
