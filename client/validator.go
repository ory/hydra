// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"

	"github.com/ory/herodot"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/ipx"
)

var supportedAuthTokenSigningAlgs = []string{
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

func isSupportedAuthTokenSigningAlg(alg string) bool {
	return slices.Contains(supportedAuthTokenSigningAlgs, alg)
}

type validatorRegistry interface {
	x.HTTPClientProvider
	config.Provider
}

type Validator struct {
	r validatorRegistry
}

func NewValidator(r validatorRegistry) *Validator {
	return &Validator{r: r}
}

func (v *Validator) Validate(ctx context.Context, c *Client) error {
	if c.TokenEndpointAuthMethod == "" {
		c.TokenEndpointAuthMethod = "client_secret_basic"
	} else if c.TokenEndpointAuthMethod == "private_key_jwt" {
		if len(c.JSONWebKeysURI) == 0 && c.GetJSONWebKeys() == nil {
			return errors.WithStack(ErrInvalidClientMetadata.WithHint("When token_endpoint_auth_method is 'private_key_jwt', either jwks or jwks_uri must be set."))
		}
		if c.TokenEndpointAuthSigningAlgorithm != "" && !isSupportedAuthTokenSigningAlg(c.TokenEndpointAuthSigningAlgorithm) {
			return errors.WithStack(ErrInvalidClientMetadata.WithHint("Only RS256, RS384, RS512, PS256, PS384, PS512, ES256, ES384 and ES512 are supported as algorithms for private key authentication."))
		}
	}

	if len(c.JSONWebKeysURI) > 0 && c.GetJSONWebKeys() != nil {
		return errors.WithStack(ErrInvalidClientMetadata.WithHint("Fields jwks and jwks_uri can not both be set, you must choose one."))
	}

	if jsonWebKeys := c.GetJSONWebKeys(); jsonWebKeys != nil {
		for _, k := range jsonWebKeys.Keys {
			if !k.Valid() {
				return errors.WithStack(ErrInvalidClientMetadata.WithHint("Invalid JSON web key in set."))
			}
		}
	}

	if v.r.Config().ClientHTTPNoPrivateIPRanges() {
		values := map[string]string{
			"jwks_uri":               c.JSONWebKeysURI,
			"backchannel_logout_uri": c.BackChannelLogoutURI,
		}

		for k, v := range c.RequestURIs {
			values[fmt.Sprintf("request_uris.%d", k)] = v
		}

		if err := ipx.AreAllAssociatedIPsAllowed(values); err != nil {
			return errors.WithStack(ErrInvalidClientMetadata.WithHintf("Client IP address is not allowed: %s", err))
		}
	}

	if c.TermsOfServiceURI != "" {
		u, err := url.ParseRequestURI(c.TermsOfServiceURI)
		if err != nil {
			return errors.WithStack(ErrInvalidClientMetadata.WithHint("Field tos_uri must be a valid URI."))
		}

		if u.Scheme != "https" && u.Scheme != "http" {
			return errors.WithStack(ErrInvalidClientMetadata.WithHintf("tos_uri %s must use https:// or http:// as HTTP scheme.", c.TermsOfServiceURI))
		}

	}

	if len(c.Secret) > 0 && len(c.Secret) < 6 {
		return errors.WithStack(ErrInvalidClientMetadata.WithHint("Field client_secret must contain a secret that is at least 6 characters long."))
	}

	if len(c.Scope) == 0 {
		c.Scope = strings.Join(v.r.Config().DefaultClientScope(ctx), " ")
	}

	for k, origin := range c.AllowedCORSOrigins {
		u, err := url.Parse(origin)
		if err != nil {
			return errors.WithStack(ErrInvalidClientMetadata.WithHintf("Origin URL %s from allowed_cors_origins could not be parsed: %s", origin, err))
		}

		if u.Scheme != "https" && u.Scheme != "http" {
			return errors.WithStack(ErrInvalidClientMetadata.WithHintf("Origin URL %s must use https:// or http:// as HTTP scheme.", origin))
		}

		if u.User != nil && len(u.User.String()) > 0 {
			return errors.WithStack(ErrInvalidClientMetadata.WithHintf("Origin URL %s has HTTP user and/or password set which is not allowed.", origin))
		}

		u.Path = strings.TrimRight(u.Path, "/")
		if len(u.Path)+len(u.RawQuery)+len(u.Fragment) > 0 {
			return errors.WithStack(ErrInvalidClientMetadata.WithHintf("Origin URL %s must have an empty path, query, and fragment but one of the parts is not empty.", origin))
		}

		c.AllowedCORSOrigins[k] = u.String()
	}

	// has to be 0 because it is not supposed to be set
	c.SecretExpiresAt = 0

	if len(c.SectorIdentifierURI) > 0 {
		if err := v.ValidateSectorIdentifierURL(ctx, c.SectorIdentifierURI, c.GetRedirectURIs()); err != nil {
			return err
		}
	}

	if c.UserinfoSignedResponseAlg == "" {
		c.UserinfoSignedResponseAlg = "none"
	}

	if c.UserinfoSignedResponseAlg != "none" && c.UserinfoSignedResponseAlg != "RS256" {
		return errors.WithStack(ErrInvalidClientMetadata.WithHint("Field userinfo_signed_response_alg can either be 'none' or 'RS256'."))
	}

	redirs := make([]*url.URL, len(c.RedirectURIs))
	for i, r := range c.RedirectURIs {
		if strings.Contains(r, "#") {
			return errors.WithStack(ErrInvalidRedirectURI.WithHint("Redirect URIs must not contain fragments (#)."))
		}
		var err error
		redirs[i], err = url.ParseRequestURI(r)
		if err != nil {
			return errors.WithStack(ErrInvalidRedirectURI.WithHintf("Unable to parse redirect URL: %s", r))
		}
	}

	if c.SubjectType != "" {
		if !slices.Contains(v.r.Config().SubjectTypesSupported(ctx, c), c.SubjectType) {
			return errors.WithStack(ErrInvalidClientMetadata.WithHintf("Subject type %s is not supported by server, only %v are allowed.", c.SubjectType, v.r.Config().SubjectTypesSupported(ctx, c)))
		}
	} else {
		supportedTypes := v.r.Config().SubjectTypesSupported(ctx, c)
		if slices.Contains(supportedTypes, "public") {
			c.SubjectType = "public"
		} else {
			c.SubjectType = supportedTypes[0]
		}
	}

	for _, l := range c.PostLogoutRedirectURIs {
		u, err := url.ParseRequestURI(l)
		if err != nil {
			return errors.WithStack(ErrInvalidClientMetadata.WithHintf("Unable to parse post_logout_redirect_uri: %s", l))
		}

		if !slices.ContainsFunc(redirs, func(r *url.URL) bool {
			return r.Scheme == u.Scheme && r.Hostname() == u.Hostname() && r.Port() == u.Port()
		}) {
			return errors.WithStack(ErrInvalidClientMetadata.
				WithHintf("post_logout_redirect_uri %q must match the domain, port, scheme of at least one of the registered redirect URIs but did not", l),
			)
		}
	}

	if c.AccessTokenStrategy != "" {
		s, err := config.ToAccessTokenStrategyType(c.AccessTokenStrategy)
		if err != nil {
			return errors.WithStack(ErrInvalidClientMetadata.
				WithHintf("invalid access token strategy: %v", err))
		}
		// Canonicalize, just in case.
		c.AccessTokenStrategy = string(s)
	}

	return nil
}

func (v *Validator) ValidateDynamicRegistration(ctx context.Context, c *Client) error {
	if c.Metadata != nil {
		return errors.WithStack(ErrInvalidClientMetadata.WithHint(`"metadata" cannot be set for dynamic client registration`))
	}
	if c.AccessTokenStrategy != "" {
		return errors.WithStack(herodot.ErrBadRequest.WithReasonf("It is not allowed to choose your own access token strategy."))
	}
	if c.SkipConsent {
		return errors.WithStack(ErrInvalidRequest.WithDescription(`"skip_consent" cannot be set for dynamic client registration`))
	}
	if c.SkipLogoutConsent.Bool {
		return errors.WithStack(ErrInvalidRequest.WithDescription(`"skip_logout_consent" cannot be set for dynamic client registration`))
	}

	return v.Validate(ctx, c)
}

func (v *Validator) ValidateSectorIdentifierURL(ctx context.Context, location string, redirectURIs []string) error {
	l, err := url.Parse(location)
	if err != nil {
		return errors.WithStack(ErrInvalidClientMetadata.WithHintf("Value of sector_identifier_uri could not be parsed because %s.", err))
	}

	if l.Scheme != "https" {
		return errors.WithStack(ErrInvalidClientMetadata.WithDebug("Value sector_identifier_uri must be an HTTPS URL but it is not."))
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, "GET", location, nil)
	if err != nil {
		return errors.WithStack(ErrInvalidClientMetadata.WithDebugf("Value sector_identifier_uri must be an HTTPS URL but it is not: %s", err.Error()))
	}
	response, err := v.r.HTTPClient(ctx).Do(req)
	if err != nil {
		return errors.WithStack(ErrInvalidClientMetadata.WithDebug(fmt.Sprintf("Unable to connect to URL set by sector_identifier_uri: %s", err)))
	}
	defer response.Body.Close() //nolint:errcheck
	response.Body = io.NopCloser(io.LimitReader(response.Body, 5<<20 /* 5 MiB */))

	var urls []string
	if err := json.NewDecoder(response.Body).Decode(&urls); err != nil {
		return errors.WithStack(ErrInvalidClientMetadata.WithDebug(fmt.Sprintf("Unable to decode values from sector_identifier_uri: %s", err)))
	}

	if len(urls) == 0 {
		return errors.WithStack(ErrInvalidClientMetadata.WithDebug("Array from sector_identifier_uri contains no items"))
	}

	for _, r := range redirectURIs {
		if !slices.Contains(urls, r) {
			return errors.WithStack(ErrInvalidClientMetadata.WithDebug(fmt.Sprintf("Redirect URL \"%s\" does not match values from sector_identifier_uri.", r)))
		}
	}

	return nil
}
