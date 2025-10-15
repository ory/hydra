// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import "github.com/ory/x/sqlxx"

type StateTransitionOption func(*Flow)

func WithConsentRequestID(id string) StateTransitionOption {
	return func(f *Flow) {
		f.ConsentRequestID = sqlxx.NullString(id)
	}
}

func WithConsentSkip(skip bool) StateTransitionOption {
	return func(f *Flow) {
		f.ConsentSkip = skip
	}
}

func WithConsentCSRF(csrf string) StateTransitionOption {
	return func(f *Flow) {
		f.ConsentCSRF = sqlxx.NullString(csrf)
	}
}

func WithID(id string) StateTransitionOption {
	return func(f *Flow) {
		f.ID = id
	}
}

func (f *Flow) ToStateConsentUnused(opts ...StateTransitionOption) {
	f.State = FlowStateConsentUnused

	for _, opt := range opts {
		opt(f)
	}
}
