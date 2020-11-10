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

	"github.com/pborman/uuid"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/x/stringslice"
	"github.com/ory/x/stringsx"
)

var (
	supportedAuthTokenSigningAlgs = []string{
		"RS256",
		"RS384",
		"RS512",
		"PS256",
		"PS384",
		"PS512",
		"ES256",
		"ES384",
		"ES512",
	}
)

type Validator struct {
	c    *http.Client
	conf Configuration
}

func NewValidator(conf Configuration) *Validator {
	return &Validator{
		c:    http.DefaultClient,
		conf: conf,
	}
}

func NewValidatorWithClient(conf Configuration, client *http.Client) *Validator {
	return &Validator{
		c:    client,
		conf: conf,
	}
}

func (v *Validator) Validate(c *Client) error {
	id := uuid.New()
	c.OutfacingID = stringsx.Coalesce(c.OutfacingID, id)

	if c.TokenEndpointAuthMethod == "" {
		c.TokenEndpointAuthMethod = "client_secret_basic"
	} else if c.TokenEndpointAuthMethod == "private_key_jwt" {
		if len(c.JSONWebKeysURI) == 0 && c.JSONWebKeys == nil {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHint("When token_endpoint_auth_method is 'private_key_jwt', either jwks or jwks_uri must be set."))
		}
		if c.TokenEndpointAuthSigningAlgorithm != "" && !isSupportedAuthTokenSigningAlg(c.TokenEndpointAuthSigningAlgorithm) {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Only RS256, RS384, RS512, PS256, PS384, PS512, ES256, ES384 and ES512 are supported as algorithms for private key authentication."))
		}
	}

	if len(c.JSONWebKeysURI) > 0 && c.JSONWebKeys != nil {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Fields jwks and jwks_uri can not both be set, you must choose one."))
	}

	if len(c.Secret) > 0 && len(c.Secret) < 6 {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Field client_secret must contain a secret that is at least 6 characters long."))
	}

	if len(c.Scope) == 0 {
		c.Scope = strings.Join(v.conf.DefaultClientScope(), " ")
	}

	for k, origin := range c.AllowedCORSOrigins {
		u, err := url.Parse(origin)
		if err != nil {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHintf("Origin URL %s from allowed_cors_origins could not be parsed: %s", origin, err))
		}

		if u.Scheme != "https" && u.Scheme != "http" {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHintf("Origin URL %s must use https:// or http:// as HTTP scheme.", origin))
		}

		if u.User != nil && len(u.User.String()) > 0 {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHintf("Origin URL %s has HTTP user and/or password set which is not allowed.", origin))
		}

		u.Path = strings.TrimRight(u.Path, "/")
		if len(u.Path)+len(u.RawQuery)+len(u.Fragment) > 0 {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHintf("Origin URL %s must have an empty path, query, and fragment but one of the parts is not empty.", origin))
		}

		c.AllowedCORSOrigins[k] = u.String()
	}

	// has to be 0 because it is not supposed to be set
	c.SecretExpiresAt = 0

	if len(c.SectorIdentifierURI) > 0 {
		if err := v.ValidateSectorIdentifierURL(c.SectorIdentifierURI, c.GetRedirectURIs()); err != nil {
			return err
		}
	}

	if c.UserinfoSignedResponseAlg == "" {
		c.UserinfoSignedResponseAlg = "none"
	}

	if c.UserinfoSignedResponseAlg != "none" && c.UserinfoSignedResponseAlg != "RS256" {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Field userinfo_signed_response_alg can either be 'none' or 'RS256'."))
	}

	var redirs []url.URL
	for _, r := range c.RedirectURIs {
		u, err := url.ParseRequestURI(r)
		if err != nil {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHintf("Unable to parse redirect URL: %s", r))
		}
		redirs = append(redirs, *u)

		if strings.Contains(r, "#") {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHint("Redirect URIs must not contain fragments (#)."))
		}
	}

	if c.SubjectType != "" {
		if !stringslice.Has(v.conf.SubjectTypesSupported(), c.SubjectType) {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHintf("Subject type %s is not supported by server, only %v are allowed.", c.SubjectType, v.conf.SubjectTypesSupported()))
		}
	} else {
		if stringslice.Has(v.conf.SubjectTypesSupported(), "public") {
			c.SubjectType = "public"
		} else {
			c.SubjectType = v.conf.SubjectTypesSupported()[0]
		}
	}

	for _, l := range c.PostLogoutRedirectURIs {
		u, err := url.ParseRequestURI(l)
		if err != nil {
			return errors.WithStack(fosite.ErrInvalidRequest.WithHintf("Unable to parse post_logout_redirect_uri: %s", l))
		}

		var found bool
		for _, r := range redirs {
			if r.Hostname() == u.Hostname() &&
				r.Port() == u.Port() &&
				r.Scheme == u.Scheme {
				found = true
			}
		}

		if !found {
			return errors.WithStack(fosite.ErrInvalidRequest.
				WithHintf(`post_logout_redirect_uri "%s" must match the domain, port, scheme of at least one of the registered redirect URIs but did not'`, l),
			)
		}
	}

	return nil
}

func (v *Validator) ValidateSectorIdentifierURL(location string, redirectURIs []string) error {
	l, err := url.Parse(location)
	if err != nil {
		return errors.WithStack(fosite.ErrInvalidRequest.WithHintf("Value of sector_identifier_uri could not be parsed because %s.", err))
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

func isSupportedAuthTokenSigningAlg(alg string) bool {
	for _, sAlg := range supportedAuthTokenSigningAlgs {
		if alg == sAlg {
			return true
		}
	}
	return false
}
