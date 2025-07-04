package hasherx

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

var ErrUnknownHashAlgorithm = errors.New("unknown hash algorithm")

// Compare the given password with the given hash.
func Compare(ctx context.Context, password []byte, hash []byte) error {
	switch {
	case IsBcryptHash(hash):
		return CompareBcrypt(ctx, password, hash)
	case IsArgon2idHash(hash):
		return CompareArgon2id(ctx, password, hash)
	case IsArgon2iHash(hash):
		return CompareArgon2i(ctx, password, hash)
	case IsPbkdf2Hash(hash):
		return ComparePbkdf2(ctx, password, hash)
	default:
		return errors.WithStack(ErrUnknownHashAlgorithm)
	}
}

func CompareBcrypt(_ context.Context, password []byte, hash []byte) error {
	if err := validateBcryptPasswordLength(password); err != nil {
		return err
	}

	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return err
	}

	return nil
}

func CompareArgon2id(_ context.Context, password []byte, hash []byte) error {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeArgon2idHash(string(hash))
	if err != nil {
		return err
	}

	mem := uint64(p.Memory)
	if mem > math.MaxUint32 {
		return errors.WithStack(ErrInvalidHash)
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey(password, salt, p.Iterations, uint32(mem), p.Parallelism, p.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return nil
	}
	return errors.WithStack(ErrMismatchedHashAndPassword)
}

func CompareArgon2i(_ context.Context, password []byte, hash []byte) error {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeArgon2idHash(string(hash))
	if err != nil {
		return err
	}

	mem := uint64(p.Memory)
	if mem > math.MaxUint32 {
		return errors.WithStack(ErrInvalidHash)
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.Key(password, salt, p.Iterations, uint32(mem), p.Parallelism, p.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return nil
	}
	return errors.WithStack(ErrMismatchedHashAndPassword)
}

func ComparePbkdf2(_ context.Context, password []byte, hash []byte) error {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodePbkdf2Hash(string(hash))
	if err != nil {
		return err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := pbkdf2.Key(password, salt, int(p.Iterations), int(p.KeyLength), getPseudorandomFunctionForPbkdf2(p.Algorithm))

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return nil
	}
	return errors.WithStack(ErrMismatchedHashAndPassword)
}

var (
	isBcryptHash   = regexp.MustCompile(`^\$2[abzy]?\$`)
	isArgon2idHash = regexp.MustCompile(`^\$argon2id\$`)
	isArgon2iHash  = regexp.MustCompile(`^\$argon2i\$`)
	isPbkdf2Hash   = regexp.MustCompile(`^\$pbkdf2-sha[0-9]{1,3}\$`)
)

func IsBcryptHash(hash []byte) bool {
	return isBcryptHash.Match(hash)
}

func IsArgon2idHash(hash []byte) bool {
	return isArgon2idHash.Match(hash)
}

func IsArgon2iHash(hash []byte) bool {
	return isArgon2iHash.Match(hash)
}

func IsPbkdf2Hash(hash []byte) bool {
	return isPbkdf2Hash.Match(hash)
}

func decodeArgon2idHash(encodedHash string) (p *Argon2Config, salt, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = new(Argon2Config)
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	saltLength := uint(len(salt))
	if saltLength > math.MaxUint32 {
		return nil, nil, nil, ErrInvalidHash
	}
	p.SaltLength = uint32(saltLength)

	hash, err = base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	keyLength := uint(len(hash))
	if keyLength > math.MaxUint32 {
		return nil, nil, nil, ErrInvalidHash
	}
	p.KeyLength = uint32(keyLength)

	return p, salt, hash, nil
}

// decodePbkdf2Hash decodes PBKDF2 encoded password hash.
// format: $pbkdf2-<digest>$i=<iterations>,l=<length>$<salt>$<hash>
func decodePbkdf2Hash(encodedHash string) (p *PBKDF2Config, salt, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 5 {
		return nil, nil, nil, ErrInvalidHash
	}

	p = new(PBKDF2Config)
	digestParts := strings.SplitN(parts[1], "-", 2)
	if len(digestParts) != 2 {
		return nil, nil, nil, ErrInvalidHash
	}
	p.Algorithm = digestParts[1]

	_, err = fmt.Sscanf(parts[2], "i=%d,l=%d", &p.Iterations, &p.KeyLength)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(parts[3])
	if err != nil {
		return nil, nil, nil, err
	}
	saltLength := uint(len(salt))
	if saltLength > math.MaxUint32 {
		return nil, nil, nil, ErrInvalidHash
	}
	p.SaltLength = uint32(saltLength)

	hash, err = base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	keyLength := uint(len(hash))
	if keyLength > math.MaxUint32 {
		return nil, nil, nil, ErrInvalidHash
	}
	p.KeyLength = uint32(keyLength)

	return p, salt, hash, nil
}
