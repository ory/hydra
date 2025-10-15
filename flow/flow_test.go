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
	f.RequestURL = r.RequestURL
	f.SessionID = r.SessionID
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

func TestFlow_HandleDeviceUserAuthRequest(t *testing.T) {
	for _, state := range []State{DeviceFlowStateUnused, DeviceFlowStateInitialized} {
		t.Run("HandleDeviceUserAuthRequest should ignore RequestedAt in its argument and copy the other fields", func(t *testing.T) {
			f := Flow{}
			assert.NoError(t, faker.FakeData(&f))
			f.State = state

			r := HandledDeviceUserAuthRequest{}
			assert.NoError(t, faker.FakeData(&r))
			f.RequestURL = r.RequestURL

			assert.NoError(t, f.HandleDeviceUserAuthRequest(&r))

			assert.WithinDuration(t, time.Time(f.DeviceHandledAt), time.Now(), time.Second)
			assert.Equal(t, r.Client, f.Client)
			assert.EqualValues(t, r.DeviceCodeRequestID, f.DeviceCodeRequestID)
		})
	}

	t.Run("should fail with invalid state", func(t *testing.T) {
		f := Flow{State: FlowStateLoginUnused}
		r := HandledDeviceUserAuthRequest{}
		assert.ErrorContains(t, f.HandleDeviceUserAuthRequest(&r), "invalid flow state")
	})

	t.Run("should fail when in used state", func(t *testing.T) {
		f := Flow{State: DeviceFlowStateUsed}
		r := HandledDeviceUserAuthRequest{}
		assert.ErrorContains(t, f.HandleDeviceUserAuthRequest(&r), "invalid flow state")
	})
}

func TestFlow_GetLoginRequest(t *testing.T) {
	t.Run("GetLoginRequest should set all fields on its return value", func(t *testing.T) {
		expected := LoginRequest{}
		assert.NoError(t, faker.FakeData(&expected))
		f := Flow{State: FlowStateLoginUsed}
		f.setLoginRequest(&expected)

		actual := f.GetLoginRequest()
		assert.Equal(t, expected, *actual)
	})
}

func TestFlow_UpdateFlowWithHandledLoginRequest(t *testing.T) {
	t.Run(
		"HandleLoginRequest should ignore RequestedAt in its argument and copy the other fields",
		func(t *testing.T) {
			f := Flow{}
			assert.NoError(t, faker.FakeData(&f))
			f.State = FlowStateLoginUnused

			r := HandledLoginRequest{}
			assert.NoError(t, faker.FakeData(&r))
			r.Subject = f.Subject
			r.ForceSubjectIdentifier = f.ForceSubjectIdentifier

			assert.NoError(t, f.HandleLoginRequest(&r))

			assert.Equal(t, r.Subject, f.Subject)
			assert.Equal(t, r.ForceSubjectIdentifier, f.ForceSubjectIdentifier)
			assert.Equal(t, r.Remember, f.LoginRemember)
			assert.Equal(t, r.RememberFor, f.LoginRememberFor)
			assert.Equal(t, r.ExtendSessionLifespan, f.LoginExtendSessionLifespan)
			assert.Equal(t, r.ACR, f.ACR)
			assert.Equal(t, r.AMR, f.AMR)
			assert.Equal(t, r.IdentityProviderSessionID, f.IdentityProviderSessionID.String())
			assert.Equal(t, r.Context, f.Context)
		},
	)
}

func TestFlow_InvalidateLoginRequest(t *testing.T) {
	for _, state := range []State{FlowStateLoginUnused, FlowStateLoginInitialized} {
		t.Run("InvalidateLoginRequest should transition the flow into FlowStateLoginUsed", func(t *testing.T) {
			f := Flow{
				ID:      "t3-id",
				Subject: "t3-sub",
				State:   state,
			}
			assert.NoError(t, f.HandleLoginRequest(&HandledLoginRequest{
				Subject: "t3-sub",
			}))
			assert.NoError(t, f.InvalidateLoginRequest())
			assert.Equal(t, FlowStateLoginUsed, f.State)
		})
	}
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

	f.State = FlowStateConsentUnused
	f.ConsentHandledAt = sqlxx.NullTime(time.Now())

	fGood := deepcopy.Copy(f).(Flow)
	require.NoError(t, f.HandleConsentRequest(&expected))

	t.Run("HandleConsentRequest should fail when already handled", func(t *testing.T) {
		fBad := deepcopy.Copy(fGood).(Flow)
		fBad.State = FlowStateConsentUsed
		assert.ErrorContains(t, fBad.HandleConsentRequest(&expected), "invalid flow state")
	})

	t.Run("HandleConsentRequest should fail when State is FlowStateLoginUsed", func(t *testing.T) {
		fBad := deepcopy.Copy(fGood).(Flow)
		fBad.State = FlowStateLoginUsed
		require.ErrorContains(t, fBad.HandleConsentRequest(&expected), "invalid flow state")
	})

	t.Run("HandleConsentRequest should pass with legacy FlowStateConsentInitialized", func(t *testing.T) {
		f := deepcopy.Copy(fGood).(Flow)
		f.State = FlowStateConsentInitialized
		require.NoError(t, f.HandleConsentRequest(&expected))

		assert.Equal(t, expected.GrantedScope, f.GrantedScope)
		assert.Equal(t, expected.GrantedAudience, f.GrantedAudience)
		assert.WithinDuration(t, time.Now(), time.Time(f.ConsentHandledAt), 5*time.Second)
		assert.Nil(t, f.ConsentError)
		assert.EqualValues(t, expected.Session.IDToken, f.SessionIDToken)
		assert.EqualValues(t, expected.Session.AccessToken, f.SessionAccessToken)
	})

	require.NoError(t, fGood.HandleConsentRequest(&expected))

	assert.Equal(t, expected.GrantedScope, fGood.GrantedScope)
	assert.Equal(t, expected.GrantedAudience, fGood.GrantedAudience)
	assert.WithinDuration(t, time.Now(), time.Time(fGood.ConsentHandledAt), 5*time.Second)
	assert.Nil(t, fGood.ConsentError)
	assert.EqualValues(t, expected.Session.IDToken, fGood.SessionIDToken)
	assert.EqualValues(t, expected.Session.AccessToken, fGood.SessionAccessToken)
}

func TestFlow_HandleConsentError(t *testing.T) {
	for _, state := range []State{FlowStateConsentInitialized, FlowStateConsentUnused, FlowStateConsentError} {
		f := Flow{}
		require.NoError(t, faker.FakeData(&f))
		f.State = state

		expected := RequestDeniedError{}
		require.NoError(t, faker.FakeData(&expected))

		require.NoError(t, f.HandleConsentError(&expected))
		assert.Equal(t, FlowStateConsentError, f.State)
		assert.WithinDuration(t, time.Now(), time.Time(f.ConsentHandledAt), 5*time.Second)
		assert.Equal(t, &expected, f.ConsentError)

		assert.Zero(t, f.ConsentRemember)
		assert.Zero(t, f.ConsentRememberFor)
		assert.Zero(t, f.GrantedScope)
		assert.Zero(t, f.GrantedAudience)
	}
}
