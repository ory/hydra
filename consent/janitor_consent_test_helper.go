// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"time"

	"github.com/ory/x/sqlxx"
)

func NewHandledLoginRequest(challenge string, hasError bool, requestedAt time.Time, authenticatedAt sqlxx.NullTime) *HandledLoginRequest {
	var deniedErr *RequestDeniedError
	if hasError {
		deniedErr = &RequestDeniedError{
			Name:        "consent request denied",
			Description: "some description",
			Hint:        "some hint",
			Code:        403,
			Debug:       "some debug",
			valid:       true,
		}
	}

	return &HandledLoginRequest{
		ID:              challenge,
		Error:           deniedErr,
		WasHandled:      true,
		RequestedAt:     requestedAt,
		AuthenticatedAt: authenticatedAt,
	}
}

func NewHandledConsentRequest(challenge string, hasError bool, requestedAt time.Time, authenticatedAt sqlxx.NullTime) *AcceptOAuth2ConsentRequest {
	var deniedErr *RequestDeniedError
	if hasError {
		deniedErr = &RequestDeniedError{
			Name:        "consent request denied",
			Description: "some description",
			Hint:        "some hint",
			Code:        403,
			Debug:       "some debug",
			valid:       true,
		}
	}

	return &AcceptOAuth2ConsentRequest{
		ID:              challenge,
		HandledAt:       sqlxx.NullTime(time.Now().Round(time.Second)),
		Error:           deniedErr,
		RequestedAt:     requestedAt,
		AuthenticatedAt: authenticatedAt,
		WasHandled:      true,
	}
}
