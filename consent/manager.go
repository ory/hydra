// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/flow"
	keysetpagination "github.com/ory/x/pagination/keysetpagination_v2"
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
		GetConsentRequest(ctx context.Context, challenge string) (*flow.OAuth2ConsentRequest, error)
		HandleConsentRequest(ctx context.Context, f *flow.Flow, r *flow.AcceptOAuth2ConsentRequest) (*flow.OAuth2ConsentRequest, error)

		RevokeSubjectConsentSession(ctx context.Context, user string) error
		RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error
		RevokeConsentSessionByID(ctx context.Context, consentRequestID string) error

		VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*flow.Flow, error)
		FindGrantedAndRememberedConsentRequest(ctx context.Context, client, user string) (*flow.Flow, error)
		FindSubjectsGrantedConsentRequests(ctx context.Context, user string, pageOpts ...keysetpagination.Option) ([]flow.Flow, *keysetpagination.Paginator, error)
		FindSubjectsSessionGrantedConsentRequests(ctx context.Context, subject, sid string, pageOpts ...keysetpagination.Option) ([]flow.Flow, *keysetpagination.Paginator, error)
		CountSubjectsGrantedConsentRequests(ctx context.Context, user string) (int, error)

		// Cookie management
		GetRememberedLoginSession(ctx context.Context, loginSessionFromCookie *flow.LoginSession, id string) (*flow.LoginSession, error)
		CreateLoginSession(ctx context.Context, session *flow.LoginSession) error
		DeleteLoginSession(ctx context.Context, id string) (deletedSession *flow.LoginSession, err error)
		RevokeSubjectLoginSession(ctx context.Context, user string) error
		ConfirmLoginSession(ctx context.Context, loginSession *flow.LoginSession) error

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

		CreateDeviceUserAuthRequest(ctx context.Context, req *flow.DeviceUserAuthRequest) (*flow.Flow, error)
		GetDeviceUserAuthRequest(ctx context.Context, challenge string) (*flow.DeviceUserAuthRequest, error)
		HandleDeviceUserAuthRequest(ctx context.Context, f *flow.Flow, challenge string, r *flow.HandledDeviceUserAuthRequest) (*flow.DeviceUserAuthRequest, error)
		VerifyAndInvalidateDeviceUserAuthRequest(ctx context.Context, verifier string) (*flow.HandledDeviceUserAuthRequest, error)

		NetworkID(ctx context.Context) uuid.UUID
	}

	ManagerProvider interface {
		ConsentManager() Manager
	}
)
