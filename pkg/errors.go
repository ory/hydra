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

package pkg

import (
	"net/http"
	"reflect"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/ory/fosite"
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

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func LogError(err error, logger log.FieldLogger) {
	extra := map[string]interface{}{}
	if logger == nil {
		logger = log.New()
	}

	if e, ok := errors.Cause(err).(*fosite.RFC6749Error); ok {
		if e.Debug != "" {
			extra["debug"] = e.Debug
		}
		if e.Hint != "" {
			extra["hint"] = e.Hint
		}
		if e.Description != "" {
			extra["description"] = e.Description
		}
	}

	logger.WithError(err).WithFields(extra).Errorln("An error occurred")
	if e, ok := err.(stackTracer); ok {
		logger.Debugf("Stack trace: %+v", e.StackTrace())
	} else {
		logger.Debugf("Stack trace could not be recovered from error type %s", reflect.TypeOf(err))
	}
}
