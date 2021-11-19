package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bxcodec/faker/v3"

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

func (f *Flow) setConsentRequest(r *consent.ConsentRequest) {
	f.ConsentChallengeID = sqlxx.NullString(r.ID)
	f.RequestedScope = r.RequestedScope
	f.RequestedAudience = r.RequestedAudience
	f.ConsentSkip = r.Skip
	f.Subject = r.Subject
	f.OpenIDConnectContext = r.OpenIDConnectContext
	f.Client = r.Client
	f.ClientID = r.ClientID
	f.RequestURL = r.RequestURL
	f.ID = r.LoginChallenge.String()
	f.SessionID = r.LoginSessionID
	f.ACR = r.ACR
	f.AMR = r.AMR
	f.Context = r.Context
	f.CHWasHandled = r.WasHandled
	f.ForceSubjectIdentifier = r.ForceSubjectIdentifier
	f.ConsentVerifier = r.Verifier
	f.ConsentCSRF = r.CSRF
	f.LoginAuthenticatedAt = r.AuthenticatedAt
	f.RequestedAt = r.RequestedAt
}

func (f *Flow) setHandledConsentRequest(r *consent.HandledConsentRequest) {
	f.ConsentChallengeID = sqlxx.NullString(r.ID)
	f.CHGrantedScope = r.GrantedScope
	f.CHGrantedAudience = r.GrantedAudience
	f.CHRemember = r.Remember
	f.CHRememberFor = r.RememberFor
	f.CHHandledAt = r.HandledAt
	f.CHWasHandled = r.WasHandled
	f.CHError = r.Error
	f.RequestedAt = r.RequestedAt
	f.LoginAuthenticatedAt = r.AuthenticatedAt
}

func TestFlow_GetLoginRequest(t *testing.T) {
	t.Run("GetLoginRequest should set all fields on its return value", func(t *testing.T) {
		f := Flow{}
		expected := consent.LoginRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		f.setLoginRequest(&expected)
		actual := f.GetLoginRequest()
		assert.Equal(t, expected, *actual)
	})
}

func TestFlow_GetHandledLoginRequest(t *testing.T) {
	t.Run("GetHandledLoginRequest should set all fields on its return value", func(t *testing.T) {
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
	t.Run("NewFlow and GetLoginRequest should use all LoginRequest fields", func(t *testing.T) {
		expected := &consent.LoginRequest{}
		assert.NoError(t, faker.FakeData(expected))
		actual := NewFlow(expected).GetLoginRequest()
		assert.Equal(t, expected, actual)
	})
}

func TestFlow_HandleLoginRequest(t *testing.T) {
	t.Run(
		"HandleLoginRequest should ignore RequestedAt in its argument and copy the other fields",
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
	t.Run("InitializeConsent should transition the flow into FlowStateConsentInitialized", func(t *testing.T) {
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
	t.Run("InitializeConsent should fail when flow.WasHandled is true", func(t *testing.T) {
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

func TestFlow_GetConsentRequest(t *testing.T) {
	t.Run("GetConsentRequest should set all fields on its return value", func(t *testing.T) {
		f := Flow{}
		expected := consent.ConsentRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		f.setConsentRequest(&expected)
		actual := f.GetConsentRequest()
		assert.Equal(t, expected, *actual)
	})
}

func TestFlow_GetHandledConsentRequest(t *testing.T) {
	t.Run("GetHandledConsentRequest should set all fields on its return value", func(t *testing.T) {
		f := Flow{}
		expected := &consent.HandledConsentRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		expected.Session = &consent.ConsentRequestSessionData{}
		expected.ConsentRequest = nil
		f.setHandledConsentRequest(expected)
		actual := f.GetHandledConsentRequest()
		actual.ConsentRequest = nil
		assert.Equal(t, expected, actual)
	})
}
