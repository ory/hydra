// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/mohae/deepcopy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/sqlxx"
)

func (f *Flow) setLoginRequest(r *LoginRequest) {
	f.ID = r.ID
	f.RequestedScope = r.RequestedScope
	f.RequestedAudience = r.RequestedAudience
	f.LoginSkip = r.Skip
	f.Subject = r.Subject
	f.OpenIDConnectContext = r.OpenIDConnectContext
	f.Client = r.Client
	f.ClientID = r.ClientID
	f.RequestURL = r.RequestURL
	f.SessionID = r.SessionID
	f.LoginWasUsed = r.WasHandled
	f.ForceSubjectIdentifier = r.ForceSubjectIdentifier
	f.LoginVerifier = r.Verifier
	f.LoginCSRF = r.CSRF
	f.LoginAuthenticatedAt = r.AuthenticatedAt
	f.RequestedAt = r.RequestedAt
}

func (f *Flow) setHandledLoginRequest(r *HandledLoginRequest) {
	f.ID = r.ID
	f.LoginRemember = r.Remember
	f.LoginRememberFor = r.RememberFor
	f.LoginExtendSessionLifespan = r.ExtendSessionLifespan
	f.ACR = r.ACR
	f.AMR = r.AMR
	f.Subject = r.Subject
	f.IdentityProviderSessionID = sqlxx.NullString(r.IdentityProviderSessionID)
	f.ForceSubjectIdentifier = r.ForceSubjectIdentifier
	f.Context = r.Context
	f.LoginWasUsed = r.WasHandled
	f.LoginError = r.Error
	f.RequestedAt = r.RequestedAt
	f.LoginAuthenticatedAt = r.AuthenticatedAt
}

func (f *Flow) setConsentRequest(r OAuth2ConsentRequest) {
	f.ConsentRequestID = sqlxx.NullString(r.ConsentRequestID)
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
	f.ConsentWasHandled = r.WasHandled
	f.ForceSubjectIdentifier = r.ForceSubjectIdentifier
	f.ConsentVerifier = sqlxx.NullString(r.Verifier)
	f.ConsentCSRF = sqlxx.NullString(r.CSRF)
	f.LoginAuthenticatedAt = r.AuthenticatedAt
	f.RequestedAt = r.RequestedAt
}

func (f *Flow) setHandledConsentRequest(r AcceptOAuth2ConsentRequest) {
	f.ConsentRequestID = sqlxx.NullString(r.ConsentRequestID)
	f.GrantedScope = r.GrantedScope
	f.GrantedAudience = r.GrantedAudience
	f.ConsentRemember = r.Remember
	f.ConsentRememberFor = &r.RememberFor
	f.ConsentHandledAt = r.HandledAt
	f.ConsentWasHandled = r.WasHandled
	f.ConsentError = r.Error
	f.RequestedAt = r.RequestedAt
	f.LoginAuthenticatedAt = r.AuthenticatedAt
	f.SessionIDToken = r.SessionIDToken
	f.SessionAccessToken = r.SessionAccessToken
	if r.Context != nil {
		f.Context = r.Context
	}
}

func (f *Flow) setDeviceRequest(r *DeviceUserAuthRequest) {
	f.DeviceChallengeID = sqlxx.NullString(r.ID)
	f.DeviceCSRF = sqlxx.NullString(r.CSRF)
	f.DeviceVerifier = sqlxx.NullString(r.Verifier)
	f.Client = r.Client
	f.RequestURL = r.RequestURL
	f.RequestedAt = r.RequestedAt
	f.RequestedScope = r.RequestedScope
	f.RequestedAudience = r.RequestedAudience
	f.DeviceWasUsed = sqlxx.NullBool{Bool: r.WasHandled, Valid: true}
	f.DeviceHandledAt = r.HandledAt
}

func (f *Flow) setHandledDeviceRequest(r *HandledDeviceUserAuthRequest) {
	f.DeviceChallengeID = sqlxx.NullString(r.ID)
	f.Client = r.Client
	f.RequestURL = r.RequestURL
	f.RequestedAt = r.RequestedAt
	f.RequestedScope = r.RequestedScope
	f.RequestedAudience = r.RequestedAudience
	f.DeviceError = r.Error
	f.RequestedAt = r.RequestedAt
	f.DeviceCodeRequestID = sqlxx.NullString(r.DeviceCodeRequestID)
	f.DeviceWasUsed = sqlxx.NullBool{Bool: r.WasHandled, Valid: true}
	f.DeviceHandledAt = r.HandledAt
}

func TestFlow_GetDeviceUserAuthRequest(t *testing.T) {
	t.Run("GetDeviceUserAuthRequest should set all fields on its return value", func(t *testing.T) {
		f := Flow{}
		expected := DeviceUserAuthRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		f.setDeviceRequest(&expected)
		actual := f.GetDeviceUserAuthRequest()
		assert.Equal(t, expected, *actual)
	})
}

