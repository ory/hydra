// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"context"
	"net/http"

	"github.com/ory/fosite"
)

var _ Strategy = new(DefaultStrategy)

type Strategy interface {
	HandleOAuth2AuthorizationRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, req fosite.AuthorizeRequester) (*AcceptOAuth2ConsentRequest, error)
	HandleOpenIDConnectLogout(ctx context.Context, w http.ResponseWriter, r *http.Request) (*LogoutResult, error)
	HandleHeadlessLogout(ctx context.Context, w http.ResponseWriter, r *http.Request, sid string) error
	ObfuscateSubjectIdentifier(ctx context.Context, cl fosite.Client, subject, forcedIdentifier string) (string, error)
}
