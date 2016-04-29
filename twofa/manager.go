package twofa

import (
	"github.com/go-errors/errors"
	"github.com/pquerna/otp/totp"
)

type Manager interface {
	Generate(subject string) error
	Validate(subject string) error
}

type TOTP struct {
	Issuer string
	Period uint
}

func (m *TOTP) Generate(subject string) error {
	_, err := totp.Generate(totp.GenerateOpts{
		// Name of the issuing Organization/Company.
		Issuer: m.Issuer,
		// Name of the User's Account (eg, email address)
		AccountName: subject,
		// Number of seconds a TOTP hash is valid for. Defaults to 30 seconds.
		Period: m.Period,
	})
	if err != nil {
		return errors.New(err)
	}

	return nil
}
