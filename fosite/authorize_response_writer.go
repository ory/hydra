// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel/trace"

	"github.com/ory/x/errorsx"
	"github.com/ory/x/otelx"
)

func (f *Fosite) NewAuthorizeResponse(ctx context.Context, ar AuthorizeRequester, session Session) (_ AuthorizeResponder, err error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("github.com/ory/hydra/v2/fosite").Start(ctx, "Fosite.NewAuthorizeResponse")
	defer otelx.End(span, &err)

	var resp = &AuthorizeResponse{
		Header:     http.Header{},
		Parameters: url.Values{},
	}

	ctx = context.WithValue(ctx, AuthorizeRequestContextKey, ar)
	ctx = context.WithValue(ctx, AuthorizeResponseContextKey, resp)

	ar.SetSession(session)
	for _, h := range f.Config.GetAuthorizeEndpointHandlers(ctx) {
		if err := h.HandleAuthorizeEndpointRequest(ctx, ar, resp); err != nil {
			return nil, err
		}
	}

	if !ar.DidHandleAllResponseTypes() {
		return nil, errorsx.WithStack(ErrUnsupportedResponseType)
	}

	if ar.GetDefaultResponseMode() == ResponseModeFragment && ar.GetResponseMode() == ResponseModeQuery {
		return nil, ErrUnsupportedResponseMode.WithHintf("Insecure response_mode '%s' for the response_type '%s'.", ar.GetResponseMode(), ar.GetResponseTypes())
	}

	return resp, nil
}
