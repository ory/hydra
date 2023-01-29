// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

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

	t.Logf("%s", buf.String())

	assert.True(t, strings.Contains(buf.String(), "trace"))

	LogError(r, errors.Wrap(new(errStackTracer), ""), l)
}

func TestLogErrorDoesNotPanic(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "https://hydra/some/endpoint", nil)
	if err != nil {
		t.Fatal(err)
	}

	LogError(r, errors.New("asdf"), nil)
}
