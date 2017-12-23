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

package pkg

import (
	"net/http"
	"reflect"

	"github.com/ory/fosite"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotFound = &RichError{
		Status: http.StatusNotFound,
		error:  errors.New("Not found"),
	}
)

type RichError struct {
	Status int
	error
}

func (e *RichError) StatusCode() int {
	return e.Status
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func LogError(err error, logger log.FieldLogger) {
	extra := map[string]interface{}{}
	if logger == nil {
		logger = log.New()
	}

	if e, ok := err.(*fosite.RFC6749Error); ok {
		if e.Debug != "" {
			extra["debug"] = e.Debug
		}
		if e.Hint != "" {
			extra["hint"] = e.Hint
		}
	} else if e, ok := errors.Cause(err).(*fosite.RFC6749Error); ok {
		if e.Debug != "" {
			extra["debug"] = e.Debug
		}
		if e.Hint != "" {
			extra["hint"] = e.Hint
		}
	}

	if e, ok := errors.Cause(err).(stackTracer); ok {
		logger.WithError(err).WithFields(extra).Errorln("An error occurred")
		logger.Debugf("Stack trace: %+v", e.StackTrace())
	} else if e, ok := err.(stackTracer); ok {
		logger.WithError(err).WithFields(extra).Errorln("An error occurred")
		logger.Debugf("Stack trace: %+v", e.StackTrace())
	} else {
		logger.WithError(err).WithFields(extra).Errorln("An error occurred")
		logger.Debugf("Stack trace could not be recovered from error type %s", reflect.TypeOf(err))
	}
}
