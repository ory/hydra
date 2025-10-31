// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"net/http"
)

// DeviceResponse represents the device authorization response
type DeviceResponse struct {
	Header                  http.Header
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete,omitempty"`
	ExpiresIn               int64  `json:"expires_in"`
	Interval                int    `json:"interval,omitempty"`
}

// NewDeviceResponse returns a new DeviceResponse
func NewDeviceResponse() *DeviceResponse {
	return &DeviceResponse{}
}

// GetDeviceCode returns the response's device_code
func (d *DeviceResponse) GetDeviceCode() string {
	return d.DeviceCode
}

// SetDeviceCode sets the response's device_code
func (d *DeviceResponse) SetDeviceCode(code string) {
	d.DeviceCode = code
}

// GetUserCode returns the response's user_code
func (d *DeviceResponse) GetUserCode() string {
	return d.UserCode
}

// SetUserCode sets the response's user_code
func (d *DeviceResponse) SetUserCode(code string) {
	d.UserCode = code
}

// GetVerificationURI returns the response's verification uri
func (d *DeviceResponse) GetVerificationURI() string {
	return d.VerificationURI
}

// SetVerificationURI sets the response's verification uri
func (d *DeviceResponse) SetVerificationURI(uri string) {
	d.VerificationURI = uri
}

// GetVerificationURIComplete returns the response's complete verification uri if set
func (d *DeviceResponse) GetVerificationURIComplete() string {
	return d.VerificationURIComplete
}

// SetVerificationURIComplete sets the response's complete verification uri
func (d *DeviceResponse) SetVerificationURIComplete(uri string) {
	d.VerificationURIComplete = uri
}

// GetExpiresIn returns the response's device code and user code lifetime in seconds if set
func (d *DeviceResponse) GetExpiresIn() int64 {
	return d.ExpiresIn
}

// SetExpiresIn sets the response's device code and user code lifetime in seconds
func (d *DeviceResponse) SetExpiresIn(seconds int64) {
	d.ExpiresIn = seconds
}

// GetInterval returns the response's polling interval if set
func (d *DeviceResponse) GetInterval() int {
	return d.Interval
}

// SetInterval sets the response's polling interval
func (d *DeviceResponse) SetInterval(seconds int) {
	d.Interval = seconds
}

// GetHeader returns the response's headers
func (d *DeviceResponse) GetHeader() http.Header {
	return d.Header
}

// AddHeader adds a header to the response
func (d *DeviceResponse) AddHeader(key, value string) {
	d.Header.Add(key, value)
}
