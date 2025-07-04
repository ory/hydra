// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package serverx

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// PermanentRedirect permanently redirects (302) a path to another one.
func PermanentRedirect(to string) func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.Redirect(rw, r, to, http.StatusPermanentRedirect)
	}
}
