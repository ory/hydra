// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"

	"github.com/ory/x/errorsx"
)

// Token represets a JWT Token
// This token provide an adaptation to
// transit from [jwt-go](https://github.com/dgrijalva/jwt-go)
// to [go-jose](https://github.com/square/go-jose)
// It provides method signatures compatible with jwt-go but implemented
// using go-json
type Token struct {
	Header map[string]interface{} // The first segment of the token
	Claims MapClaims              // The second segment of the token
	Method jose.SignatureAlgorithm
	valid  bool
}

const (
	SigningMethodNone = jose.SignatureAlgorithm("none")
	// This key should be use to correctly sign and verify alg:none JWT tokens
	UnsafeAllowNoneSignatureType unsafeNoneMagicConstant = "none signing method allowed"

	JWTHeaderType      = jose.HeaderKey("typ")
	JWTHeaderTypeValue = "JWT"
)

type unsafeNoneMagicConstant string

// Valid informs if the token was verified against a given verification key
// and claims are valid
func (t *Token) Valid() bool {
	return t.valid
}

// Claims is a port from https://github.com/dgrijalva/jwt-go/blob/master/claims.go
// including its validation methods, which are not available in go-jose library
//
// > For a type to be a Claims object, it must just have a Valid method that determines
// if the token is invalid for any supported reason
type Claims interface {
	Valid() error
}

// NewWithClaims creates an unverified Token with the given claims and signing method
func NewWithClaims(method jose.SignatureAlgorithm, claims MapClaims) *Token {
	return &Token{
		Claims: claims,
		Method: method,
		Header: map[string]interface{}{},
	}
}

func (t *Token) toJoseHeader() map[jose.HeaderKey]interface{} {
	h := map[jose.HeaderKey]interface{}{
		JWTHeaderType: JWTHeaderTypeValue,
	}
	for k, v := range t.Header {
		h[jose.HeaderKey(k)] = v
	}
	return h
}

// SignedString provides a compatible `jwt-go` Token.SignedString method
//
// > Get the complete, signed token
func (t *Token) SignedString(k interface{}) (rawToken string, err error) {
	if _, ok := k.(unsafeNoneMagicConstant); ok {
		rawToken, err = unsignedToken(t)
		return

	}
	var signer jose.Signer
	key := jose.SigningKey{
		Algorithm: t.Method,
		Key:       k,
	}
	opts := &jose.SignerOptions{ExtraHeaders: t.toJoseHeader()}
	signer, err = jose.NewSigner(key, opts)
	if err != nil {
		err = errorsx.WithStack(err)
		return
	}

	// A explicit conversion from type alias MapClaims
	// to map[string]interface{} is required because the
	// go-jose CompactSerialize() only support explicit maps
	// as claims or structs but not type aliases from maps.
	claims := map[string]interface{}(t.Claims)
	rawToken, err = jwt.Signed(signer).Claims(claims).CompactSerialize()
	if err != nil {
		err = &ValidationError{Errors: ValidationErrorClaimsInvalid, Inner: err}
		return
	}
	return
}

func unsignedToken(t *Token) (string, error) {
	t.Header["alg"] = "none"
	if _, ok := t.Header[string(JWTHeaderType)]; !ok {
		t.Header[string(JWTHeaderType)] = JWTHeaderTypeValue
	}
	hbytes, err := json.Marshal(&t.Header)
	if err != nil {
		return "", errorsx.WithStack(err)
	}
	bbytes, err := json.Marshal(&t.Claims)
	if err != nil {
		return "", errorsx.WithStack(err)
	}
	h := base64.RawURLEncoding.EncodeToString(hbytes)
	b := base64.RawURLEncoding.EncodeToString(bbytes)
	return fmt.Sprintf("%v.%v.", h, b), nil
}

