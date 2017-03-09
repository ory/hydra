// Package jwt is able to generate and validate json web tokens.
// Follows https://tools.ietf.org/html/draft-ietf-oauth-json-web-token-32
package jwt

import (
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// RS256JWTStrategy is responsible for generating and validating JWT challenges
type RS256JWTStrategy struct {
	PrivateKey *rsa.PrivateKey
}

// Generate generates a new authorize code or returns an error. set secret
func (j *RS256JWTStrategy) Generate(claims jwt.Claims, header Mapper) (string, string, error) {
	if header == nil || claims == nil {
		return "", "", errors.New("Either claims or header is nil.")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header = assign(token.Header, header.ToMap())

	var sig, sstr string
	var err error
	if sstr, err = token.SigningString(); err != nil {
		return "", "", errors.WithStack(err)
	}

	if sig, err = token.Method.Sign(sstr, j.PrivateKey); err != nil {
		return "", "", errors.WithStack(err)
	}

	return fmt.Sprintf("%s.%s", sstr, sig), sig, nil
}

// Validate validates a token and returns its signature or an error if the token is not valid.
func (j *RS256JWTStrategy) Validate(token string) (string, error) {
	if _, err := j.Decode(token); err != nil {
		return "", errors.WithStack(err)
	}

	return j.GetSignature(token)
}

// Decode will decode a JWT token
func (j *RS256JWTStrategy) Decode(token string) (*jwt.Token, error) {
	// Parse the token.
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return &j.PrivateKey.PublicKey, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "Couldn't parse token")
	} else if !parsedToken.Valid {
		return nil, errors.Errorf("Token is invalid")
	}

	return parsedToken, err
}

// GetSignature will return the signature of a token
func (j *RS256JWTStrategy) GetSignature(token string) (string, error) {
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return "", errors.New("Header, body and signature must all be set")
	}
	return split[2], nil
}

// Hash will return a given hash based on the byte input or an error upon fail
func (j *RS256JWTStrategy) Hash(in []byte) ([]byte, error) {
	// SigningMethodRS256
	hash := sha256.New()
	_, err := hash.Write(in)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}
	return hash.Sum([]byte{}), nil
}

// GetSigningMethodLength will return the length of the signing method
func (j *RS256JWTStrategy) GetSigningMethodLength() int {
	return jwt.SigningMethodRS256.Hash.Size()
}

func assign(a, b map[string]interface{}) map[string]interface{} {
	for k, w := range b {
		if _, ok := a[k]; ok {
			continue
		}
		a[k] = w
	}
	return a
}
