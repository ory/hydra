// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/x/sqlxx"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/consent"
)

var _ consent.Strategy = new(consentMock)

type consentMock struct {
	deny        bool
	authTime    time.Time
	requestTime time.Time
}

func (c *consentMock) HandleOAuth2AuthorizationRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, req fosite.AuthorizeRequester) (*flow.AcceptOAuth2ConsentRequest, *flow.Flow, error) {
	if c.deny {
		return nil, nil, fosite.ErrRequestForbidden
	}

	return &flow.AcceptOAuth2ConsentRequest{
		ConsentRequest: &flow.OAuth2ConsentRequest{
			Subject: "foo",
			ACR:     "1",
		},
		AuthenticatedAt: sqlxx.NullTime(c.authTime),
		GrantedScope:    []string{"offline", "openid", "hydra.*"},
		Session: &flow.AcceptOAuth2ConsentRequestSession{
			AccessToken: map[string]interface{}{},
			IDToken:     map[string]interface{}{},
		},
		RequestedAt: c.requestTime,
	}, nil, nil
}

func (c *consentMock) HandleOpenIDConnectLogout(ctx context.Context, w http.ResponseWriter, r *http.Request) (*flow.LogoutResult, error) {
	panic("not implemented")
}

func (c *consentMock) HandleHeadlessLogout(ctx context.Context, w http.ResponseWriter, r *http.Request, sid string) error {
	panic("not implemented")
}

func (c *consentMock) ObfuscateSubjectIdentifier(ctx context.Context, cl fosite.Client, subject, forcedIdentifier string) (string, error) {
	if c, ok := cl.(*client.Client); ok && c.SubjectType == "pairwise" {
		panic("not implemented")
	} else if !ok {
		return "", errors.New("Unable to type assert OAuth 2.0 Client to *client.Client")
	}
	return subject, nil
}
