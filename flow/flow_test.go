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
	f.ForceSubjectIdentifier = r.ForceSubjectIdentifier
	f.LoginVerifier = r.Verifier
	f.LoginCSRF = r.CSRF
	f.LoginAuthenticatedAt = r.AuthenticatedAt
	f.RequestedAt = r.RequestedAt
}

func (f *Flow) setConsentRequest(r OAuth2ConsentRequest) {
	f.ConsentRequestID = sqlxx.NullString(r.ConsentRequestID)
	f.RequestedScope = r.RequestedScope
	f.RequestedAudience = r.RequestedAudience
	f.ConsentSkip = r.Skip
	f.Subject = r.Subject
	f.OpenIDConnectContext = r.OpenIDConnectContext
	f.Client = r.Client
	f.RequestURL = r.RequestURL
	f.ID = r.LoginChallenge.String()
	f.SessionID = r.LoginSessionID
	f.ACR = r.ACR
	f.AMR = r.AMR
	f.Context = r.Context
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
		expected := LoginRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		f := Flow{
			State:        FlowStateLoginUsed,
			LoginWasUsed: true,
		}
		f.setLoginRequest(&expected)
		expected.WasHandled = true // this depends on the flow state and is not copied from the request

		actual := f.GetLoginRequest()
		assert.Equal(t, expected, *actual)
	})
}

func TestFlow_UpdateFlowWithHandledLoginRequest(t *testing.T) {
	t.Run(
		"UpdateFlowWithHandledLoginRequest should ignore RequestedAt in its argument and copy the other fields",
		func(t *testing.T) {
			f := Flow{}
			assert.NoError(t, faker.FakeData(&f))
			f.State = FlowStateLoginInitialized
			f.LoginWasUsed = false

			r := HandledLoginRequest{}
			assert.NoError(t, faker.FakeData(&r))
			r.Subject = f.Subject
			r.ForceSubjectIdentifier = f.ForceSubjectIdentifier

			assert.NoError(t, f.UpdateFlowWithHandledLoginRequest(&r))

			assert.Equal(t, r.Subject, f.Subject)
			assert.Equal(t, r.ForceSubjectIdentifier, f.ForceSubjectIdentifier)
			assert.Equal(t, r.Remember, f.LoginRemember)
			assert.Equal(t, r.RememberFor, f.LoginRememberFor)
			assert.Equal(t, r.ExtendSessionLifespan, f.LoginExtendSessionLifespan)
			assert.Equal(t, r.ACR, f.ACR)
			assert.Equal(t, r.AMR, f.AMR)
			assert.Equal(t, r.Error, f.LoginError)
			assert.Equal(t, r.IdentityProviderSessionID, f.IdentityProviderSessionID.String())
			assert.Equal(t, r.Context, f.Context)
		},
	)
}

func TestFlow_InvalidateLoginRequest(t *testing.T) {
	t.Run("InvalidateLoginRequest should transition the flow into FlowStateLoginUsed", func(t *testing.T) {
		f := Flow{
			ID:      "t3-id",
			Subject: "t3-sub",
			State:   FlowStateLoginInitialized,
		}
		assert.NoError(t, f.UpdateFlowWithHandledLoginRequest(&HandledLoginRequest{
			Subject: "t3-sub",
		}))
		assert.NoError(t, f.InvalidateLoginRequest())
		assert.Equal(t, FlowStateLoginUsed, f.State)
	})
	t.Run("InvalidateLoginRequest should fail when flow is in used state", func(t *testing.T) {
		f := Flow{
			Subject: "t3-sub",
			State:   FlowStateLoginUsed,
		}
		assert.ErrorContains(t, f.InvalidateLoginRequest(), "invalid flow state")
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

	expected.Session = &AcceptOAuth2ConsentRequestSession{
		IDToken:     sqlxx.MapStringInterface{"claim1": "value1", "claim2": "value2"},
		AccessToken: sqlxx.MapStringInterface{"claim3": "value3", "claim4": "value4"},
	}

	f.State = FlowStateConsentInitialized
	f.ConsentWasHandled = false
	f.ConsentHandledAt = sqlxx.NullTime(time.Now())

	fGood := deepcopy.Copy(f).(Flow)
	fGood.ConsentHandledAt = f.ConsentHandledAt // deepcopy cannot copy time.Time, as it only has unexported fields
	require.NoError(t, f.HandleConsentRequest(&expected))

	t.Run("HandleConsentRequest should fail when already handled", func(t *testing.T) {
		fBad := deepcopy.Copy(fGood).(Flow)
		fBad.State = FlowStateConsentUsed
		fBad.ConsentHandledAt = f.ConsentHandledAt // deepcopy cannot copy time.Time, as it only has unexported fields
		assert.ErrorContains(t, fBad.HandleConsentRequest(&expected), "invalid flow state")
	})

	t.Run("HandleConsentRequest should fail when State is FlowStateLoginUsed", func(t *testing.T) {
		fBad := deepcopy.Copy(fGood).(Flow)
		fBad.State = FlowStateLoginUsed
		fBad.ConsentHandledAt = f.ConsentHandledAt // deepcopy cannot copy time.Time, as it only has unexported fields
		require.ErrorContains(t, fBad.HandleConsentRequest(&expected), "invalid flow state")
	})

	require.NoError(t, fGood.HandleConsentRequest(&expected))

	assert.Equal(t, expected.GrantedScope, fGood.GrantedScope)
	assert.Equal(t, expected.GrantedAudience, fGood.GrantedAudience)
	assert.WithinDuration(t, time.Now(), time.Time(fGood.ConsentHandledAt), 5*time.Second)
	assert.Equal(t, expected.Error, fGood.ConsentError)
	assert.EqualValues(t, expected.Session.IDToken, fGood.SessionIDToken)
	assert.EqualValues(t, expected.Session.AccessToken, fGood.SessionAccessToken)
}
