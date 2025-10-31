// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

type UserCodeState int16

const (
	// User code is active
	UserCodeUnused = UserCodeState(0)
	// User code has been accepted
	UserCodeAccepted = UserCodeState(1)
	// User code has been rejected
	UserCodeRejected = UserCodeState(2)
)

// DeviceRequest is an implementation of DeviceRequester
type DeviceRequest struct {
	UserCodeState UserCodeState
	Request
}

func (d *DeviceRequest) GetUserCodeState() UserCodeState {
	return d.UserCodeState
}

func (d *DeviceRequest) SetUserCodeState(state UserCodeState) {
	d.UserCodeState = state
}

func (d *DeviceRequest) Sanitize(allowedParameters []string) Requester {
	r, _ := d.Request.Sanitize(allowedParameters).(*Request)
	d.Request = *r
	return d
}

// NewDeviceRequest returns a new device request
func NewDeviceRequest() *DeviceRequest {
	return &DeviceRequest{
		UserCodeState: UserCodeUnused,
		Request:       *NewRequest(),
	}
}
