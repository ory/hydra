// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/ory/hydra/v2/client"
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
	RevokeLoginSessionConsentSession(ctx context.Context, loginSessionId string) error
	RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error
	RevokeSubjectClientLoginSessionConsentSession(ctx context.Context, user, client, loginSessionId string) error

	VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*AcceptOAuth2ConsentRequest, error)
	FindGrantedAndRememberedConsentRequests(ctx context.Context, client, user string) ([]AcceptOAuth2ConsentRequest, error)
	FindSubjectsGrantedConsentRequests(ctx context.Context, user string, limit, offset int) ([]AcceptOAuth2ConsentRequest, error)
	FindSubjectsSessionGrantedConsentRequests(ctx context.Context, user, sid string, limit, offset int) ([]AcceptOAuth2ConsentRequest, error)
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

	ListUserSessionAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject, sid string) ([]client.LoginSessionClient, error)
	ListUserSessionAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject, sid string) ([]client.LoginSessionClient, error)
	ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject string) ([]client.LoginSessionClient, error)

	CreateLogoutRequest(ctx context.Context, request *LogoutRequest) error
	GetLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error)
	AcceptLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error)
	RejectLogoutRequest(ctx context.Context, challenge string) error
	VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*LogoutRequest, error)
}
