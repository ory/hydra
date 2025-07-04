// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package httpx

import "net/http"

type chanHandler <-chan http.HandlerFunc

var _ http.Handler = chanHandler(nil)

func (c chanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	(<-c)(w, r)
}

// NewChanHandler returns a new handler and corresponding channel for sending handler funcs.
// Useful for testing. The argument buf specifies the channel capacity, so pass 0 for a sync handler.
func NewChanHandler(buf int) (http.Handler, chan<- http.HandlerFunc) {
	c := make(chan http.HandlerFunc, buf)
	return chanHandler(c), c
}
