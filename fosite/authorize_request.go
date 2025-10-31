// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"net/url"
)

type ResponseModeType string

const (
	ResponseModeDefault  = ResponseModeType("")
	ResponseModeFormPost = ResponseModeType("form_post")
	ResponseModeQuery    = ResponseModeType("query")
	ResponseModeFragment = ResponseModeType("fragment")
)

// AuthorizeRequest is an implementation of AuthorizeRequester
type AuthorizeRequest struct {
	ResponseTypes        Arguments        `json:"responseTypes" gorethink:"responseTypes"`
	RedirectURI          *url.URL         `json:"redirectUri" gorethink:"redirectUri"`
	State                string           `json:"state" gorethink:"state"`
	HandledResponseTypes Arguments        `json:"handledResponseTypes" gorethink:"handledResponseTypes"`
	ResponseMode         ResponseModeType `json:"ResponseModes" gorethink:"ResponseModes"`
	DefaultResponseMode  ResponseModeType `json:"DefaultResponseMode" gorethink:"DefaultResponseMode"`

	Request
}

func NewAuthorizeRequest() *AuthorizeRequest {
	return &AuthorizeRequest{
		ResponseTypes:        Arguments{},
		HandledResponseTypes: Arguments{},
		Request:              *NewRequest(),
		ResponseMode:         ResponseModeDefault,
		// The redirect URL must be unset / nil for redirect detection to work properly:
		// RedirectURI:          &url.URL{},
	}
}

func (d *AuthorizeRequest) IsRedirectURIValid() bool {
	if d.GetRedirectURI() == nil {
		return false
	}

	raw := d.GetRedirectURI().String()
	if d.GetClient() == nil {
		return false
	}

	redirectURI, err := MatchRedirectURIWithClientRedirectURIs(raw, d.GetClient())
	if err != nil {
		return false
	}
	return IsValidRedirectURI(redirectURI)
}

func (d *AuthorizeRequest) GetResponseTypes() Arguments {
	return d.ResponseTypes
}

func (d *AuthorizeRequest) GetState() string {
	return d.State
}

func (d *AuthorizeRequest) GetRedirectURI() *url.URL {
	return d.RedirectURI
}

func (d *AuthorizeRequest) SetResponseTypeHandled(name string) {
	d.HandledResponseTypes = append(d.HandledResponseTypes, name)
}

func (d *AuthorizeRequest) DidHandleAllResponseTypes() bool {
	for _, rt := range d.ResponseTypes {
		if !d.HandledResponseTypes.Has(rt) {
			return false
		}
	}

	return len(d.ResponseTypes) > 0
}

func (d *AuthorizeRequest) GetResponseMode() ResponseModeType {
	return d.ResponseMode
}

func (d *AuthorizeRequest) SetDefaultResponseMode(defaultResponseMode ResponseModeType) {
	if d.ResponseMode == ResponseModeDefault {
		d.ResponseMode = defaultResponseMode
	}
	d.DefaultResponseMode = defaultResponseMode
}

func (d *AuthorizeRequest) GetDefaultResponseMode() ResponseModeType {
	return d.DefaultResponseMode
}
