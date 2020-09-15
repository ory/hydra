package oauth2

import (
	"context"
	"io"

	"github.com/gobuffalo/pop/v5"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
)

type Manager interface {
	Connection(_ context.Context) *pop.Connection
	connection(ctx context.Context) *pop.Connection
	transaction(ctx context.Context, f func(context.Context, *pop.Connection) error) error
	MigrationStatus(_ context.Context, w io.Writer) error
	MigrateDown(_ context.Context, steps int) error
	MigrateUp(_ context.Context) error
	MigrateUpTo(_ context.Context, steps int) (int, error)
	PrepareMigration(_ context.Context) error
	migrateOldMigrationTables() error
	CreateConsentRequest(ctx context.Context, req *consent.ConsentRequest) error
	GetConsentRequest(ctx context.Context, challenge string) (*consent.ConsentRequest, error)
	HandleConsentRequest(ctx context.Context, challenge string, r *consent.HandledConsentRequest) (cr *consent.ConsentRequest, err error)
	RevokeSubjectConsentSession(ctx context.Context, user string) error
	RevokeSubjectClientConsentSession(ctx context.Context, user string, client string) error
	revokeConsentSession(whereStmt string, whereArgs ...interface{}) func(context.Context, *pop.Connection) error
	VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*consent.HandledConsentRequest, error)
	FindGrantedAndRememberedConsentRequests(ctx context.Context, client string, subject string) ([]consent.HandledConsentRequest, error)
	FindSubjectsGrantedConsentRequests(ctx context.Context, subject string, limit int, offset int) ([]consent.HandledConsentRequest, error)
	resolveHandledConsentRequests(requests []consent.HandledConsentRequest) ([]consent.HandledConsentRequest, error)
	CountSubjectsGrantedConsentRequests(ctx context.Context, subject string) (int, error)
	GetRememberedLoginSession(ctx context.Context, id string) (*consent.LoginSession, error)
	CreateLoginSession(ctx context.Context, session *consent.LoginSession) error
	DeleteLoginSession(ctx context.Context, id string) error
	RevokeSubjectLoginSession(ctx context.Context, subject string) error
	ConfirmLoginSession(ctx context.Context, id string, subject string, remember bool) error
	CreateLoginRequest(ctx context.Context, req *consent.LoginRequest) error
	GetLoginRequest(ctx context.Context, challenge string) (*consent.LoginRequest, error)
	HandleLoginRequest(ctx context.Context, challenge string, r *consent.HandledLoginRequest) (lr *consent.LoginRequest, err error)
	VerifyAndInvalidateLoginRequest(ctx context.Context, verifier string) (*consent.HandledLoginRequest, error)
	CreateForcedObfuscatedLoginSession(ctx context.Context, session *consent.ForcedObfuscatedLoginSession) error
	GetForcedObfuscatedLoginSession(ctx context.Context, client string, obfuscated string) (*consent.ForcedObfuscatedLoginSession, error)
	ListUserAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject string, sid string) ([]client.Client, error)
	ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject string, sid string) ([]client.Client, error)
	listUserAuthenticatedClients(ctx context.Context, subject string, sid string, channel string) ([]client.Client, error)
	CreateLogoutRequest(ctx context.Context, request *consent.LogoutRequest) error
	GetLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error)
	AcceptLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error)
	RejectLogoutRequest(ctx context.Context, challenge string) error
	VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*consent.LogoutRequest, error)
	GetClient(ctx context.Context, id string) (fosite.Client, error)
	CreateClient(ctx context.Context, c *client.Client) error
	UpdateClient(ctx context.Context, c *client.Client) error
	DeleteClient(ctx context.Context, id string) error
	GetClients(ctx context.Context, limit int, offset int) ([]client.Client, error)
	CountClients(ctx context.Context) (int, error)
	GetConcreteClient(ctx context.Context, id string) (*client.Client, error)
	Authenticate(ctx context.Context, id string, secret []byte) (*client.Client, error)
}
