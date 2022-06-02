package testhelpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ory/fosite"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/hydra/driver"
	"github.com/ory/x/sqlxx"
)

type JanitorSessionTestHelper struct {
	uniqueName string
	reg        driver.Registry
}

func NewJanitorSessionTestHelper(reg driver.Registry, uniqueName string) *JanitorSessionTestHelper {
	return &JanitorSessionTestHelper{reg: reg, uniqueName: uniqueName}
}

func (h *JanitorSessionTestHelper) CreateEnvironmentForSession(t *testing.T, ctx context.Context, id string) {
	clientDTO := h.CreateClient(t, ctx)
	h.CreateLoginRequest(t, ctx, clientDTO, id)
	h.CreateConsentRequest(t, ctx, id)
}

func (h *JanitorSessionTestHelper) ValidateSessionExist(t *testing.T, ctx context.Context, id string) {
	session, err := h.reg.ConsentManager().GetRememberedLoginSession(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, session)
	require.NotZero(t, session.ID)
}

func (h *JanitorSessionTestHelper) ValidateSessionNotExist(t *testing.T, ctx context.Context, id string) {
	session, err := h.reg.ConsentManager().GetRememberedLoginSession(ctx, id)
	require.Error(t, err)
	rpcErr := fosite.ErrorToRFC6749Error(err)
	require.Equal(t, fosite.ErrNotFound.StatusCode(), rpcErr.StatusCode())
	require.Nil(t, session)
}

func (h *JanitorSessionTestHelper) CreateClient(t *testing.T, ctx context.Context) *client.Client {
	clientReq := &client.Client{
		OutfacingID:  fmt.Sprintf("%s_flush-login-consent-1", h.uniqueName),
		RedirectURIs: []string{"http://redirect"},
	}
	require.NoError(t, h.reg.ClientManager().CreateClient(ctx, clientReq))
	return clientReq
}

func (h *JanitorSessionTestHelper) CreateLoginSession(t *testing.T, ctx context.Context, sessionID string) {
	require.NoError(t, h.reg.ConsentManager().CreateLoginSession(ctx, &consent.LoginSession{
		ID:       sessionID,
		Subject:  "sub",
		Remember: true,
	}))
}

func (h *JanitorSessionTestHelper) CreateLoginRequest(t *testing.T, ctx context.Context, clientDTO *client.Client, sessionID string) {
	require.NoError(t, h.reg.ConsentManager().CreateLoginRequest(ctx, &consent.LoginRequest{
		ID:              fmt.Sprintf("%s_flush-login-1", h.uniqueName),
		RequestedScope:  []string{"foo", "bar"},
		Subject:         fmt.Sprintf("%s_flush-login-1", h.uniqueName),
		Client:          clientDTO,
		RequestURL:      "http://redirect",
		RequestedAt:     time.Now().Round(time.Second),
		AuthenticatedAt: sqlxx.NullTime(time.Now().Round(time.Second)),
		Verifier:        fmt.Sprintf("%s_flush-login-1", h.uniqueName),
		SessionID:       sqlxx.NullString(sessionID),
	}))
}

func (h *JanitorSessionTestHelper) CreateConsentRequest(t *testing.T, ctx context.Context, sessionID string) {
	require.NoError(t, h.reg.ConsentManager().CreateConsentRequest(ctx, &consent.ConsentRequest{
		ID:                   fmt.Sprintf("%s_flush-consent-1", h.uniqueName),
		RequestedScope:       []string{"foo", "bar"},
		Subject:              fmt.Sprintf("%s_flush-consent-1", h.uniqueName),
		OpenIDConnectContext: nil,
		ClientID:             fmt.Sprintf("%s_flush-login-consent-1", h.uniqueName),
		RequestURL:           "http://redirect",
		LoginChallenge:       sqlxx.NullString(fmt.Sprintf("%s_flush-login-1", h.uniqueName)),
		RequestedAt:          time.Now().Round(time.Second),
		Verifier:             fmt.Sprintf("%s_flush-consent-1", h.uniqueName),
		LoginSessionID:       sqlxx.NullString(sessionID),
	}))
}
