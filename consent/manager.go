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

type ForcedObfuscatedAuthenticationSession struct {
	ClientID          string `db:"client_id"`
	Subject           string `db:"subject"`
	SubjectObfuscated string `db:"subject_obfuscated"`
}

type Manager interface {
	CreateConsentRequest(*ConsentRequest) error
	GetConsentRequest(challenge string) (*ConsentRequest, error)
	HandleConsentRequest(challenge string, r *HandledConsentRequest) (*ConsentRequest, error)
	RevokeUserConsentSession(user string) error
	RevokeUserClientConsentSession(user, client string) error

	VerifyAndInvalidateConsentRequest(verifier string) (*HandledConsentRequest, error)
	FindPreviouslyGrantedConsentRequests(client string, user string) ([]HandledConsentRequest, error)
	FindPreviouslyGrantedConsentRequestsByUser(user string, limit, offset int) ([]HandledConsentRequest, error)

	// Cookie management
	GetAuthenticationSession(id string) (*AuthenticationSession, error)
	CreateAuthenticationSession(*AuthenticationSession) error
	DeleteAuthenticationSession(id string) error
	RevokeUserAuthenticationSession(user string) error

	CreateAuthenticationRequest(*AuthenticationRequest) error
	GetAuthenticationRequest(challenge string) (*AuthenticationRequest, error)
	HandleAuthenticationRequest(challenge string, r *HandledAuthenticationRequest) (*AuthenticationRequest, error)
	VerifyAndInvalidateAuthenticationRequest(verifier string) (*HandledAuthenticationRequest, error)

	CreateForcedObfuscatedAuthenticationSession(*ForcedObfuscatedAuthenticationSession) error
	GetForcedObfuscatedAuthenticationSession(client, obfuscated string) (*ForcedObfuscatedAuthenticationSession, error)
}
