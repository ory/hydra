// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ory/fosite"
)

// Request a Verifiable Credential
//
// swagger:parameters createVerifiableCredential
//
//lint:ignore U1000 Used to generate Swagger and OpenAPI definitions
type createVerifiableCredentialRequest struct {
	// in: body
	Body CreateVerifiableCredentialRequestBody
}

// CreateVerifiableCredentialRequestBody contains the request body to request a verifiable credential.
//
// swagger:parameters createVerifiableCredentialRequestBody
type CreateVerifiableCredentialRequestBody struct {
	Format string                     `json:"format"`
	Types  []string                   `json:"types"`
	Proof  *VerifiableCredentialProof `json:"proof"`
}

// VerifiableCredentialProof contains the proof of a verifiable credential.
//
// swagger:parameters verifiableCredentialProof
type VerifiableCredentialProof struct {
	ProofType string `json:"proof_type"`
	JWT       string `json:"jwt"`
}

// VerifiableCredentialResponse contains the verifiable credential.
//
// swagger:model verifiableCredentialResponse
type VerifiableCredentialResponse struct {
	Format     string `json:"format"`
	Credential string `json:"credential_draft_00"`
}

// VerifiableCredentialPrimingResponse contains the nonce to include in the proof-of-possession JWT.
//
// swagger:model verifiableCredentialPrimingResponse
type VerifiableCredentialPrimingResponse struct {
	Format         string `json:"format"`
	Nonce          string `json:"c_nonce"`
	NonceExpiresIn int64  `json:"c_nonce_expires_in"`

	fosite.RFC6749ErrorJson
}

type VerifableCredentialClaims struct {
	jwt.RegisteredClaims
	VerifiableCredential VerifiableCredentialClaim `json:"vc"`
}
type VerifiableCredentialClaim struct {
	Context []string       `json:"@context"`
	Subject map[string]any `json:"credentialSubject"`
	Type    []string       `json:"type"`
}

func (v *VerifableCredentialClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}

func (v *VerifableCredentialClaims) ToMapClaims() (res map[string]any, err error) {
	res = map[string]any{}

	bs, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bs, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
