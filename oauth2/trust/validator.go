// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package trust

import (
	"github.com/ory/x/errorsx"
)

type GrantValidator struct {
}

func NewGrantValidator() *GrantValidator {
	return &GrantValidator{}
}

func (v *GrantValidator) Validate(request createGrantRequest) error {
	if request.Issuer == "" {
		return errorsx.WithStack(ErrMissingRequiredParameter.WithHint("Field 'issuer' is required."))
	}

	if request.Subject == "" && !request.AllowAnySubject {
		return errorsx.WithStack(ErrMissingRequiredParameter.WithHint("One of 'subject' or 'allow_any_subject' field must be set."))
	}

	if request.Subject != "" && request.AllowAnySubject {
		return errorsx.WithStack(ErrMissingRequiredParameter.WithHint("Both 'subject' and 'allow_any_subject' fields cannot be set at the same time."))
	}

	if request.ExpiresAt.IsZero() {
		return errorsx.WithStack(ErrMissingRequiredParameter.WithHint("Field 'expires_at' is required."))
	}

	if request.PublicKeyJWK.KeyID == "" {
		return errorsx.WithStack(ErrMissingRequiredParameter.WithHint("Field 'jwk' must contain JWK with kid header."))
	}

	return nil
}
