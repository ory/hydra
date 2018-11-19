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
	"net/http"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/pkg"
)

var hasher = &fosite.BCrypt{}
var oauth2OpqaueStrategy = &compose.CommonStrategy{
	CoreStrategy:               compose.NewOAuth2HMACStrategy(fc, []byte("some super secret secret secret secret"), nil),
	OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(fc, pkg.MustINSECURELOWENTROPYRSAKEYFORTEST()),
}
var oauth2JWTStrategy = &compose.CommonStrategy{
	CoreStrategy:               compose.NewOAuth2JWTStrategy(pkg.MustINSECURELOWENTROPYRSAKEYFORTEST(), compose.NewOAuth2HMACStrategy(fc, []byte("some super secret secret secret secret"), nil)),
	OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(fc, pkg.MustINSECURELOWENTROPYRSAKEYFORTEST()),
}

var fc = &compose.Config{
	AccessTokenLifespan:        time.Second * 2,
	SendDebugMessagesToClients: true,
}

type consentMock struct {
	deny        bool
	authTime    time.Time
	requestTime time.Time
}

func (c *consentMock) HandleOAuth2AuthorizationRequest(w http.ResponseWriter, r *http.Request, req fosite.AuthorizeRequester) (*consent.HandledConsentRequest, error) {
	if c.deny {
		return nil, fosite.ErrRequestForbidden
	}

	return &consent.HandledConsentRequest{
		ConsentRequest: &consent.ConsentRequest{
			Subject:           "foo",
			SubjectIdentifier: "foo",
			ACR:               "1",
		},
		AuthenticatedAt: c.authTime,
		GrantedScope:    []string{"offline", "openid", "hydra.*"},
		Session: &consent.ConsentRequestSessionData{
			AccessToken: map[string]interface{}{},
			IDToken:     map[string]interface{}{},
		},
		RequestedAt: c.requestTime,
	}, nil
}