func newToken(parsedToken *jwt.JSONWebToken, claims MapClaims) (*Token, error) {
	token := &Token{Claims: claims}
	if len(parsedToken.Headers) != 1 {
		return nil, &ValidationError{text: fmt.Sprintf("only one header supported, got %v", len(parsedToken.Headers)), Errors: ValidationErrorMalformed}
	}

	// copy headers
	h := parsedToken.Headers[0]
	token.Header = map[string]interface{}{
		"alg": h.Algorithm,
	}
	if h.KeyID != "" {
		token.Header["kid"] = h.KeyID
	}
	for k, v := range h.ExtraHeaders {
		token.Header[string(k)] = v
	}

	token.Method = jose.SignatureAlgorithm(h.Algorithm)

	return token, nil
}

// Parse methods use this callback function to supply
// the key for verification.  The function receives the parsed,
// but unverified Token.  This allows you to use properties in the
// Header of the token (such as `kid`) to identify which key to use.
type Keyfunc func(*Token) (interface{}, error)

func Parse(tokenString string, keyFunc Keyfunc) (*Token, error) {
	return ParseWithClaims(tokenString, MapClaims{}, keyFunc)
}

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
// If everything is kosher, err will be nil
func ParseWithClaims(rawToken string, claims MapClaims, keyFunc Keyfunc) (*Token, error) {
	// Parse the token.
	parsedToken, err := jwt.ParseSigned(rawToken)
	if err != nil {
		return &Token{}, &ValidationError{Errors: ValidationErrorMalformed, text: err.Error()}
	}

	// fill unverified claims
	// This conversion is required because go-jose supports
	// only marshalling structs or maps but not alias types from maps
	//
	// The KeyFunc(*Token) function requires the claims to be set into the
	// Token, that is an unverified token, therefore an UnsafeClaimsWithoutVerification is done first
	// then with the returned key, the claims gets verified.
	if err := parsedToken.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return nil, &ValidationError{Errors: ValidationErrorClaimsInvalid, text: err.Error()}
	}

	// creates an usafe token
	token, err := newToken(parsedToken, claims)
	if err != nil {
		return nil, err
	}

	if keyFunc == nil {
		// keyFunc was not provided.  short circuiting validation
		return token, &ValidationError{Errors: ValidationErrorUnverifiable, text: "no Keyfunc was provided."}
	}

	// Call keyFunc callback to get verification key
	verificationKey, err := keyFunc(token)
	if err != nil {
		// keyFunc returned an error
		if ve, ok := err.(*ValidationError); ok {
			return token, ve
		}
		return token, &ValidationError{Errors: ValidationErrorUnverifiable, Inner: err}
	}
	if verificationKey == nil {
		return token, &ValidationError{Errors: ValidationErrorSignatureInvalid, text: "keyfunc returned a nil verification key"}
	}
	// To verify signature go-jose requires a pointer to
	// public key instead of the public key value.
	// The pointer values provides that pointer.
	// E.g. transform rsa.PublicKey -> *rsa.PublicKey
	verificationKey = pointer(verificationKey)

	// verify signature with returned key
	_, validNoneKey := verificationKey.(*unsafeNoneMagicConstant)
	isSignedToken := !(token.Method == SigningMethodNone && validNoneKey)
	if isSignedToken {
		if err := parsedToken.Claims(verificationKey, &claims); err != nil {
			return token, &ValidationError{Errors: ValidationErrorSignatureInvalid, text: err.Error()}
		}
	}

	// Validate claims
	// This validation is performed to be backwards compatible
	// with jwt-go library behavior
	if err := claims.Valid(); err != nil {
		if e, ok := err.(*ValidationError); !ok {
			err = &ValidationError{Inner: e, Errors: ValidationErrorClaimsInvalid}
		}
		return token, err
	}

	// set token as verified and validated
	token.valid = true
	return token, nil
}

// if underline value of v is not a pointer
// it creates a pointer of it and returns it
func pointer(v interface{}) interface{} {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		value := reflect.New(reflect.ValueOf(v).Type())
		value.Elem().Set(reflect.ValueOf(v))
		return value.Interface()
	}
	return v
}
