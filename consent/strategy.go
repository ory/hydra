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
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package consent

import (
	"context"
	"net/http"

	"github.com/ory/fosite"
)

var _ Strategy = new(DefaultStrategy)

type Strategy interface {
	HandleOAuth2AuthorizationRequest(w http.ResponseWriter, r *http.Request, req fosite.AuthorizeRequester) (*HandledConsentRequest, error)
	HandleOpenIDConnectLogout(w http.ResponseWriter, r *http.Request) (*LogoutResult, error)
	ExecuteBackChannelLogoutBySubject(ctx context.Context, r *http.Request, subject string) error
	ExecuteBackChannelLogoutBySession(ctx context.Context, r *http.Request, subject, sid string) error
	ExecuteBackChannelLogoutByClient(ctx context.Context, r *http.Request, subject, client string) error
	ExecuteBackChannelLogoutByClientSession(ctx context.Context, r *http.Request, subject, client, sid string) error
}
