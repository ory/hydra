// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import (
	"net/http"
)

func GetResponseMeta(w http.ResponseWriter) (status, size int) {
	switch t := w.(type) {
	case interface{ Status() int }:
		status = t.Status()
	}

	switch t := w.(type) {
	case interface{ Size() int }:
		size = t.Size()
	case interface{ Written() int64 }:
		size = int(t.Written())
	}

	return
}
