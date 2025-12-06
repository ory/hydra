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

	"github.com/pkg/errors"
)

type TokenIntrospector interface {
	IntrospectToken(ctx context.Context, token string, tokenUse TokenUse, accessRequest AccessRequester, scopes []string) (TokenUse, error)
}

func AccessTokenFromRequest(req *http.Request) string {
	// According to https://tools.ietf.org/html/rfc6750 you can pass tokens through:
	// - Form-Encoded Body Parameter. Recommended, more likely to appear. e.g.: Authorization: Bearer mytoken123
	// - URI Query Parameter e.g. access_token=mytoken123

	auth := req.Header.Get("Authorization")
	split := strings.SplitN(auth, " ", 2)
	if len(split) != 2 || !strings.EqualFold(split[0], "bearer") {
		// Nothing in Authorization header, try access_token
		// Empty string returned if there's no such parameter
		if err := req.ParseMultipartForm(1 << 20); err != nil && err != http.ErrNotMultipart {
			return ""
		}
		return req.Form.Get("access_token")
	}

	return split[1]
}

func (f *Fosite) IntrospectToken(ctx context.Context, token string, tokenUse TokenUse, session Session, scopes ...string) (_ TokenUse, _ AccessRequester, err error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("github.com/ory/hydra/v2/fosite").Start(ctx, "Fosite.IntrospectToken")
	defer otelx.End(span, &err)

	var found = false
	var foundTokenUse TokenUse = ""

	ar := NewAccessRequest(session)
	for _, validator := range f.Config.GetTokenIntrospectionHandlers(ctx) {
		tu, err := validator.IntrospectToken(ctx, token, tokenUse, ar, scopes)
		if err == nil {
			found = true
			foundTokenUse = tu
		} else if errors.Is(err, ErrUnknownRequest) {
			// do nothing
		} else {
			rfcerr := ErrorToRFC6749Error(err)
			return "", nil, errorsx.WithStack(rfcerr)
		}
	}

	if !found {
		return "", nil, errorsx.WithStack(ErrRequestUnauthorized.WithHint("Unable to find a suitable validation strategy for the token, thus it is invalid."))
	}

	return foundTokenUse, ar, nil
}
