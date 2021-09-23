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
)

var (
	ErrNotFound = &fosite.RFC6749Error{
		CodeField:        http.StatusNotFound,
		ErrorField:       http.StatusText(http.StatusNotFound),
		DescriptionField: "Unable to locate the requested resource",
	}
	ErrConflict = &fosite.RFC6749Error{
		CodeField:        http.StatusConflict,
		ErrorField:       http.StatusText(http.StatusConflict),
		DescriptionField: "Unable to process the requested resource because of conflict in the current state",
	}
)

func LogError(r *http.Request, err error, logger *logrusx.Logger) {
	if logger == nil {
		logger = logrusx.New("", "")
	}

	logger.WithRequest(r).
		WithError(err).Errorln("An error occurred")
}
