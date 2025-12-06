// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ory/hydra/v2/fosite/internal/gen"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnsignedToken(t *testing.T) {
	var testCases = []struct {
		name         string
		jwtHeaders   map[string]interface{}
		expectedType string
	}{
		{
			name:         "set JWT as 'typ' when the the type is not specified in the headers",
			jwtHeaders:   map[string]interface{}{},
			expectedType: "JWT",
		},
		{
			name:         "'typ' set explicitly",
			jwtHeaders:   map[string]interface{}{"typ": "at+jwt"},
			expectedType: "at+jwt",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key := UnsafeAllowNoneSignatureType
			token := NewWithClaims(SigningMethodNone, MapClaims{
				"aud": "foo",
				"exp": time.Now().UTC().Add(time.Hour).Unix(),
				"iat": time.Now().UTC().Unix(),
				"sub": "nestor",
			})
			token.Header = tc.jwtHeaders
			rawToken, err := token.SignedString(key)
			require.NoError(t, err)
			require.NotEmpty(t, rawToken)
			parts := strings.Split(rawToken, ".")
			require.Len(t, parts, 3)
			require.Empty(t, parts[2])
			tk, err := jwt.ParseSigned(rawToken)
			require.NoError(t, err)
			require.Len(t, tk.Headers, 1)
			require.Equal(t, tc.expectedType, tk.Headers[0].ExtraHeaders[jose.HeaderKey("typ")])
		})
	}
}

func TestJWTHeaders(t *testing.T) {
	var testCases = []struct {
		name         string
		jwtHeaders   map[string]interface{}
		expectedType string
	}{
		{
			name:         "set JWT as 'typ' when the the type is not specified in the headers",
			jwtHeaders:   map[string]interface{}{},
			expectedType: "JWT",
		},
		{
			name:         "'typ' set explicitly",
			jwtHeaders:   map[string]interface{}{"typ": "at+jwt"},
			expectedType: "at+jwt",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rawToken := makeSampleTokenWithCustomHeaders(nil, jose.RS256, tc.jwtHeaders, gen.MustRSAKey())
			tk, err := jwt.ParseSigned(rawToken)
			require.NoError(t, err)
			require.Len(t, tk.Headers, 1)
			require.Equal(t, tk.Headers[0].Algorithm, "RS256")
			require.Equal(t, tc.expectedType, tk.Headers[0].ExtraHeaders[jose.HeaderKey("typ")])
		})
	}
}

var keyFuncError error = fmt.Errorf("error loading key")
var (
	jwtTestDefaultKey *rsa.PublicKey = parseRSAPublicKeyFromPEM(defaultPubKeyPEM)
	defaultKeyFunc    Keyfunc        = func(t *Token) (interface{}, error) { return jwtTestDefaultKey, nil }
	emptyKeyFunc      Keyfunc        = func(t *Token) (interface{}, error) { return nil, nil }
	errorKeyFunc      Keyfunc        = func(t *Token) (interface{}, error) { return nil, keyFuncError }
	nilKeyFunc        Keyfunc        = nil
)

