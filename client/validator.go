/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ory/fosite"
	"github.com/ory/go-convenience/stringslice"
	"github.com/ory/go-convenience/stringsx"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

type Validator struct {
	c                   *http.Client
	DefaultClientScopes []string
	SubjectTypes        []string
}

func NewValidator(
	defaultClientScopes []string,
	subjectTypes []string,
) *Validator {
	if len(subjectTypes) == 0 {
		subjectTypes = []string{"public"}
	}

	subjectTypes = stringslice.Filter(subjectTypes, func(s string) bool {
		return !(s == "public" || s == "pairwise")
	})

	return &Validator{
		c:                   http.DefaultClient,
		DefaultClientScopes: defaultClientScopes,
		SubjectTypes:        subjectTypes,
	}
}

func (v *Validator) Validate(c *Client) error {
	id := uuid.New()
	c.ClientID = stringsx.Coalesce(c.ClientID, id)

	if c.TokenEndpointAuthMethod == "" {
		c.TokenEndpointAuthMethod = "client_secret_basic"
	} else {
		if len(c.JSONWebKeysURI) == 0 && c.JSONWebKeys == nil && c.TokenEndpointAuthMethod == "private_key_jwt" {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHint("When token_endpoint_auth_method is \"private_key_jwt\", either jwks or jwks_uri must be set."))
		}
	}

	if len(c.JSONWebKeysURI) > 0 && c.JSONWebKeys != nil {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Fields jwks and jwks_uri can not both be set, you must choose one."))
	}

	if len(c.Secret) > 0 && len(c.Secret) < 6 {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Field client_secret must contain a secret that is at least 6 characters long."))
	}

	if len(c.Scope) == 0 {
		c.Scope = strings.Join(v.DefaultClientScopes, " ")
	}

	// has to be 0 because it is not supposed to be set
	c.SecretExpiresAt = 0

	if len(c.SectorIdentifierURI) > 0 {
		if err := v.validateSectorIdentifierURL(c.SectorIdentifierURI, c.GetRedirectURIs()); err != nil {
			return err
		}
	}

	if c.UserinfoSignedResponseAlg == "" {
		c.UserinfoSignedResponseAlg = "none"
	}

	if c.UserinfoSignedResponseAlg != "none" && c.UserinfoSignedResponseAlg != "RS256" {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Field userinfo_signed_response_alg can either be \"none\" or \"RS256\"."))
	}

	for _, r := range c.RedirectURIs {
		if strings.Contains(r, "#") {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Redirect URIs must not contain fragments (#)"))
		}
	}

	if c.SubjectType != "" {
		if !stringslice.Has(v.SubjectTypes, c.SubjectType) {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHint(fmt.Sprintf("Subject type %s is not supported by server, only %v are allowed.", c.SubjectType, v.SubjectTypes)))
		}
	} else {
		if !stringslice.Has(v.SubjectTypes, "public") {
			c.SubjectType = "public"
		} else {
			c.SubjectType = v.SubjectTypes[0]
		}
	}

	return nil
}

func (v *Validator) validateSectorIdentifierURL(location string, redirectURIs []string) error {
	l, err := url.Parse(location)
	if err != nil {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint(fmt.Sprintf("Value of sector_identifier_uri could not be parsed: %s", err)))
	}

	if l.Scheme != "https" {
		return errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Value sector_identifier_uri must be an HTTPS URL but it is not."))
	}

	response, err := v.c.Get(location)
	if err != nil {
		return errors.WithStack(fosite.ErrInvalidRequest.WithDebug(fmt.Sprintf("Unable to connect to URL set by sector_identifier_uri: %s", err)))
	}
	defer response.Body.Close()

	var urls []string
	if err := json.NewDecoder(response.Body).Decode(&urls); err != nil {
		return errors.WithStack(fosite.ErrInvalidRequest.WithDebug(fmt.Sprintf("Unable to decode values from sector_identifier_uri: %s", err)))
	}

	if len(urls) == 0 {
		return errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Array from sector_identifier_uri contains no items"))
	}

	for _, r := range redirectURIs {
		if !stringslice.Has(urls, r) {
			return errors.WithStack(fosite.ErrInvalidRequest.WithDebug(fmt.Sprintf("Redirect URL \"%s\" does not match values from sector_identifier_uri.", r)))
		}
	}

	return nil
}
