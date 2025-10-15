// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/x/sqlxx"
)

func TestStateTransition(t *testing.T) {
	t.Run("case=ToStateConsentUnused", func(t *testing.T) {
		testCases := []struct {
			name     string
			flowID   string
			opts     []StateTransitionOption
			expected *Flow
		}{
			{
				name:   "with all options",
				flowID: "test-flow-1",
				opts: []StateTransitionOption{
					WithConsentRequestID("consent-req-123"),
					WithConsentSkip(true),
					WithConsentCSRF("csrf-789"),
					WithID("new-flow-id"),
				},
				expected: &Flow{
					ID:               "new-flow-id",
					State:            FlowStateConsentUnused,
					ConsentRequestID: sqlxx.NullString("consent-req-123"),
					ConsentSkip:      true,
					ConsentCSRF:      sqlxx.NullString("csrf-789"),
				},
			},
			{
				name:   "with partial options",
				flowID: "test-flow-2",
				opts: []StateTransitionOption{
					WithConsentRequestID("consent-req-456"),
				},
				expected: &Flow{
					ID:               "test-flow-2",
					State:            FlowStateConsentUnused,
					ConsentRequestID: sqlxx.NullString("consent-req-456"),
					ConsentSkip:      false,
					ConsentCSRF:      sqlxx.NullString(""),
				},
			},
			{
				name:   "with no options",
				flowID: "test-flow-3",
				opts:   []StateTransitionOption{},
				expected: &Flow{
					ID:               "test-flow-3",
					State:            FlowStateConsentUnused,
					ConsentRequestID: sqlxx.NullString(""),
					ConsentSkip:      false,
					ConsentCSRF:      sqlxx.NullString(""),
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				f := &Flow{
					ID: tc.flowID,
				}

				f.ToStateConsentUnused(tc.opts...)

				assert.Equal(t, tc.expected.ID, f.ID)
				assert.Equal(t, tc.expected.State, f.State)
				assert.Equal(t, tc.expected.ConsentRequestID, f.ConsentRequestID)
				assert.Equal(t, tc.expected.ConsentSkip, f.ConsentSkip)
				assert.Equal(t, tc.expected.ConsentCSRF, f.ConsentCSRF)
			})
		}
	})

	t.Run("case=functional_options_work_independently", func(t *testing.T) {
		f := &Flow{ID: "test-flow"}

		// Test WithConsentRequestID
		WithConsentRequestID("test-consent-id")(f)
		assert.Equal(t, sqlxx.NullString("test-consent-id"), f.ConsentRequestID)

		// Test WithConsentSkip
		WithConsentSkip(true)(f)
		assert.True(t, f.ConsentSkip)

		// Test WithConsentCSRF
		WithConsentCSRF("test-csrf")(f)
		assert.Equal(t, sqlxx.NullString("test-csrf"), f.ConsentCSRF)

		// Test WithID
		WithID("new-id")(f)
		assert.Equal(t, "new-id", f.ID)
	})

	t.Run("case=state_transition_preserves_existing_fields", func(t *testing.T) {
		f := &Flow{
			ID:      "original-id",
			Subject: "test-subject",
			Client:  &client.Client{ID: "test-client"},
		}

		f.ToStateConsentUnused(
			WithConsentRequestID("new-consent-id"),
		)

		// State should be updated
		assert.Equal(t, FlowStateConsentUnused, f.State)
		assert.Equal(t, sqlxx.NullString("new-consent-id"), f.ConsentRequestID)

		// Other fields should be preserved
		assert.Equal(t, "original-id", f.ID)
		assert.Equal(t, "test-subject", f.Subject)
		assert.NotNil(t, f.Client)
		assert.Equal(t, "test-client", f.Client.ID)
	})
}
