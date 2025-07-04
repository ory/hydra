package hasherx

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"time"

	"github.com/ory/x/otelx"

	"github.com/inhies/go-bytesize"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash               = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion       = errors.New("incompatible version of argon2")
	ErrMismatchedHashAndPassword = errors.New("passwords do not match")
)

type (
	// Argon2Config is the configuration for a Argon2 hasher.
	Argon2Config struct {
		// Memory is the amount of memory to use.
		Memory bytesize.ByteSize `json:"memory"`

		// Iterations is the number of iterations to use.
		Iterations uint32 `json:"iterations"`

		// Parallelism is the number of threads to use.
		Parallelism uint8 `json:"parallelism"`

		// SaltLength is the length of the salt to use.
		SaltLength uint32 `json:"salt_length"`

		// KeyLength is the length of the key to use.
		KeyLength uint32 `json:"key_length"`

		// ExpectedDuration is the expected duration of the hash.
		ExpectedDuration time.Duration `json:"expected_duration"`

		// ExpectedDeviation is the expected deviation of the hash.
		ExpectedDeviation time.Duration `json:"expected_deviation"`

		// DedicatedMemory is the amount of dedicated memory to use.
		DedicatedMemory bytesize.ByteSize `json:"dedicated_memory"`
	}
	// Argon2 is a hasher that uses the Argon2 algorithm.
	Argon2 struct {
		c Argon2Configurator
	}
	// Argon2Configurator is a function that returns the Argon2 configuration.
	Argon2Configurator interface {
		HasherArgon2Config(ctx context.Context) *Argon2Config
	}
)

func NewHasherArgon2(c Argon2Configurator) *Argon2 {
	return &Argon2{c: c}
}

func toKB(mem bytesize.ByteSize) (uint32, error) {
	kb := uint64(mem / bytesize.KB)
	if kb > math.MaxUint32 {
		return 0, errors.Errorf("memory %v is too large", mem)
	}
	return uint32(kb), nil
}

// Generate generates a hash for the given password.
func (h *Argon2) Generate(ctx context.Context, password []byte) (_ []byte, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "hash.Argon2.Generate")
	defer otelx.End(span, &err)
	p := h.c.HasherArgon2Config(ctx)
	span.SetAttributes(attribute.String("argon2.config", fmt.Sprintf("#%v", p)))

	salt := make([]byte, p.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	mem, err := toKB(p.Memory)
	if err != nil {
		return nil, err
	}
	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey(password, salt, p.Iterations, mem, p.Parallelism, p.KeyLength)

	var b bytes.Buffer
	if _, err := fmt.Fprintf(
		&b,
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, mem, p.Iterations, p.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.WithStack(err)
	}

	return b.Bytes(), nil
}

// Understands checks if the given hash is in the correct format.
func (h *Argon2) Understands(hash []byte) bool {
	return IsArgon2idHash(hash)
}
