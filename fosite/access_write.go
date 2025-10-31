// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"context"
	"encoding/json"
	"net/http"
)

func (f *Fosite) WriteAccessResponse(ctx context.Context, rw http.ResponseWriter, requester AccessRequester, responder AccessResponder) {
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")

	js, err := json.Marshal(responder.ToMap())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")

	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write(js)
}
