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

		// ListClientsWithLogoutURLsForSubjectAndSID returns the clients for
		// which to call front-channel and back-channel logout endpoints when
		// the subject logs out of the session with ID sid.
		//
		// One of both of these lists may be empty.
		//
		// It is no error if no authentication session can be found for the
		// given subject and sid. Both lists will be empty in that case, and a
		// nil error is returned.
		ListClientsWithLogoutURLsForSubjectAndSID(ctx context.Context, subject, sid string) (withFrontChannelURL, withBackChannelURL []client.Client, err error)
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
	// LogoutManager handles the stateless logout flow. Logout challenges and
	// verifiers are AEAD-encrypted, self-contained blobs; nothing is persisted.
	//
	// Because the flow is stateless, there are no enforceable state
	// transitions between accept and reject: a rejected challenge remains
	// technically decodable, and an already-issued verifier cannot be
	// invalidated by a later reject. This is safe because completing a logout
	// is gated on deleting the login session — the single, atomic,
	// database-backed state transition of the flow. A verifier whose session
	// is gone redirects without side effects, so replaying any of these blobs
	// cannot log a user out of a session that was not already targeted.
	LogoutManager interface {
		// CreateLogoutChallenge encodes the logout request into a stateless
		// logout challenge.
		CreateLogoutChallenge(ctx context.Context, request *flow.LogoutRequest) (challenge string, err error)

		// GetLogoutRequest decodes a logout challenge. It returns
		// ErrorLogoutFlowExpired if the embedded expiry has passed.
		GetLogoutRequest(ctx context.Context, challenge string) (*flow.LogoutRequest, error)

		// AcceptLogoutRequest exchanges a valid logout challenge for a logout
		// verifier. The verifier carries only the fields needed to complete
		// the logout. It returns ErrorLogoutFlowExpired if the embedded
		// expiry has passed.
		AcceptLogoutRequest(ctx context.Context, challenge string) (verifier string, err error)

		// RejectLogoutRequest validates the logout challenge. There is no
		// state to delete; see the interface documentation.
		RejectLogoutRequest(ctx context.Context, challenge string) error

		// VerifyAndInvalidateLogoutRequest decodes a logout verifier.
		// Invalidation is enforced by the caller deleting the login session;
		// see the interface documentation.
		VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*flow.LogoutRequest, error)
	}

	ManagerProvider interface {
		ConsentManager() Manager
	}
	ObfuscatedSubjectManagerProvider interface {
		ObfuscatedSubjectManager() ObfuscatedSubjectManager
	}
	LoginManagerProvider interface {
		LoginManager() LoginManager
	}
	LogoutManagerProvider interface {
		LogoutManager() LogoutManager
	}
)
