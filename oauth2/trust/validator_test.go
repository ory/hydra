package trust

import (
	"testing"
	"time"

	"gopkg.in/square/go-jose.v2"
)

func TestEmptyIssuerIsInvalid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:    "",
		Subject:   "valid-subject",
		Domain:    "",
		ExpiresAt: time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	if err := v.Validate(r); err == nil {
		t.Error("an empty issuer should not be valid")
	}
}

func TestEmptySubjectIsInvalid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:    "valid-issuer",
		Subject:   "",
		Domain:    "",
		ExpiresAt: time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	if err := v.Validate(r); err == nil {
		t.Error("an empty subject should not be valid")
	}
}

func TestEmptyExpiresAtIsInvalid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:    "valid-issuer",
		Subject:   "valid-subject",
		Domain:    "",
		ExpiresAt: time.Time{},
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
		Issuer:    "valid-issuer",
		Subject:   "valid-subject",
		Domain:    "",
		ExpiresAt: time.Now().Add(time.Hour * 10),
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
		Issuer:    "valid-issuer",
		Subject:   "valid-subject",
		Domain:    "",
		ExpiresAt: time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	if err := v.Validate(r); err != nil {
		t.Error("A request with an issuer, a subject, an expiration and a public key should be valid")
	}
}

func TestDomainIsValid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:    "valid-issuer",
		Subject:   "",
		Domain:    "ory.sh",
		ExpiresAt: time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	if err := v.Validate(r); err != nil {
		t.Error("A request with a domain restriction and an empty subject should be valid")
	}
}

func TestAnySubjectWithSubjectIsNotValid(t *testing.T) {
	v := GrantValidator{}

	r := createGrantRequest{
		Issuer:    "valid-issuer",
		Subject:   "valid-subject",
		Domain:    "ory.sh",
		ExpiresAt: time.Now().Add(time.Hour * 10),
		PublicKeyJWK: jose.JSONWebKey{
			KeyID: "valid-key-id",
		},
	}

	if err := v.Validate(r); err == nil {
		t.Error("A request with a domain restriction and a non-empty subject should not be valid")
	}
}
