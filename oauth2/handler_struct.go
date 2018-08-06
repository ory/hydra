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

package oauth2

import (
	"net/url"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	OAuth2  fosite.OAuth2Provider
	Consent consent.Strategy
	Storage pkg.FositeStorer

	H herodot.Writer

	ForcedHTTP bool
	ErrorURL   url.URL

	AccessTokenLifespan time.Duration
	IDTokenLifespan     time.Duration
	CookieStore         sessions.Store

	OpenIDJWTStrategy      *jwk.RS256JWTStrategy
	AccessTokenJWTStrategy *jwk.RS256JWTStrategy

	L logrus.FieldLogger

	ScopeStrategy fosite.ScopeStrategy

	IssuerURL string

	ClaimsSupported  string
	ScopesSupported  string
	SubjectTypes     []string
	UserinfoEndpoint string
}
