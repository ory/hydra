// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/trace"

	"github.com/ory/x/errorsx"
	"github.com/ory/x/otelx"

	"github.com/ory/hydra/v2/fosite/i18n"
)

// NewDeviceRequest parses an http Request returns a Device request
func (f *Fosite) NewDeviceRequest(ctx context.Context, r *http.Request) (_ DeviceRequester, err error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("github.com/ory/hydra/v2/fosite").Start(ctx, "Fosite.NewDeviceRequest")
	defer otelx.End(span, &err)

	request := NewDeviceRequest()
	request.Lang = i18n.GetLangFromRequest(f.Config.GetMessageCatalog(ctx), r)

	if r.Method != http.MethodPost {
		return request, errorsx.WithStack(ErrInvalidRequest.WithHintf("HTTP method is '%s', expected 'POST'.", r.Method))
	}
	if err := r.ParseForm(); err != nil {
		return request, errorsx.WithStack(ErrInvalidRequest.WithHint("Unable to parse HTTP body, make sure to send a properly formatted form request body.").WithWrap(err).WithDebug(err.Error()))
	}
	if len(r.PostForm) == 0 {
		return request, errorsx.WithStack(ErrInvalidRequest.WithHint("The POST body can not be empty."))
	}
	request.Form = r.PostForm

	client, clientErr := f.AuthenticateClient(ctx, r, r.PostForm)
	if clientErr != nil {
		return request, clientErr
	}
	if client.GetID() != request.Form.Get("client_id") {
		return request, errorsx.WithStack(ErrInvalidRequest.WithHint("Provided client_id mismatch."))
	}
	request.Client = client

	if !client.GetGrantTypes().Has(string(GrantTypeDeviceCode)) {
		return request, errorsx.WithStack(ErrInvalidGrant.WithHint("The requested OAuth 2.0 Client does not have the 'urn:ietf:params:oauth:grant-type:device_code' grant."))
	}

	if err := f.validateDeviceScope(ctx, r, request); err != nil {
		return request, err
	}

	if err := f.validateAudience(ctx, r, request); err != nil {
		return request, err
	}

	return request, nil
}

func (f *Fosite) validateDeviceScope(ctx context.Context, r *http.Request, request *DeviceRequest) error {
	scopes := RemoveEmpty(strings.Split(request.Form.Get("scope"), " "))
	scopeStrategy := f.Config.GetScopeStrategy(ctx)
	for _, scope := range scopes {
		if !scopeStrategy(request.Client.GetScopes(), scope) {
			return errorsx.WithStack(ErrInvalidScope.WithHintf("The OAuth 2.0 Client is not allowed to request scope '%s'.", scope))
		}
	}
	request.SetRequestedScopes(scopes)
	return nil
}