func TestFlow_GetHandledDeviceUserAuthRequest(t *testing.T) {
	t.Run("GetHandledDeviceUserAuthRequest should set all fields on its return value", func(t *testing.T) {
		f := Flow{}
		expected := HandledDeviceUserAuthRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		f.setHandledDeviceRequest(&expected)
		actual := f.GetHandledDeviceUserAuthRequest()
		assert.NotNil(t, actual.Request)
		expected.Request = nil
		actual.Request = nil
		assert.Equal(t, expected, *actual)
	})
}

func TestFlow_NewDeviceFlow(t *testing.T) {
	t.Run("NewDeviceFlow and GetDeviceUserAuthRequest should use all DeviceUserAuthRequest fields", func(t *testing.T) {
		expected := &DeviceUserAuthRequest{}
		assert.NoError(t, faker.FakeData(expected))
		actual := NewDeviceFlow(expected).GetDeviceUserAuthRequest()
		assert.Equal(t, expected, actual)
	})
}

func TestFlow_HandleDeviceUserAuthRequest(t *testing.T) {
	t.Run(
		"HandleDeviceUserAuthRequest should ignore RequestedAt in its argument and copy the other fields",
		func(t *testing.T) {
			f := Flow{}
			assert.NoError(t, faker.FakeData(&f))
			f.State = DeviceFlowStateInitialized

			r := HandledDeviceUserAuthRequest{}
			assert.NoError(t, faker.FakeData(&r))
			r.ID = f.DeviceChallengeID.String()
			f.DeviceWasUsed = sqlxx.NullBool{Bool: false, Valid: true}
			f.RequestedAudience = r.RequestedAudience
			f.RequestedScope = r.RequestedScope
			f.RequestURL = r.RequestURL

			assert.NoError(t, f.HandleDeviceUserAuthRequest(&r))

			actual := f.GetHandledDeviceUserAuthRequest()
			assert.NotEqual(t, r.RequestedAt, actual.RequestedAt)
			r.Request = f.GetDeviceUserAuthRequest()
			actual.RequestedAt = r.RequestedAt
			assert.Equal(t, r, *actual)
		},
	)
}

func TestFlow_GetLoginRequest(t *testing.T) {
	t.Run("GetLoginRequest should set all fields on its return value", func(t *testing.T) {
		f := Flow{}
		expected := LoginRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		f.setLoginRequest(&expected)
		actual := f.GetLoginRequest()
		assert.Equal(t, expected, *actual)
	})
}

func TestFlow_GetHandledLoginRequest(t *testing.T) {
	t.Run("GetHandledLoginRequest should set all fields on its return value", func(t *testing.T) {
		f := Flow{}
		expected := HandledLoginRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		f.setHandledLoginRequest(&expected)
		actual := f.GetHandledLoginRequest()
		assert.NotNil(t, actual.LoginRequest)
		expected.LoginRequest = nil
		actual.LoginRequest = nil
		assert.Equal(t, expected, actual)
	})
}

func TestFlow_NewFlow(t *testing.T) {
	t.Run("NewFlow and GetLoginRequest should use all LoginRequest fields", func(t *testing.T) {
		expected := &LoginRequest{}
		assert.NoError(t, faker.FakeData(expected))
		actual := NewFlow(expected).GetLoginRequest()
		assert.Equal(t, expected, actual)
	})
}

func TestFlow_UpdateFlowWithHandledLoginRequest(t *testing.T) {
	t.Run(
		"UpdateFlowWithHandledLoginRequest should ignore RequestedAt in its argument and copy the other fields",
		func(t *testing.T) {
			f := Flow{}
			assert.NoError(t, faker.FakeData(&f))
			f.State = FlowStateLoginInitialized

			r := HandledLoginRequest{}
			assert.NoError(t, faker.FakeData(&r))
			r.ID = f.ID
			r.Subject = f.Subject
			r.ForceSubjectIdentifier = f.ForceSubjectIdentifier
			f.LoginWasUsed = false

			assert.NoError(t, f.UpdateFlowWithHandledLoginRequest(&r))

			actual := f.GetHandledLoginRequest()
			assert.NotEqual(t, r.RequestedAt, actual.RequestedAt)
			r.LoginRequest = f.GetLoginRequest()
			actual.RequestedAt = r.RequestedAt
			assert.Equal(t, r, actual)
		},
	)
}

