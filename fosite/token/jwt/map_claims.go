// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwt

import (
	"bytes"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"time"

	jjson "github.com/go-jose/go-jose/v3/json"

	"github.com/ory/x/errorsx"
)

var TimeFunc = time.Now

// MapClaims provides backwards compatible validations not available in `go-jose`.
// It was taken from [here](https://raw.githubusercontent.com/form3tech-oss/jwt-go/master/map_claims.go).
//
// Claims type that uses the map[string]interface{} for JSON decoding
// This is the default claims type if you don't supply one
type MapClaims map[string]interface{}

// Compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m MapClaims) VerifyAudience(cmp string, req bool) bool {
	var aud []string
	switch v := m["aud"].(type) {
	case []string:
		aud = v
	case []interface{}:
		for _, a := range v {
			vs, ok := a.(string)
			if !ok {
				return false
			}
			aud = append(aud, vs)
		}
	case string:
		aud = append(aud, v)
	default:
		return false
	}
	return verifyAud(aud, cmp, req)
}

// Compares the exp claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m MapClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	if v, ok := m.toInt64("exp"); ok {
		return verifyExp(v, cmp, req)
	}
	return !req
}

// Compares the iat claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m MapClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	if v, ok := m.toInt64("iat"); ok {
		return verifyIat(v, cmp, req)
	}
	return !req
}

// Compares the iss claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m MapClaims) VerifyIssuer(cmp string, req bool) bool {
	iss, _ := m["iss"].(string)
	return verifyIss(iss, cmp, req)
}

// Compares the nbf claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (m MapClaims) VerifyNotBefore(cmp int64, req bool) bool {
	if v, ok := m.toInt64("nbf"); ok {
		return verifyNbf(v, cmp, req)
	}

	return !req
}

func (m MapClaims) toInt64(claim string) (int64, bool) {
	switch t := m[claim].(type) {
	case float64:
		return int64(t), true
	case int64:
		return t, true
	case json.Number:
		v, err := t.Int64()
		if err == nil {
			return v, true
		}
		vf, err := t.Float64()
		if err != nil {
			return 0, false
		}

		return int64(vf), true
	}
	return 0, false
}

// Validates time based claims "exp, iat, nbf".
// There is no accounting for clock skew.
// As well, if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (m MapClaims) Valid() error {
	vErr := new(ValidationError)
	now := TimeFunc().Unix()

	if !m.VerifyExpiresAt(now, false) {
		vErr.Inner = errors.New("Token is expired")
		vErr.Errors |= ValidationErrorExpired
	}

	if !m.VerifyIssuedAt(now, false) {
		vErr.Inner = errors.New("Token used before issued")
		vErr.Errors |= ValidationErrorIssuedAt
	}

	if !m.VerifyNotBefore(now, false) {
		vErr.Inner = errors.New("Token is not valid yet")
		vErr.Errors |= ValidationErrorNotValidYet
	}

	if vErr.valid() {
		return nil
	}

	return vErr
}

func (m MapClaims) UnmarshalJSON(b []byte) error {
	// This custom unmarshal allows to configure the
	// go-jose decoding settings since there is no other way
	// see https://github.com/square/go-jose/issues/353.
	// If issue is closed with a better solution
	// this custom Unmarshal method can be removed
	d := jjson.NewDecoder(bytes.NewReader(b))
	mp := map[string]interface{}(m)
	d.SetNumberType(jjson.UnmarshalIntOrFloat)
	if err := d.Decode(&mp); err != nil {
		return errorsx.WithStack(err)
	}

	return nil
}

func verifyAud(aud []string, cmp string, required bool) bool {
	if len(aud) == 0 {
		return !required
	}

	for _, a := range aud {
		if subtle.ConstantTimeCompare([]byte(a), []byte(cmp)) != 0 {
			return true
		}
	}
	return false
}

func verifyExp(exp int64, now int64, required bool) bool {
	if exp == 0 {
		return !required
	}
	return now <= exp
}

func verifyIat(iat int64, now int64, required bool) bool {
	if iat == 0 {
		return !required
	}
	return now >= iat
}

func verifyIss(iss string, cmp string, required bool) bool {
	if iss == "" {
		return !required
	}
	if subtle.ConstantTimeCompare([]byte(iss), []byte(cmp)) != 0 {
		return true
	} else {
		return false
	}
}

func verifyNbf(nbf int64, now int64, required bool) bool {
	if nbf == 0 {
		return !required
	}
	return now >= nbf
}