// Many test cases where taken from https://github.com/dgrijalva/jwt-go/blob/master/parser_test.go
// Test cases related to json.Number where excluded because that is not supported by go-jose,
// it is not used in fosite and therefore not supported.
func TestParser_Parse(t *testing.T) {
	var (
		defaultES256PrivateKey = gen.MustES256Key()
		defaultSigningKey      = parseRSAPrivateKeyFromPEM(defaultPrivateKeyPEM)
		publicECDSAKey         = func(*Token) (interface{}, error) { return &defaultES256PrivateKey.PublicKey, nil }
		noneKey                = func(*Token) (interface{}, error) { return UnsafeAllowNoneSignatureType, nil }
		randomKey              = func(*Token) (interface{}, error) {
			k, err := rsa.GenerateKey(rand.Reader, 2048)
			require.NoError(t, err)
			return &k.PublicKey, nil
		}
	)
	type expected struct {
		errors  uint32
		keyFunc Keyfunc
		valid   bool
		claims  MapClaims
	}
	type generate struct {
		claims     MapClaims
		signingKey interface{}             // defaultSigningKey
		method     jose.SignatureAlgorithm // default RS256
	}
	type given struct {
		name        string
		tokenString string
		generate    *generate
	}
	var jwtTestData = []struct {
		expected
		given
	}{
		{
			given: given{
				name:        "basic",
				tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
			},
			expected: expected{
				keyFunc: defaultKeyFunc,
				claims:  MapClaims{"foo": "bar"},
				valid:   true,
				errors:  0,
			},
		},
		{
			given: given{
				name: "basic expired",
				generate: &generate{
					claims: MapClaims{"foo": "bar", "exp": time.Now().Unix() - 100},
				},
			},
			expected: expected{
				keyFunc: defaultKeyFunc,
				claims:  MapClaims{"foo": "bar", "exp": time.Now().Unix() - 100},
				valid:   false,
				errors:  ValidationErrorExpired,
			},
		},
		{
			given: given{
				name: "basic nbf",
				generate: &generate{
					claims: MapClaims{"foo": "bar", "nbf": time.Now().Unix() + 100},
				},
			},
			expected: expected{
				keyFunc: defaultKeyFunc,
				claims:  MapClaims{"foo": "bar", "nbf": time.Now().Unix() + 100},
				valid:   false,
				errors:  ValidationErrorNotValidYet,
			},
		},
		{
			given: given{
				name: "expired and nbf",
				generate: &generate{
					claims: MapClaims{"foo": "bar", "nbf": time.Now().Unix() + 100, "exp": time.Now().Unix() - 100},
				},
			},
			expected: expected{
				keyFunc: defaultKeyFunc,
				claims:  MapClaims{"foo": "bar", "nbf": time.Now().Unix() + 100, "exp": time.Now().Unix() - 100},
				valid:   false,
				errors:  ValidationErrorNotValidYet | ValidationErrorExpired,
			},
		},
		{
			given: given{
				name:        "basic invalid",
				tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.EhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
			},
			expected: expected{
				keyFunc: defaultKeyFunc,
				claims:  MapClaims{"foo": "bar"},
				valid:   false,
				errors:  ValidationErrorSignatureInvalid,
			},
		},
		{
			given: given{
				name:        "basic nokeyfunc",
				tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
			},
			expected: expected{
				keyFunc: nilKeyFunc,
				claims:  MapClaims{"foo": "bar"},
				valid:   false,
				errors:  ValidationErrorUnverifiable,
			},
		},
		{
			given: given{
				name:        "basic nokey",
				tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
			},
			expected: expected{
				keyFunc: emptyKeyFunc,
				claims:  MapClaims{"foo": "bar"},
				valid:   false,
				errors:  ValidationErrorSignatureInvalid,
			},
		},
		{
			given: given{
				name:        "basic errorkey",
				tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
				generate: &generate{
					claims: MapClaims{"foo": "bar"},
				},
			},
			expected: expected{
				keyFunc: errorKeyFunc,
				claims:  MapClaims{"foo": "bar"},
				valid:   false,
				errors:  ValidationErrorUnverifiable,
			},
		},
		{
			given: given{
				name: "valid signing method",
				generate: &generate{
					claims: MapClaims{"foo": "bar"},
				},
			},
			expected: expected{
				keyFunc: defaultKeyFunc,
				claims:  MapClaims{"foo": "bar"},
				valid:   true,
				errors:  0,
			},
		},
		{
			given: given{
				name:        "invalid",
				tokenString: "foo_invalid_token",
			},
			expected: expected{
				keyFunc: defaultKeyFunc,
				claims:  MapClaims(nil),
				valid:   false,
				errors:  ValidationErrorMalformed,
			},
		},
		{
			given: given{
				name:        "valid format invalid content",
				tokenString: "foo.bar.baz",
			},
			expected: expected{
				keyFunc: defaultKeyFunc,
				claims:  MapClaims(nil),
				valid:   false,
				errors:  ValidationErrorMalformed,
			},
		},
		{
			given: given{
				name:        "wrong key, expected ECDSA got RSA",
				tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
			},
			expected: expected{
				keyFunc: publicECDSAKey,
				claims:  MapClaims{"foo": "bar"},
				valid:   false,
				errors:  ValidationErrorSignatureInvalid,
			},
		},
		{
			given: given{
				name:        "should fail, got RSA but found no key",
				tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
			},
			expected: expected{
				keyFunc: emptyKeyFunc,
				claims:  MapClaims{"foo": "bar"},
				valid:   false,
				errors:  ValidationErrorSignatureInvalid,
			},
		},
		{
			given: given{
				name:        "key does not match",
				tokenString: "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
			},
			expected: expected{
				keyFunc: randomKey,
				claims:  MapClaims{"foo": "bar"},
				valid:   false,
				errors:  ValidationErrorSignatureInvalid,
			},
		},
		{
			given: given{
				name: "used before issued",
				generate: &generate{
					claims: MapClaims{"foo": "bar", "iat": time.Now().Unix() + 500},
				},
			},
			expected: expected{
				keyFunc: defaultKeyFunc,
				claims:  MapClaims{"foo": "bar", "iat": time.Now().Unix() + 500},
				valid:   false,
				errors:  ValidationErrorIssuedAt,
			},
		},
		{
			given: given{
				name: "valid ECDSA signing method",
				generate: &generate{
					claims:     MapClaims{"foo": "bar"},
					signingKey: defaultES256PrivateKey,
					method:     jose.ES256,
				},
			},
			expected: expected{
				keyFunc: publicECDSAKey,
				claims:  MapClaims{"foo": "bar"},
				valid:   true,
				errors:  0,
			},
		},
		{
			given: given{
				name: "should pass, valid NONE signing method",
				generate: &generate{
					claims:     MapClaims{"foo": "bar"},
					signingKey: UnsafeAllowNoneSignatureType,
					method:     SigningMethodNone,
				},
			},
			expected: expected{
				keyFunc: noneKey,
				claims:  MapClaims{"foo": "bar"},
				valid:   true,
				errors:  0,
			},
		},
		{
			given: given{
				name: "should fail, expected RS256 but got NONE",
				generate: &generate{
					claims:     MapClaims{"foo": "bar"},
					signingKey: UnsafeAllowNoneSignatureType,
					method:     SigningMethodNone,
				},
			},
			expected: expected{
				keyFunc: defaultKeyFunc,
				claims:  MapClaims{"foo": "bar"},
				valid:   false,
				errors:  ValidationErrorSignatureInvalid,
			},
		},
		{
			given: given{
				name: "should fail, expected ECDSA but got NONE",
				generate: &generate{
					claims:     MapClaims{"foo": "bar"},
					signingKey: UnsafeAllowNoneSignatureType,
					method:     SigningMethodNone,
				},
			},
			expected: expected{
				keyFunc: publicECDSAKey,
				claims:  MapClaims{"foo": "bar"},
				valid:   false,
				errors:  ValidationErrorSignatureInvalid,
			},
		},
	}

	// Iterate over test data set and run tests
	for _, data := range jwtTestData {
		t.Run(data.name, func(t *testing.T) {
			if data.generate != nil {
				signingKey := data.generate.signingKey
				method := data.generate.method
				if signingKey == nil {
					// use test defaults
					signingKey = defaultSigningKey
					method = jose.RS256
				}
				data.tokenString = makeSampleToken(data.generate.claims, method, signingKey)
			}

			// Parse the token
			var token *Token
			var err error

			// Figure out correct claims type
			token, err = ParseWithClaims(data.tokenString, MapClaims{}, data.keyFunc)
			// Verify result matches expectation
			assert.EqualValues(t, data.claims, token.Claims)
			if data.valid && err != nil {
				t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
			}

			if !data.valid && err == nil {
				t.Errorf("[%v] Invalid token passed validation", data.name)
			}

			if (err == nil && !token.Valid()) || (err != nil && token.Valid()) {
				t.Errorf("[%v] Inconsistent behavior between returned error and token.Valid", data.name)
			}

			if data.errors != 0 {
				if err == nil {
					t.Errorf("[%v] Expecting error.  Didn't get one.", data.name)
				} else {

					ve := err.(*ValidationError)
					// compare the bitfield part of the error
					if e := ve.Errors; e != data.errors {
						t.Errorf("[%v] Errors don't match expectation.  %v != %v", data.name, e, data.errors)
					}

					if err.Error() == keyFuncError.Error() && ve.Inner != keyFuncError {
						t.Errorf("[%v] Inner error does not match expectation.  %v != %v", data.name, ve.Inner, keyFuncError)
					}
				}
			}
		})
	}
}

func makeSampleToken(c MapClaims, m jose.SignatureAlgorithm, key interface{}) string {
	token := NewWithClaims(m, c)
	s, e := token.SignedString(key)

	if e != nil {
		panic(e.Error())
	}

	return s
}

func makeSampleTokenWithCustomHeaders(c MapClaims, m jose.SignatureAlgorithm, headers map[string]interface{}, key interface{}) string {
	token := NewWithClaims(m, c)
	token.Header = headers
	s, e := token.SignedString(key)

	if e != nil {
		panic(e.Error())
	}

	return s
}

func parseRSAPublicKeyFromPEM(key []byte) *rsa.PublicKey {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		panic("not possible to decode")
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			panic(err)
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		panic("not an *rsa.PublicKey")
	}

	return pkey
}

func parseRSAPrivateKeyFromPEM(key []byte) *rsa.PrivateKey {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		panic("unable to decode")
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			panic(err)
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		panic("not an rsa private key")
	}

	return pkey
}

var (
	defaultPubKeyPEM = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----`)
	defaultPrivateKeyPEM = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----`)
)
