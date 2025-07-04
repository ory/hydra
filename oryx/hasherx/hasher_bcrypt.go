package hasherx

import (
	"context"

	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"golang.org/x/crypto/bcrypt"
)

// ErrBcryptPasswordLengthReached is returned when the password is longer than 72 bytes.
var ErrBcryptPasswordLengthReached = errors.Errorf("passwords are limited to a maximum length of 72 characters")

type (
	// Bcrypt is a hasher that uses the bcrypt algorithm.
	Bcrypt struct {
		c BCryptConfigurator
	}
	// BCryptConfig is the configuration for the bcrypt hasher.
	BCryptConfig struct {
		Cost uint32 `json:"cost"`
	}
	// BCryptConfigurator is the interface that must be implemented by a configuration provider for the bcrypt hasher.
	BCryptConfigurator interface {
		HasherBcryptConfig(ctx context.Context) *BCryptConfig
	}
)

func NewHasherBcrypt(c BCryptConfigurator) *Bcrypt {
	return &Bcrypt{c: c}
}

// Generate generates a hash for the given password.
func (h *Bcrypt) Generate(ctx context.Context, password []byte) ([]byte, error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "hash.Bcrypt.Generate")
	defer span.End()

	if err := validateBcryptPasswordLength(password); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	cost := int(h.c.HasherBcryptConfig(ctx).Cost)
	span.SetAttributes(attribute.Int("bcrypt.cost", cost))
	hash, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return hash, nil
}

func validateBcryptPasswordLength(password []byte) error {
	// Bcrypt truncates the password to the first 72 bytes, following the OpenBSD implementation,
	// so if password is longer than 72 bytes, function returns an error
	// See https://en.wikipedia.org/wiki/Bcrypt#User_input
	if len(password) > 72 {
		return ErrBcryptPasswordLengthReached
	}
	return nil
}

// Understands checks if the given hash is in the correct format.
func (h *Bcrypt) Understands(hash []byte) bool {
	return IsBcryptHash(hash)
}
