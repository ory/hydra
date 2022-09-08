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
	"time"

	"github.com/gofrs/uuid"

	"github.com/ory/hydra/client"
)

type ForcedObfuscatedLoginSession struct {
	ClientID          string    `db:"client_id"`
	Subject           string    `db:"subject"`
	SubjectObfuscated string    `db:"subject_obfuscated"`
	NID               uuid.UUID `db:"nid"`
}

func (_ ForcedObfuscatedLoginSession) TableName() string {
	return "hydra_oauth2_obfuscated_authentication_session"
}

type Manager interface {
	CreateConsentRequest(ctx context.Context, req *OAuth2ConsentRequest) error
	GetConsentRequest(ctx context.Context, challenge string) (*OAuth2ConsentRequest, error)
	HandleConsentRequest(ctx context.Context, r *AcceptOAuth2ConsentRequest) (*OAuth2ConsentRequest, error)
	RevokeSubjectConsentSession(ctx context.Context, user string) error
	RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error

	VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*AcceptOAuth2ConsentRequest, error)
	FindGrantedAndRememberedConsentRequests(ctx context.Context, client, user string) ([]AcceptOAuth2ConsentRequest, error)
	FindSubjectsGrantedConsentRequests(ctx context.Context, user string, limit, offset int) ([]AcceptOAuth2ConsentRequest, error)
	CountSubjectsGrantedConsentRequests(ctx context.Context, user string) (int, error)

	// Cookie management
	GetRememberedLoginSession(ctx context.Context, id string) (*LoginSession, error)
	CreateLoginSession(ctx context.Context, session *LoginSession) error
	DeleteLoginSession(ctx context.Context, id string) error
	RevokeSubjectLoginSession(ctx context.Context, user string) error
	ConfirmLoginSession(ctx context.Context, id string, authTime time.Time, subject string, remember bool) error

	CreateLoginRequest(ctx context.Context, req *LoginRequest) error
	GetLoginRequest(ctx context.Context, challenge string) (*LoginRequest, error)
	HandleLoginRequest(ctx context.Context, challenge string, r *HandledLoginRequest) (*LoginRequest, error)
	VerifyAndInvalidateLoginRequest(ctx context.Context, verifier string) (*HandledLoginRequest, error)

	CreateForcedObfuscatedLoginSession(ctx context.Context, session *ForcedObfuscatedLoginSession) error
	GetForcedObfuscatedLoginSession(ctx context.Context, client, obfuscated string) (*ForcedObfuscatedLoginSession, error)

	ListUserAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error)
	ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error)

	CreateLogoutRequest(ctx context.Context, request *LogoutRequest) error
	GetLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error)
	AcceptLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error)
	RejectLogoutRequest(ctx context.Context, challenge string) error
	VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*LogoutRequest, error)
}
