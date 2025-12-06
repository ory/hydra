// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"net/http"
)

// ResponseModeHandler provides a contract for handling custom response modes
type ResponseModeHandler interface {
	// ResponseModes returns a set of supported response modes handled
	// by the interface implementation.
	//
	// In an authorize request with any of the provide response modes
	// methods `WriteAuthorizeResponse` and `WriteAuthorizeError` will be
	// invoked to write the successful or error authorization responses respectively.
	ResponseModes() ResponseModeTypes

	// WriteAuthorizeResponse writes successful responses
	//
	// Following headers are expected to be set by default:
	// header.Set("Cache-Control", "no-store")
	// header.Set("Pragma", "no-cache")
	WriteAuthorizeResponse(ctx context.Context, rw http.ResponseWriter, ar AuthorizeRequester, resp AuthorizeResponder)

	// WriteAuthorizeError writes error responses
	//
	// Following headers are expected to be set by default:
	// header.Set("Cache-Control", "no-store")
	// header.Set("Pragma", "no-cache")
	WriteAuthorizeError(ctx context.Context, rw http.ResponseWriter, ar AuthorizeRequester, err error)
}

type ResponseModeTypes []ResponseModeType

func (rs ResponseModeTypes) Has(item ResponseModeType) bool {
	for _, r := range rs {
		if r == item {
			return true
		}
	}
	return false
}

func NewDefaultResponseModeHandler() *DefaultResponseModeHandler {
	return new(DefaultResponseModeHandler)
}

type DefaultResponseModeHandler struct{}

func (d *DefaultResponseModeHandler) ResponseModes() ResponseModeTypes { return nil }
func (d *DefaultResponseModeHandler) WriteAuthorizeResponse(ctx context.Context, rw http.ResponseWriter, ar AuthorizeRequester, resp AuthorizeResponder) {
}
func (d *DefaultResponseModeHandler) WriteAuthorizeError(ctx context.Context, rw http.ResponseWriter, ar AuthorizeRequester, err error) {
}
