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
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ory/x/logrusx"
)

type errStackTracer struct{}

func (s *errStackTracer) StackTrace() errors.StackTrace {
	return errors.StackTrace{}
}

func (s *errStackTracer) Error() string {
	return "foo"
}

func TestLogError(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "https://hydra/some/endpoint", nil)
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer([]byte{})
	l := logrusx.New("", "", logrusx.ForceLevel(logrus.TraceLevel))
	l.Logger.Out = buf
	LogError(r, errors.New("asdf"), l)

	t.Logf("%s", string(buf.Bytes()))

	assert.True(t, strings.Contains(string(buf.Bytes()), "trace"))

	LogError(r, errors.Wrap(new(errStackTracer), ""), l)
}

func TestLogErrorDoesNotPanic(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "https://hydra/some/endpoint", nil)
	if err != nil {
		t.Fatal(err)
	}

	LogError(r, errors.New("asdf"), nil)
}
