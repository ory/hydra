// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"encoding/json"
	"net/http"
)

// WriteDeviceResponse writes the device response
func (f *Fosite) WriteDeviceResponse(ctx context.Context, rw http.ResponseWriter, requester DeviceRequester, responder DeviceResponder) {
	// Set custom headers, e.g. "X-MySuperCoolCustomHeader" or "X-DONT-CACHE-ME"...
	wh := rw.Header()
	rh := responder.GetHeader()
	for k := range rh {
		wh.Set(k, rh.Get(k))
	}

	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")

	deviceResponse := &DeviceResponse{
		DeviceCode:              responder.GetDeviceCode(),
		UserCode:                responder.GetUserCode(),
		VerificationURI:         responder.GetVerificationURI(),
		VerificationURIComplete: responder.GetVerificationURIComplete(),
		ExpiresIn:               responder.GetExpiresIn(),
		Interval:                responder.GetInterval(),
	}

	r, err := json.Marshal(deviceResponse)
	rw.Write(r)
	if err != nil {
		http.Error(rw, ErrServerError.WithWrap(err).WithDebug(err.Error()).Error(), http.StatusInternalServerError)
		return
	}
}
