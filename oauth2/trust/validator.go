// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package trust

import "github.com/pkg/errors"

func validateGrant(request createGrantRequest) error {
	if request.Issuer == "" {
		return errors.WithStack(ErrMissingRequiredParameter.WithHint("Field 'issuer' is required."))
	}

	if request.Subject == "" && !request.AllowAnySubject {
		return errors.WithStack(ErrMissingRequiredParameter.WithHint("One of 'subject' or 'allow_any_subject' field must be set."))
	}

	if request.Subject != "" && request.AllowAnySubject {
		return errors.WithStack(ErrMissingRequiredParameter.WithHint("Both 'subject' and 'allow_any_subject' fields cannot be set at the same time."))
	}

	if request.ExpiresAt.IsZero() {
		return errors.WithStack(ErrMissingRequiredParameter.WithHint("Field 'expires_at' is required."))
	}

	if request.PublicKeyJWK.KeyID == "" {
		return errors.WithStack(ErrMissingRequiredParameter.WithHint("Field 'jwk' must contain JWK with kid header."))
	}

	return nil
}
