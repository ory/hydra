// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package trust

import (
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite"
)

func TestEmptyIssuerIsInvalid(t *testing.T) {
	r := createGrantRequest{
		Issuer:          "",
		Subject:         "valid-subject",
		AllowAnySubject: false,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	err := &fosite.RFC6749Error{}
	require.ErrorAs(t, validateGrant(r), &err)
	assert.Equal(t, "Field 'issuer' is required.", err.HintField)
}

func TestEmptySubjectAndNoAnySubjectFlagIsInvalid(t *testing.T) {
	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "",
		AllowAnySubject: false,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	err := &fosite.RFC6749Error{}
	require.ErrorAs(t, validateGrant(r), &err)
	assert.Equal(t, "One of 'subject' or 'allow_any_subject' field must be set.", err.HintField)
}

func TestEmptySubjectWithAnySubjectFlagIsValid(t *testing.T) {
	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "",
		AllowAnySubject: true,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	assert.NoError(t, validateGrant(r))
}

func TestNonEmptySubjectWithAnySubjectFlagIsInvalid(t *testing.T) {
	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "some-issuer",
		AllowAnySubject: true,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	err := &fosite.RFC6749Error{}
	require.ErrorAs(t, validateGrant(r), &err)
	assert.Equal(t, "Both 'subject' and 'allow_any_subject' fields cannot be set at the same time.", err.HintField)
}

func TestEmptyExpiresAtIsInvalid(t *testing.T) {
	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "valid-subject",
		AllowAnySubject: false,
		ExpiresAt:       time.Time{},
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	err := &fosite.RFC6749Error{}
	require.ErrorAs(t, validateGrant(r), &err)
	assert.Equal(t, "Field 'expires_at' is required.", err.HintField)
}

func TestEmptyPublicKeyIdIsInvalid(t *testing.T) {
	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "valid-subject",
		AllowAnySubject: false,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "",
		},
	}

	err := &fosite.RFC6749Error{}
	require.ErrorAs(t, validateGrant(r), &err)
	assert.Equal(t, "Field 'jwk' must contain JWK with kid header.", err.HintField)
}

func TestIsValid(t *testing.T) {
	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "valid-subject",
		AllowAnySubject: false,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	assert.NoError(t, validateGrant(r))
}
