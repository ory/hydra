// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ory/x/logrusx"
)

func TestLogAudit(t *testing.T) {
	for k, tc := range []struct {
		d              string
		message        interface{}
		expectContains []string
	}{
		{
			d:              "This should log \"access allowed\" because no errors are given",
			message:        nil,
			expectContains: []string{"msg=access allowed"},
		},
		{
			d:              "This should log \"access denied\" because an error is given",
			message:        errors.New("asdf"),
			expectContains: []string{"msg=access denied"},
		},
	} {
		t.Run(fmt.Sprintf("case=%d/description=%s", k, tc.d), func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "https://hydra/some/endpoint", nil)
			if err != nil {
				t.Fatal(err)
			}
			buf := bytes.NewBuffer([]byte{})
			l := logrusx.NewAudit("", "", logrusx.ForceLevel(logrus.TraceLevel))
			l.Logger.Out = buf
			LogAudit(r, tc.message, l)

			assert.Contains(t, buf.String(), "audience=audit")
			for _, expectContain := range tc.expectContains {
				assert.Contains(t, buf.String(), expectContain)
			}
		})
	}
}
