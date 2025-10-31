// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/ory/hydra/v2/fosite"
)

func TestWriteDeviceUserResponse(t *testing.T) {
	oauth2 := &Fosite{Config: &Config{
		DeviceAndUserCodeLifespan:      time.Minute,
		DeviceAuthTokenPollingInterval: time.Minute,
		DeviceVerificationURL:          "http://ory.sh",
	}}
	ctx := context.Background()

	rw := httptest.NewRecorder()
	ar := &DeviceRequest{}
	resp := &DeviceResponse{}
	resp.SetUserCode("AAAA")
	resp.SetDeviceCode("BBBB")
	resp.SetInterval(int(
		oauth2.Config.GetDeviceAuthTokenPollingInterval(ctx).Round(time.Second).Seconds(),
	))
	resp.SetExpiresIn(int64(
		oauth2.Config.GetDeviceAndUserCodeLifespan(ctx),
	))
	resp.SetVerificationURI(oauth2.Config.GetDeviceVerificationURL(ctx))
	resp.SetVerificationURIComplete(
		oauth2.Config.GetDeviceVerificationURL(ctx) + "?user_code=" + resp.GetUserCode(),
	)

	oauth2.WriteDeviceResponse(context.Background(), rw, ar, resp)

	assert.Equal(t, 200, rw.Code)

	body, err := io.ReadAll(rw.Body)
	require.NoError(t, err)

	wroteDeviceResponse := DeviceResponse{}
	err = json.Unmarshal(body, &wroteDeviceResponse)
	require.NoError(t, err)

	assert.Equal(t, resp.GetUserCode(), wroteDeviceResponse.UserCode)
	assert.Equal(t, resp.GetDeviceCode(), wroteDeviceResponse.DeviceCode)
	assert.Equal(t, resp.GetVerificationURI(), wroteDeviceResponse.VerificationURI)
	assert.Equal(t, resp.GetVerificationURIComplete(), wroteDeviceResponse.VerificationURIComplete)
	assert.Equal(t, resp.GetInterval(), wroteDeviceResponse.Interval)
	assert.Equal(t, resp.GetExpiresIn(), wroteDeviceResponse.ExpiresIn)
}