func TestFlow_InvalidateLoginRequest(t *testing.T) {
	t.Run("InvalidateLoginRequest should transition the flow into FlowStateLoginUsed", func(t *testing.T) {
		f := NewFlow(&LoginRequest{
			ID:         "t3-id",
			Subject:    "t3-sub",
			WasHandled: false,
		})
		assert.NoError(t, f.UpdateFlowWithHandledLoginRequest(&HandledLoginRequest{
			ID:         "t3-id",
			Subject:    "t3-sub",
			WasHandled: false,
		}))
		assert.NoError(t, f.InvalidateLoginRequest())
		assert.Equal(t, FlowStateLoginUsed, f.State)
		assert.Equal(t, true, f.LoginWasUsed)
	})
	t.Run("InvalidateLoginRequest should fail when flow.LoginWasUsed is true", func(t *testing.T) {
		f := NewFlow(&LoginRequest{
			ID:         "t3-id",
			Subject:    "t3-sub",
			WasHandled: false,
		})
		assert.NoError(t, f.UpdateFlowWithHandledLoginRequest(&HandledLoginRequest{
			ID:         "t3-id",
			Subject:    "t3-sub",
			WasHandled: true,
		}))
		err := f.InvalidateLoginRequest()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "verifier has already been used")
	})
}

func TestFlow_GetConsentRequest(t *testing.T) {
	t.Run("GetConsentRequest should set all fields on its return value", func(t *testing.T) {
		f := Flow{}
		expected := OAuth2ConsentRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		f.setConsentRequest(expected)
		actual := f.GetConsentRequest(expected.Challenge)
		assert.Equal(t, expected, *actual)
	})
}

func TestFlow_HandleConsentRequest(t *testing.T) {
	f := Flow{}
	require.NoError(t, faker.FakeData(&f))

	expected := AcceptOAuth2ConsentRequest{}
	require.NoError(t, faker.FakeData(&expected))

	expected.ConsentRequestID = string(f.ConsentRequestID)
	expected.HandledAt = sqlxx.NullTime(time.Now())
	expected.RequestedAt = f.RequestedAt
	expected.Session = &AcceptOAuth2ConsentRequestSession{
		IDToken:     sqlxx.MapStringInterface{"claim1": "value1", "claim2": "value2"},
		AccessToken: sqlxx.MapStringInterface{"claim3": "value3", "claim4": "value4"},
	}
	expected.SessionIDToken = expected.Session.IDToken
	expected.SessionAccessToken = expected.Session.AccessToken

	f.State = FlowStateConsentInitialized
	f.ConsentWasHandled = false

	fGood := deepcopy.Copy(f).(Flow)
	eGood := deepcopy.Copy(expected).(AcceptOAuth2ConsentRequest)
	require.NoError(t, f.HandleConsentRequest(&expected))

	t.Run("HandleConsentRequest should fail when already handled", func(t *testing.T) {
		fBad := deepcopy.Copy(fGood).(Flow)
		fBad.ConsentWasHandled = true
		require.Error(t, fBad.HandleConsentRequest(&expected))
	})

	t.Run("HandleConsentRequest should fail when State is FlowStateLoginUsed", func(t *testing.T) {
		fBad := deepcopy.Copy(fGood).(Flow)
		fBad.State = FlowStateLoginUsed
		require.Error(t, fBad.HandleConsentRequest(&expected))
	})

	t.Run("HandleConsentRequest should fail when HandledAt in its argument is zero", func(t *testing.T) {
		f := deepcopy.Copy(fGood).(Flow)
		eBad := deepcopy.Copy(eGood).(AcceptOAuth2ConsentRequest)
		eBad.HandledAt = sqlxx.NullTime(time.Time{})
		require.Error(t, f.HandleConsentRequest(&eBad))
	})

	require.NoError(t, fGood.HandleConsentRequest(&expected))

	actual := f.GetHandledConsentRequest()
	require.NotNil(t, actual.ConsentRequest)
	expected.ConsentRequest = nil
	actual.ConsentRequest = nil
	require.Equal(t, &expected, actual)
}

func TestFlow_GetHandledConsentRequest(t *testing.T) {
	t.Run("GetHandledConsentRequest should set all fields on its return value", func(t *testing.T) {
		f := Flow{}
		expected := AcceptOAuth2ConsentRequest{}

		assert.NoError(t, faker.FakeData(&expected))
		expected.ConsentRequest = nil
		expected.Session = &AcceptOAuth2ConsentRequestSession{
			IDToken:     sqlxx.MapStringInterface{"claim1": "value1", "claim2": "value2"},
			AccessToken: sqlxx.MapStringInterface{"claim3": "value3", "claim4": "value4"},
		}
		expected.SessionIDToken = expected.Session.IDToken
		expected.SessionAccessToken = expected.Session.AccessToken

		f.setHandledConsentRequest(expected)
		actual := f.GetHandledConsentRequest()

		assert.NotNil(t, actual.ConsentRequest)
		actual.ConsentRequest = nil

		assert.Equal(t, expected, *actual)
	})
}
