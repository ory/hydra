package flow

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/x/sqlxx"
)

func newTimeIterator(baseTime time.Time) func() time.Time {
	i := 0
	return func() time.Time {
		i += 1
		return baseTime.Add(time.Second * time.Duration(i))
	}
}

var nextSecond = newTimeIterator(time.Now())

func TestFlow_NewFlow(t *testing.T) {
	clientID := uuid.Must(uuid.NewV4())
	nonDefaultLoginRequest := &consent.LoginRequest{
		ID:                     "t1",
		RequestedScope:         sqlxx.StringSlicePipeDelimiter{"t1-requested_scope"},
		RequestedAudience:      sqlxx.StringSlicePipeDelimiter{"t1-requested_audience"},
		Skip:                   true,
		Subject:                "t1-subject",
		OpenIDConnectContext:   &consent.OpenIDConnectContext{Display: "t1-display"},
		RequestURL:             "http://request/t1",
		SessionID:              sqlxx.NullString("t1-auth_session"),
		Verifier:               "t1-verifier",
		CSRF:                   "t1-csrf",
		WasHandled:             true,
		Client:                 &client.Client{ID: clientID},
		ClientID:               clientID.String(),
		ForceSubjectIdentifier: "t1-force-subject-id",
		AuthenticatedAt:        sqlxx.NullTime(nextSecond()),
		RequestedAt:            nextSecond(),
	}

	roundtrip := NewFlow(nonDefaultLoginRequest).GetLoginRequest()
	assert.Equal(t, nonDefaultLoginRequest, roundtrip)
}

func TestFlow_HandleLoginRequest(t *testing.T) {
	t2ID := uuid.Must(uuid.NewV4())
	tests := []struct {
		name    string
		flow    Flow
		arg     *consent.HandledLoginRequest
		wantErr bool
	}{
		{
			name: "HandleLoginRequest ignores RequestedAt in its argument and copies the other fields",
			flow: Flow{
				ID:                     "t2-id",
				RequestedScope:         sqlxx.StringSlicePipeDelimiter{"t2-requested_scope"},
				RequestedAudience:      sqlxx.StringSlicePipeDelimiter{"t2-requested_audience"},
				Skip:                   false,
				Subject:                "t2-sub-invalid",
				OpenIDConnectContext:   &consent.OpenIDConnectContext{Display: "t2-display"},
				Client:                 &client.Client{ID: t2ID},
				ClientID:               t2ID.String(),
				RequestURL:             "http://test/t2",
				SessionID:              sqlxx.NullString("t2-auth_session"),
				Verifier:               "t2-verifier",
				CSRF:                   "t2-csrf",
				LoginInitializedAt:     sqlxx.NullTime(nextSecond()),
				RequestedAt:            nextSecond(),
				State:                  FlowStateLoginInitialized,
				Remember:               false,
				RememberFor:            0,
				ACR:                    "",
				AMR:                    []string{},
				ForceSubjectIdentifier: "t2-force-sub-id-invalid",
				Context:                []byte("{\"v\": \"t2-invlid\"}"),
				WasHandled:             false,
				Error:                  nil,
				LoginAuthenticatedAt:   sqlxx.NullTime(nextSecond()),
			},
			arg: &consent.HandledLoginRequest{
				ID:                     "t2-id",
				Subject:                "t2-sub-valid",
				ForceSubjectIdentifier: "t2-force-sub-id-valid",
				WasHandled:             true,
				Remember:               true,
				RememberFor:            1,
				ACR:                    "valid",
				AMR:                    []string{"t2-amr-1-valid", "t2-amr-2-valid"},
				Context:                []byte("{\"v\": \"t2-valid\"}"),
				LoginRequest:           nil,
				Error:                  &consent.RequestDeniedError{Description: "expected-error"},
				RequestedAt:            nextSecond(),
				AuthenticatedAt:        sqlxx.NullTime(nextSecond()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := &consent.HandledLoginRequest{
				ID:                     tt.flow.ID,
				Remember:               tt.arg.Remember,
				RememberFor:            tt.arg.RememberFor,
				ACR:                    tt.arg.ACR,
				AMR:                    tt.arg.AMR,
				Subject:                tt.arg.Subject,
				ForceSubjectIdentifier: tt.arg.ForceSubjectIdentifier,
				Context:                tt.arg.Context,
				WasHandled:             true,
				LoginRequest:           &consent.LoginRequest{ID: "invalid-set-later"},
				Error:                  tt.arg.Error,
				RequestedAt:            tt.flow.RequestedAt,
				AuthenticatedAt:        tt.arg.AuthenticatedAt,
			}
			if err := tt.flow.HandleLoginRequest(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("Flow.HandleLoginRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
			expected.LoginRequest = tt.flow.GetLoginRequest()
			actual := tt.flow.GetHandledLoginRequest()
			assert.Equal(t, expected, &actual)
		})
	}
}

func TestFlow_InitializeConsent(t *testing.T) {
	t.Run("InitializeConsent transitions the flow into FlowStateConsentInitialized", func(t *testing.T) {
		f := NewFlow(&consent.LoginRequest{
			ID:         "t3-id",
			Subject:    "t3-sub",
			WasHandled: false,
		})
		assert.NoError(t, f.HandleLoginRequest(&consent.HandledLoginRequest{
			ID:         "t3-id",
			Subject:    "t3-sub",
			WasHandled: false,
		}))
		assert.NoError(t, f.InitializeConsent())
		assert.Equal(t, FlowStateConsentInitialized, f.State)
		assert.Equal(t, true, f.WasHandled)
	})
	t.Run("InitializeConsent fails when flow.WasHandled is true", func(t *testing.T) {
		f := NewFlow(&consent.LoginRequest{
			ID:         "t3-id",
			Subject:    "t3-sub",
			WasHandled: false,
		})
		assert.NoError(t, f.HandleLoginRequest(&consent.HandledLoginRequest{
			ID:         "t3-id",
			Subject:    "t3-sub",
			WasHandled: true,
		}))
		err := f.InitializeConsent()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "verifier has already been used")
	})
}
