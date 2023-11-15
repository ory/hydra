// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/flow"
)

type ForcedObfuscatedLoginSession struct {
	ClientID          string    `db:"client_id"`
	Subject           string    `db:"subject"`
	SubjectObfuscated string    `db:"subject_obfuscated"`
	NID               uuid.UUID `db:"nid"`
}

func (ForcedObfuscatedLoginSession) TableName() string {
	return "hydra_oauth2_obfuscated_authentication_session"
}

type (
	Manager interface {
		CreateConsentRequest(ctx context.Context, f *flow.Flow, req *flow.OAuth2ConsentRequest) error
		GetConsentRequest(ctx context.Context, challenge string) (*flow.OAuth2ConsentRequest, error)
		HandleConsentRequest(ctx context.Context, f *flow.Flow, r *flow.AcceptOAuth2ConsentRequest) (*flow.OAuth2ConsentRequest, error)
		RevokeSubjectConsentSession(ctx context.Context, user string) error
		RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error

		VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*flow.AcceptOAuth2ConsentRequest, error)
		FindGrantedAndRememberedConsentRequests(ctx context.Context, client, user string) ([]flow.AcceptOAuth2ConsentRequest, error)
		FindSubjectsGrantedConsentRequests(ctx context.Context, user string, limit, offset int) ([]flow.AcceptOAuth2ConsentRequest, error)
		FindSubjectsSessionGrantedConsentRequests(ctx context.Context, user, sid string, limit, offset int) ([]flow.AcceptOAuth2ConsentRequest, error)
		CountSubjectsGrantedConsentRequests(ctx context.Context, user string) (int, error)

		// Cookie management
		GetRememberedLoginSession(ctx context.Context, loginSessionFromCookie *flow.LoginSession, id string) (*flow.LoginSession, error)
		CreateLoginSession(ctx context.Context, session *flow.LoginSession) error
		DeleteLoginSession(ctx context.Context, id string) (deletedSession *flow.LoginSession, err error)
		RevokeSubjectLoginSession(ctx context.Context, user string) error
		ConfirmLoginSession(ctx context.Context, loginSession *flow.LoginSession) error

		CreateLoginRequest(ctx context.Context, req *flow.LoginRequest) (*flow.Flow, error)
		GetLoginRequest(ctx context.Context, challenge string) (*flow.LoginRequest, error)
		HandleLoginRequest(ctx context.Context, f *flow.Flow, challenge string, r *flow.HandledLoginRequest) (*flow.LoginRequest, error)
		VerifyAndInvalidateLoginRequest(ctx context.Context, verifier string) (*flow.HandledLoginRequest, error)

		CreateForcedObfuscatedLoginSession(ctx context.Context, session *ForcedObfuscatedLoginSession) error
		GetForcedObfuscatedLoginSession(ctx context.Context, client, obfuscated string) (*ForcedObfuscatedLoginSession, error)

		ListUserAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error)
		ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error)

		CreateLogoutRequest(ctx context.Context, request *flow.LogoutRequest) error
		GetLogoutRequest(ctx context.Context, challenge string) (*flow.LogoutRequest, error)
		AcceptLogoutRequest(ctx context.Context, challenge string) (*flow.LogoutRequest, error)
		RejectLogoutRequest(ctx context.Context, challenge string) error
		VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*flow.LogoutRequest, error)
	}

	ManagerProvider interface {
		ConsentManager() Manager
	}
)
