// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package trust

import (
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
)

func TestEmptyIssuerIsInvalid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:          "",
		Subject:         "valid-subject",
		AllowAnySubject: false,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	if err := v.Validate(r); err == nil {
		t.Error("an empty issuer should not be valid")
	}
}

func TestEmptySubjectAndNoAnySubjectFlagIsInvalid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "",
		AllowAnySubject: false,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	if err := v.Validate(r); err == nil {
		t.Error("an empty subject should not be valid")
	}
}

func TestEmptySubjectWithAnySubjectFlagIsValid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "",
		AllowAnySubject: true,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	if err := v.Validate(r); err != nil {
		t.Error("an empty subject with the allow_any_subject flag should be valid")
	}
}

func TestNonEmptySubjectWithAnySubjectFlagIsInvalid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "some-issuer",
		AllowAnySubject: true,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	if err := v.Validate(r); err == nil {
		t.Error("a non empty subject with the allow_any_subject flag should not be valid")
	}
}

func TestEmptyExpiresAtIsInvalid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "valid-subject",
		AllowAnySubject: false,
		ExpiresAt:       time.Time{},
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	if err := v.Validate(r); err == nil {
		t.Error("an empty expiration should not be valid")
	}
}

func TestEmptyPublicKeyIdIsInvalid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "valid-subject",
		AllowAnySubject: false,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "",
		},
	}

	if err := v.Validate(r); err == nil {
		t.Error("an empty public key id should not be valid")
	}
}

func TestIsValid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:          "valid-issuer",
		Subject:         "valid-subject",
		AllowAnySubject: false,
		ExpiresAt:       time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	if err := v.Validate(r); err != nil {
		t.Error("A request with an issuer, a subject, an expiration and a public key should be valid")
	}
}
