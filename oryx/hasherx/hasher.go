package hasherx

import (
	"context"
)

// Hasher provides methods for generating and comparing password hashes.
type Hasher interface {
	// Generate returns a hash derived from the password or an error if the hash method failed.
	Generate(ctx context.Context, password []byte) ([]byte, error)

	// Understands returns whether the given hash can be understood by this hasher.
	Understands(hash []byte) bool
}

type HashProvider interface {
	Hasher() Hasher
}

const tracingComponent = "github.com/ory/kratos/hash"
