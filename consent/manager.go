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

import "context"

type ForcedObfuscatedAuthenticationSession struct {
	ClientID          string `db:"client_id"`
	Subject           string `db:"subject"`
	SubjectObfuscated string `db:"subject_obfuscated"`
}

type Manager interface {
	CreateConsentRequest(ctx context.Context, req *ConsentRequest) error
	GetConsentRequest(ctx context.Context, challenge string) (*ConsentRequest, error)
	HandleConsentRequest(ctx context.Context, challenge string, r *HandledConsentRequest) (*ConsentRequest, error)
	RevokeUserConsentSession(ctx context.Context, user string) error
	RevokeUserClientConsentSession(ctx context.Context, user, client string) error

	VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*HandledConsentRequest, error)
	FindPreviouslyGrantedConsentRequests(ctx context.Context, client, user string) ([]HandledConsentRequest, error)
	FindPreviouslyGrantedConsentRequestsByUser(ctx context.Context, user string, limit, offset int) ([]HandledConsentRequest, error)

	// Cookie management
	GetAuthenticationSession(ctx context.Context, id string) (*AuthenticationSession, error)
	CreateAuthenticationSession(ctx context.Context, session *AuthenticationSession) error
	DeleteAuthenticationSession(ctx context.Context, id string) error
	RevokeUserAuthenticationSession(ctx context.Context, user string) error

	CreateAuthenticationRequest(ctx context.Context, req *AuthenticationRequest) error
	GetAuthenticationRequest(ctx context.Context, challenge string) (*AuthenticationRequest, error)
	HandleAuthenticationRequest(ctx context.Context, challenge string, r *HandledAuthenticationRequest) (*AuthenticationRequest, error)
	VerifyAndInvalidateAuthenticationRequest(ctx context.Context, verifier string) (*HandledAuthenticationRequest, error)

	CreateForcedObfuscatedAuthenticationSession(ctx context.Context, session *ForcedObfuscatedAuthenticationSession) error
	GetForcedObfuscatedAuthenticationSession(ctx context.Context, client, obfuscated string) (*ForcedObfuscatedAuthenticationSession, error)
}
