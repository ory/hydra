// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/trace"

	"github.com/ory/x/errorsx"
	"github.com/ory/x/otelx"
)

// NewPushedAuthorizeResponse executes the handlers and builds the response
func (f *Fosite) NewPushedAuthorizeResponse(ctx context.Context, ar AuthorizeRequester, session Session) (_ PushedAuthorizeResponder, err error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("github.com/ory/hydra/v2/fosite").Start(ctx, "Fosite.NewPushedAuthorizeResponse")
	defer otelx.End(span, &err)

	// Get handlers. If no handlers are defined, this is considered a misconfigured Fosite instance.
	handlersProvider, ok := f.Config.(PushedAuthorizeRequestHandlersProvider)
	if !ok {
		return nil, errorsx.WithStack(ErrServerError.WithHint(ErrorPARNotSupported).WithDebug(DebugPARRequestsHandlerMissing))
	}

	var resp = &PushedAuthorizeResponse{
		Header: http.Header{},
		Extra:  map[string]interface{}{},
	}

	ctx = context.WithValue(ctx, AuthorizeRequestContextKey, ar)
	ctx = context.WithValue(ctx, PushedAuthorizeResponseContextKey, resp)

	ar.SetSession(session)
	for _, h := range handlersProvider.GetPushedAuthorizeEndpointHandlers(ctx) {
		if err := h.HandlePushedAuthorizeEndpointRequest(ctx, ar, resp); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

// WritePushedAuthorizeResponse writes the PAR response
func (f *Fosite) WritePushedAuthorizeResponse(ctx context.Context, rw http.ResponseWriter, ar AuthorizeRequester, resp PushedAuthorizeResponder) {
	// Set custom headers, e.g. "X-MySuperCoolCustomHeader" or "X-DONT-CACHE-ME"...
	wh := rw.Header()
	rh := resp.GetHeader()
	for k := range rh {
		wh.Set(k, rh.Get(k))
	}

	wh.Set("Cache-Control", "no-store")
	wh.Set("Pragma", "no-cache")
	wh.Set("Content-Type", "application/json;charset=UTF-8")

	js, err := json.Marshal(resp.ToMap())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")

	rw.WriteHeader(http.StatusCreated)
	_, _ = rw.Write(js)
}

// WritePushedAuthorizeError writes the PAR error
func (f *Fosite) WritePushedAuthorizeError(ctx context.Context, rw http.ResponseWriter, ar AuthorizeRequester, err error) {
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")
	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")

	sendDebugMessagesToClient := f.Config.GetSendDebugMessagesToClients(ctx)
	rfcerr := ErrorToRFC6749Error(err).WithLegacyFormat(f.Config.GetUseLegacyErrorFormat(ctx)).
		WithExposeDebug(sendDebugMessagesToClient).WithLocalizer(f.Config.GetMessageCatalog(ctx), getLangFromRequester(ar))

	js, err := json.Marshal(rfcerr)
	if err != nil {
		if sendDebugMessagesToClient {
			errorMessage := EscapeJSONString(err.Error())
			http.Error(rw, fmt.Sprintf(`{"error":"server_error","error_description":"%s"}`, errorMessage), http.StatusInternalServerError)
		} else {
			http.Error(rw, `{"error":"server_error"}`, http.StatusInternalServerError)
		}
		return
	}

	rw.WriteHeader(rfcerr.CodeField)
	_, _ = rw.Write(js)
}
