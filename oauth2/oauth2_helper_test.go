/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package oauth2_test

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/x/sqlxx"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
)

var _ consent.Strategy = new(consentMock)

type consentMock struct {
	deny        bool
	authTime    time.Time
	requestTime time.Time
}

func (c *consentMock) HandleOAuth2AuthorizationRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, req fosite.AuthorizeRequester) (*consent.AcceptOAuth2ConsentRequest, error) {
	if c.deny {
		return nil, fosite.ErrRequestForbidden
	}

	return &consent.AcceptOAuth2ConsentRequest{
		ConsentRequest: &consent.OAuth2ConsentRequest{
			Subject: "foo",
			ACR:     "1",
		},
		AuthenticatedAt: sqlxx.NullTime(c.authTime),
		GrantedScope:    []string{"offline", "openid", "hydra.*"},
		Session: &consent.AcceptOAuth2ConsentRequestSession{
			AccessToken: map[string]interface{}{},
			IDToken:     map[string]interface{}{},
		},
		RequestedAt: c.requestTime,
	}, nil
}

func (c *consentMock) HandleOpenIDConnectLogout(ctx context.Context, w http.ResponseWriter, r *http.Request) (*consent.LogoutResult, error) {
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
