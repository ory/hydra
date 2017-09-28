package warden_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActionAllowed(t *testing.T) {
	for n, w := range wardens {
		t.Run("warden="+n, func(t *testing.T) {
			for k, c := range accessRequestTokenTestCases {
				t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
					ctx, err := w.TokenAllowed(context.Background(), c.token, c.req, c.scopes...)
					if c.expectErr {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
					}

					if err == nil && c.assert != nil {
						c.assert(t, ctx)
					}
				})
			}
		})
	}
}

func TestAllowed(t *testing.T) {
	for n, w := range wardens {
		t.Run("warden="+n, func(t *testing.T) {
			for k, c := range accessRequestTestCases {
				t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
					err := w.IsAllowed(context.Background(), c.req)
					if c.expectErr {
						require.Error(t, err)
					} else {
						require.NoError(t, err)
					}
				})
			}
		})
	}
}
