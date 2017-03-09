// Package hmac is the default implementation for generating and validating challenges. It uses HMAC-SHA256 to
// generate and validate challenges.
package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ory-am/fosite"
	"github.com/pkg/errors"
)

// HMACStrategy is responsible for generating and validating challenges.
type HMACStrategy struct {
	AuthCodeEntropy int
	GlobalSecret    []byte
}

const (
	// key should be at least 256 bit long, making it
	minimumEntropy = 32

	// the secrets (client and global) should each have at least 16 characters making it harder to guess them
	minimumSecretLength = 32
)

var b64 = base64.URLEncoding.WithPadding(base64.NoPadding)

// Generate generates a token and a matching signature or returns an error.
// This method implements rfc6819 Section 5.1.4.2.2: Use High Entropy for Secrets.
func (c *HMACStrategy) Generate() (string, string, error) {
	if len(c.GlobalSecret) < minimumSecretLength/2 {
		return "", "", errors.New("Secret is not strong enough")
	}

	if c.AuthCodeEntropy < minimumEntropy {
		c.AuthCodeEntropy = minimumEntropy
	}

	// When creating secrets not intended for usage by human users (e.g.,
	// client secrets or token handles), the authorization server should
	// include a reasonable level of entropy in order to mitigate the risk
	// of guessing attacks.  The token value should be >=128 bits long and
	// constructed from a cryptographically strong random or pseudo-random
	// number sequence (see [RFC4086] for best current practice) generated
	// by the authorization server.
	key, err := RandomBytes(c.AuthCodeEntropy)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	if len(key) < c.AuthCodeEntropy {
		return "", "", errors.New("Could not read enough random data for key generation")
	}

	useSecret := append([]byte{}, c.GlobalSecret...)
	mac := hmac.New(sha256.New, useSecret)
	_, err = mac.Write(key)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	signature := mac.Sum([]byte{})
	encodedSignature := b64.EncodeToString(signature)
	encodedToken := fmt.Sprintf("%s.%s", b64.EncodeToString(key), encodedSignature)
	return encodedToken, encodedSignature, nil
}

// Validate validates a token and returns its signature or an error if the token is not valid.
func (c *HMACStrategy) Validate(token string) error {
	split := strings.Split(token, ".")
	if len(split) != 2 {
		return errors.WithStack(fosite.ErrInvalidTokenFormat)
	}

	key := split[0]
	signature := split[1]
	if key == "" || signature == "" {
		return errors.WithStack(fosite.ErrInvalidTokenFormat)
	}

	decodedSignature, err := b64.DecodeString(signature)
	if err != nil {
		return errors.WithStack(err)
	}

	decodedKey, err := b64.DecodeString(key)
	if err != nil {
		return errors.WithStack(err)
	}

	useSecret := append([]byte{}, c.GlobalSecret...)
	mac := hmac.New(sha256.New, useSecret)
	_, err = mac.Write(decodedKey)
	if err != nil {
		return errors.WithStack(err)
	}

	if !hmac.Equal(decodedSignature, mac.Sum([]byte{})) {
		// Hash is invalid
		return errors.WithStack(fosite.ErrTokenSignatureMismatch)
	}

	return nil
}

func (c *HMACStrategy) Signature(token string) string {
	split := strings.Split(token, ".")

	if len(split) != 2 {
		return ""
	}

	return split[1]
}
