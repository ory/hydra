// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/trace"

	"github.com/ory/hydra/v2/fosite/i18n"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/otelx"

	"github.com/pkg/errors"
)

// Implements
//   - https://tools.ietf.org/html/rfc6749#section-2.3.1
//     Clients in possession of a client password MAY use the HTTP Basic
//     authentication scheme as defined in [RFC2617] to authenticate with
//     the authorization server.  The client identifier is encoded using the
//     "application/x-www-form-urlencoded" encoding algorithm per
//     Appendix B, and the encoded value is used as the username; the client
//     password is encoded using the same algorithm and used as the
//     password.  The authorization server MUST support the HTTP Basic
//     authentication scheme for authenticating clients that were issued a
//     client password.
//     Including the client credentials in the request-body using the two
//     parameters is NOT RECOMMENDED and SHOULD be limited to clients unable
//     to directly utilize the HTTP Basic authentication scheme (or other
//     password-based HTTP authentication schemes).  The parameters can only
//     be transmitted in the request-body and MUST NOT be included in the
//     request URI.
//   - https://tools.ietf.org/html/rfc6749#section-3.2.1
//   - Confidential clients or other clients issued client credentials MUST
//     authenticate with the authorization server as described in
//     Section 2.3 when making requests to the token endpoint.
//   - If the client type is confidential or the client was issued client
//     credentials (or assigned other authentication requirements), the
//     client MUST authenticate with the authorization server as described
//     in Section 3.2.1.
func (f *Fosite) NewAccessRequest(ctx context.Context, r *http.Request, session Session) (_ AccessRequester, err error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("github.com/ory/hydra/v2/fosite").Start(ctx, "Fosite.NewAccessRequest")
	defer otelx.End(span, &err)

	accessRequest := NewAccessRequest(session)
	accessRequest.Request.Lang = i18n.GetLangFromRequest(f.Config.GetMessageCatalog(ctx), r)

	ctx = context.WithValue(ctx, RequestContextKey, r)
	ctx = context.WithValue(ctx, AccessRequestContextKey, accessRequest)

	if r.Method != "POST" {
		return accessRequest, errorsx.WithStack(ErrInvalidRequest.WithHintf("HTTP method is '%s', expected 'POST'.", r.Method))
	} else if err := r.ParseMultipartForm(1 << 20); err != nil && err != http.ErrNotMultipart {
		return accessRequest, errorsx.WithStack(ErrInvalidRequest.WithHint("Unable to parse HTTP body, make sure to send a properly formatted form request body.").WithWrap(err).WithDebug(err.Error()))
	} else if len(r.PostForm) == 0 {
		return accessRequest, errorsx.WithStack(ErrInvalidRequest.WithHint("The POST body can not be empty."))
	}

	accessRequest.Form = r.PostForm
	if session == nil {
		return accessRequest, errors.New("Session must not be nil")
	}

	accessRequest.SetRequestedScopes(RemoveEmpty(strings.Split(r.PostForm.Get("scope"), " ")))
	accessRequest.SetRequestedAudience(GetAudiences(r.PostForm))
	accessRequest.GrantTypes = RemoveEmpty(strings.Split(r.PostForm.Get("grant_type"), " "))
	if len(accessRequest.GrantTypes) < 1 {
		return accessRequest, errorsx.WithStack(ErrInvalidRequest.WithHint("Request parameter 'grant_type' is missing"))
	}

	client, clientErr := f.AuthenticateClient(ctx, r, r.PostForm)
	if clientErr == nil {
		accessRequest.Client = client
	}

	var found = false
	for _, loader := range f.Config.GetTokenEndpointHandlers(ctx) {
		// Is the loader responsible for handling the request?
		if !loader.CanHandleTokenEndpointRequest(ctx, accessRequest) {
			continue
		}

		// The handler **is** responsible!

		// Is the client supplied in the request? If not can this handler skip client auth?
		if !loader.CanSkipClientAuth(ctx, accessRequest) && clientErr != nil {
			// No client and handler can not skip client auth -> error.
			return accessRequest, clientErr
		}

		// All good.
		if err := loader.HandleTokenEndpointRequest(ctx, accessRequest); err == nil {
			found = true
		} else if errors.Is(err, ErrUnknownRequest) {
			// This is a duplicate because it should already have been handled by
			// `loader.CanHandleTokenEndpointRequest(accessRequest)` but let's keep it for sanity.
			//
			continue
		} else if err != nil {
			return accessRequest, err
		}
	}

	if !found {
		return nil, errorsx.WithStack(ErrInvalidRequest)
	}
	return accessRequest, nil
}
