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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var errNilDependency = errors.New("A dependency was expected to be defined but is nil. Please open an issue with the stack trace.")

func expectDependency(logger logrus.FieldLogger, dependencies ...interface{}) {
	for _, d := range dependencies {
		if d == nil {
			logger.WithError(errors.WithStack(errNilDependency)).Fatalf("A fatal issue occurred.")
		}
	}
}
