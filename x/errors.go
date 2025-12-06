// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"net/http"

	"github.com/ory/hydra/v2/fosite"
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

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
