// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceResponse(t *testing.T) {
	r := NewDeviceResponse()
	r.SetDeviceCode("device_code")
	r.SetUserCode("user_code")
	r.SetExpiresIn(5)
	r.SetVerificationURI("https://www.example.com")
	r.SetVerificationURIComplete("https://www.example.com?code=user_code")
	r.SetInterval(5)
	assert.Equal(t, "device_code", r.GetDeviceCode())
	assert.Equal(t, "user_code", r.GetUserCode())
	assert.Equal(t, int64(5), r.GetExpiresIn())
	assert.Equal(t, "https://www.example.com", r.GetVerificationURI())
	assert.Equal(t, "https://www.example.com?code=user_code", r.GetVerificationURIComplete())
	assert.Equal(t, 5, r.GetInterval())
}
