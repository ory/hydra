package flow

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/bxcodec/faker/v3"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/consent"
	"github.com/ory/x/sqlxx"
)

func (f *Flow) setLoginRequest(r *consent.LoginRequest) {
	f.ID = r.ID
	f.RequestedScope = r.RequestedScope
	f.RequestedAudience = r.RequestedAudience
	f.Skip = r.Skip
	f.Subject = r.Subject
	f.OpenIDConnectContext = r.OpenIDConnectContext
	f.Client = r.Client
	f.ClientID = r.ClientID
	f.RequestURL = r.RequestURL
	f.SessionID = r.SessionID
	f.WasHandled = r.WasHandled
	f.ForceSubjectIdentifier = r.ForceSubjectIdentifier
	f.Verifier = r.Verifier
	f.CSRF = r.CSRF
	f.LoginAuthenticatedAt = r.AuthenticatedAt
	f.RequestedAt = r.RequestedAt
}

func (f *Flow) setHandledLoginRequest(r *consent.HandledLoginRequest) {
	f.ID = r.ID
	f.Remember = r.Remember
	f.RememberFor = r.RememberFor
	f.ACR = r.ACR
	f.AMR = r.AMR
	f.Subject = r.Subject
	f.ForceSubjectIdentifier = r.ForceSubjectIdentifier
	f.Context = r.Context
	f.WasHandled = r.WasHandled
	f.Error = r.Error
	f.RequestedAt = r.RequestedAt
	f.LoginAuthenticatedAt = r.AuthenticatedAt
}

func newTimeIterator(baseTime time.Time) func() time.Time {
	i := 0
	return func() time.Time {
		i += 1
		return baseTime.Add(time.Second * time.Duration(i))
	}
}

var nextSecond = newTimeIterator(time.Now())

func TestFlow_GetLoginRequest(t *testing.T) {
	t.Run("GetLoginRequest sets all fields on its return value", func(t *testing.T) {
		f := Flow{}
		expected := consent.LoginRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		f.setLoginRequest(&expected)
		actual := f.GetLoginRequest()
		assert.Equal(t, expected, *actual)
	})
}

func TestFlow_GetHandledLoginRequest(t *testing.T) {
	t.Run("GetHandledLoginRequest sets all fields on its return value", func(t *testing.T) {
		f := Flow{}
		expected := consent.HandledLoginRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		f.setHandledLoginRequest(&expected)
		actual := f.GetHandledLoginRequest()
		expected.LoginRequest = nil
		actual.LoginRequest = nil
		assert.Equal(t, expected, actual)
	})
}

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
	t.Run(
		"HandleLoginRequest ignores RequestedAt in its argument and copies the other fields",
		func(t *testing.T) {
			f := Flow{}
			assert.NoError(t, faker.FakeData(&f))
			f.State = FlowStateLoginInitialized

			r := consent.HandledLoginRequest{}
			assert.NoError(t, faker.FakeData(&r))
			r.ID = f.ID
			r.Subject = f.Subject
			r.ForceSubjectIdentifier = f.ForceSubjectIdentifier

			assert.NoError(t, f.HandleLoginRequest(&r))

			actual := f.GetHandledLoginRequest()
			assert.NotEqual(t, r.RequestedAt, actual.RequestedAt)
			r.LoginRequest = f.GetLoginRequest()
			actual.RequestedAt = r.RequestedAt
			assert.Equal(t, r, actual)
		},
	)
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
