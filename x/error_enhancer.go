// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/herodot"
)

type enhancedError struct {
	*fosite.RFC6749Error
	RequestID string `json:"request_id"`
}

func ErrorEnhancer(r *http.Request, err error) interface{} {
	if e := new(herodot.DefaultError); errors.As(err, &e) {
		return &enhancedError{
			RFC6749Error: (&fosite.RFC6749Error{
				ErrorField:       e.Error(),
				DescriptionField: e.Reason(),
				CodeField:        e.StatusCode(),
			}).WithTrace(err),
			RequestID: r.Header.Get("X-Request-Id"),
		}
	}

	return &enhancedError{
		RFC6749Error: fosite.ErrorToRFC6749Error(err),
		RequestID:    r.Header.Get("X-Request-Id"),
	}
}
