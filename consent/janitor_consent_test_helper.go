// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"time"

	"github.com/ory/hydra/v2/flow"
	"github.com/ory/x/sqlxx"
)

func NewHandledLoginRequest(challenge string, hasError bool, requestedAt time.Time, authenticatedAt sqlxx.NullTime) *flow.HandledLoginRequest {
	var deniedErr *flow.RequestDeniedError
	if hasError {
		deniedErr = &flow.RequestDeniedError{
			Name:        "consent request denied",
			Description: "some description",
			Hint:        "some hint",
			Code:        403,
			Debug:       "some debug",
			Valid:       true,
		}
	}

	return &flow.HandledLoginRequest{
		ID:              challenge,
		Error:           deniedErr,
		WasHandled:      true,
		RequestedAt:     requestedAt,
		AuthenticatedAt: authenticatedAt,
	}
}

func NewHandledConsentRequest(challenge string, hasError bool, requestedAt time.Time, authenticatedAt sqlxx.NullTime) *flow.AcceptOAuth2ConsentRequest {
	var deniedErr *flow.RequestDeniedError
	if hasError {
		deniedErr = &flow.RequestDeniedError{
			Name:        "consent request denied",
			Description: "some description",
			Hint:        "some hint",
			Code:        403,
			Debug:       "some debug",
			Valid:       true,
		}
	}

	return &flow.AcceptOAuth2ConsentRequest{
		ID:              challenge,
		HandledAt:       sqlxx.NullTime(time.Now().Round(time.Second)),
		Error:           deniedErr,
		RequestedAt:     requestedAt,
		AuthenticatedAt: authenticatedAt,
		WasHandled:      true,
	}
}
