package memory

import (
	"context"
	"io"
	"time"

	jose "gopkg.in/square/go-jose.v2"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"

	"github.com/gobuffalo/pop/v5"

	"github.com/ory/hydra/persistence"
)

var _ persistence.Persister = new(Persister)

type Persister struct{}

func (p *Persister) CreateConsentRequest(ctx context.Context, req *consent.ConsentRequest) error {
	panic("implement me")
}

func (p *Persister) GetConsentRequest(ctx context.Context, challenge string) (*consent.ConsentRequest, error) {
	panic("implement me")
}

func (p *Persister) HandleConsentRequest(ctx context.Context, challenge string, r *consent.HandledConsentRequest) (*consent.ConsentRequest, error) {
	panic("implement me")
}

func (p *Persister) RevokeSubjectConsentSession(ctx context.Context, user string) error {
	panic("implement me")
}

func (p *Persister) RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error {
	panic("implement me")
}

func (p *Persister) VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*consent.HandledConsentRequest, error) {
	panic("implement me")
}

func (p *Persister) FindGrantedAndRememberedConsentRequests(ctx context.Context, client, user string) ([]consent.HandledConsentRequest, error) {
	panic("implement me")
}

func (p *Persister) FindSubjectsGrantedConsentRequests(ctx context.Context, user string, limit, offset int) ([]consent.HandledConsentRequest, error) {
	panic("implement me")
}

func (p *Persister) CountSubjectsGrantedConsentRequests(ctx context.Context, user string) (int, error) {
	panic("implement me")
}

func (p *Persister) GetRememberedLoginSession(ctx context.Context, id string) (*consent.LoginSession, error) {
	panic("implement me")
}

func (p *Persister) CreateLoginSession(ctx context.Context, session *consent.LoginSession) error {
	panic("implement me")
}

func (p *Persister) DeleteLoginSession(ctx context.Context, id string) error {
	panic("implement me")
}

func (p *Persister) RevokeSubjectLoginSession(ctx context.Context, user string) error {
	panic("implement me")
}

func (p *Persister) ConfirmLoginSession(ctx context.Context, id string, subject string, remember bool) error {
	panic("implement me")
}

func (p *Persister) CreateLoginRequest(ctx context.Context, req *consent.LoginRequest) error {
	panic("implement me")
}

func (p *Persister) GetLoginRequest(ctx context.Context, challenge string) (*consent.LoginRequest, error) {
	panic("implement me")
}

func (p *Persister) HandleLoginRequest(ctx context.Context, challenge string, r *consent.HandledLoginRequest) (*consent.LoginRequest, error) {
	panic("implement me")
}

func (p *Persister) VerifyAndInvalidateLoginRequest(ctx context.Context, verifier string) (*consent.HandledLoginRequest, error) {
	panic("implement me")
}

func (p *Persister) CreateForcedObfuscatedLoginSession(ctx context.Context, session *consent.ForcedObfuscatedLoginSession) error {
	panic("implement me")
}

func (p *Persister) GetForcedObfuscatedLoginSession(ctx context.Context, client, obfuscated string) (*consent.ForcedObfuscatedLoginSession, error) {
	panic("implement me")
}

func (p *Persister) ListUserAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error) {
	panic("implement me")
}

func (p *Persister) ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject, sid string) ([]client.Client, error) {
	panic("implement me")
}

func (p *Persister) CreateLogoutRequest(ctx context.Context, request *consent.LogoutRequest) error {
	panic("implement me")
}

func (p *Persister) GetLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error) {
	panic("implement me")
}

func (p *Persister) AcceptLogoutRequest(ctx context.Context, challenge string) (*consent.LogoutRequest, error) {
	panic("implement me")
}

func (p *Persister) RejectLogoutRequest(ctx context.Context, challenge string) error {
	panic("implement me")
}

func (p *Persister) VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*consent.LogoutRequest, error) {
	panic("implement me")
}

func (p *Persister) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	panic("implement me")
}

func (p *Persister) CreateClient(ctx context.Context, c *client.Client) error {
	panic("implement me")
}

func (p *Persister) UpdateClient(ctx context.Context, c *client.Client) error {
	panic("implement me")
}

func (p *Persister) DeleteClient(ctx context.Context, id string) error {
	panic("implement me")
}

func (p *Persister) GetClients(ctx context.Context, limit, offset int) ([]client.Client, error) {
	panic("implement me")
}

func (p *Persister) CountClients(ctx context.Context) (int, error) {
	panic("implement me")
}

func (p *Persister) GetConcreteClient(ctx context.Context, id string) (*client.Client, error) {
	panic("implement me")
}

func (p *Persister) Authenticate(ctx context.Context, id string, secret []byte) (*client.Client, error) {
	panic("implement me")
}

func (p *Persister) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	panic("implement me")
}

func (p *Persister) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	panic("implement me")
}

func (p *Persister) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) (err error) {
	panic("implement me")
}

func (p *Persister) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {
	panic("implement me")
}

func (p *Persister) InvalidateAuthorizeCodeSession(ctx context.Context, code string) (err error) {
	panic("implement me")
}

func (p *Persister) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	panic("implement me")
}

func (p *Persister) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	panic("implement me")
}

func (p *Persister) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	panic("implement me")
}

func (p *Persister) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	panic("implement me")
}

func (p *Persister) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	panic("implement me")
}

func (p *Persister) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	panic("implement me")
}

func (p *Persister) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) error {
	panic("implement me")
}

func (p *Persister) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	panic("implement me")
}

func (p *Persister) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) error {
	panic("implement me")
}

func (p *Persister) GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	panic("implement me")
}

func (p *Persister) CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error {
	panic("implement me")
}

func (p *Persister) DeletePKCERequestSession(ctx context.Context, signature string) error {
	panic("implement me")
}

func (p *Persister) RevokeRefreshToken(ctx context.Context, requestID string) error {
	panic("implement me")
}

func (p *Persister) RevokeAccessToken(ctx context.Context, requestID string) error {
	panic("implement me")
}

func (p *Persister) FlushInactiveAccessTokens(ctx context.Context, notAfter time.Time) error {
	panic("implement me")
}

func (p *Persister) AddKey(ctx context.Context, set string, key *jose.JSONWebKey) error {
	panic("implement me")
}

func (p *Persister) AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) error {
	panic("implement me")
}

func (p *Persister) GetKey(ctx context.Context, set, kid string) (*jose.JSONWebKeySet, error) {
	panic("implement me")
}

func (p *Persister) GetKeySet(ctx context.Context, set string) (*jose.JSONWebKeySet, error) {
	panic("implement me")
}

func (p *Persister) DeleteKey(ctx context.Context, set, kid string) error {
	panic("implement me")
}

func (p *Persister) DeleteKeySet(ctx context.Context, set string) error {
	panic("implement me")
}

func (*Persister) MigrationStatus(_ context.Context, _ io.Writer) error {
	return nil
}

func (*Persister) MigrateDown(_ context.Context, steps int) error {
	return nil
}

func (*Persister) MigrateUp(_ context.Context) error {
	return nil
}

func (*Persister) PrepareMigration(context.Context) error {
	return nil
}

func (*Persister) Connection(_ context.Context) *pop.Connection {
	return nil
}
