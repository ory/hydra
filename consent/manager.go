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
		RevokeSubjectConsentSession(ctx context.Context, subject string) error
		RevokeSubjectClientConsentSession(ctx context.Context, subject, client string) error
		RevokeConsentSessionByID(ctx context.Context, consentRequestID string) error

		CreateConsentSession(ctx context.Context, f *flow.Flow) error
		FindGrantedAndRememberedConsentRequest(ctx context.Context, client, subject string) (*flow.Flow, error)
		FindSubjectsGrantedConsentRequests(ctx context.Context, subject string, pageOpts ...keysetpagination.Option) ([]flow.Flow, *keysetpagination.Paginator, error)
		FindSubjectsSessionGrantedConsentRequests(ctx context.Context, subject, sid string, pageOpts ...keysetpagination.Option) ([]flow.Flow, *keysetpagination.Paginator, error)

		ListUserAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error)
		ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error)
	}
	ObfuscatedSubjectManager interface {
		CreateForcedObfuscatedLoginSession(ctx context.Context, session *ForcedObfuscatedLoginSession) error
		GetForcedObfuscatedLoginSession(ctx context.Context, client, obfuscated string) (*ForcedObfuscatedLoginSession, error)
	}
	LoginManager interface {
		GetRememberedLoginSession(ctx context.Context, id string) (*flow.LoginSession, error)
		DeleteLoginSession(ctx context.Context, id string) (deletedSession *flow.LoginSession, err error)
		RevokeSubjectLoginSession(ctx context.Context, subject string) error
		ConfirmLoginSession(ctx context.Context, loginSession *flow.LoginSession) error
	}
	LogoutManager interface {
		CreateLogoutRequest(ctx context.Context, request *flow.LogoutRequest) error
		GetLogoutRequest(ctx context.Context, challenge string) (*flow.LogoutRequest, error)
		AcceptLogoutRequest(ctx context.Context, challenge string) (*flow.LogoutRequest, error)
		RejectLogoutRequest(ctx context.Context, challenge string) error
		VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*flow.LogoutRequest, error)
	}

	ManagerProvider                  interface{ ConsentManager() Manager }
	ObfuscatedSubjectManagerProvider interface {
		ObfuscatedSubjectManager() ObfuscatedSubjectManager
	}
	LoginManagerProvider  interface{ LoginManager() LoginManager }
	LogoutManagerProvider interface{ LogoutManager() LogoutManager }
)
