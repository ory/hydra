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

	if request.Domain != "" && request.Subject != "" {
		return errorsx.WithStack(ErrMissingRequiredParameter.WithHint("Fields 'subject' and 'domain' are mutually exclusive, both cannot be set at the same time."))
	}

	if request.Subject == "" && request.Domain == "" {
		return errorsx.WithStack(ErrMissingRequiredParameter.WithHint("Both 'subject' and 'domain' fields are empty, one of them is required."))
	}

	if request.ExpiresAt.IsZero() {
		return errorsx.WithStack(ErrMissingRequiredParameter.WithHint("Field 'expires_at' is required."))
	}

	if request.PublicKeyJWK.KeyID == "" {
		return errorsx.WithStack(ErrMissingRequiredParameter.WithHint("Field 'jwk' must contain JWK with kid header."))
	}

	return nil
}
