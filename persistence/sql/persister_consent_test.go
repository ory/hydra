// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/flow"
	"github.com/ory/hydra/v2/internal/testhelpers"
	hydrasql "github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/networkx"
	"github.com/ory/x/servicelocatorx"
	"github.com/ory/x/sqlxx"
	"github.com/ory/x/uuidx"
)

// TestRevokeSubjectLoginSessionBatched is intentionally NOT parallel: it lowers
// the package-level RevokeSubjectLoginSessionBatchSize, which must not race with
// the parallel TestManagers suite. Non-parallel tests complete in the sequential
// phase before paused parallel tests resume, so the mutation is safe.
func TestRevokeSubjectLoginSessionBatched(t *testing.T) {
	original := hydrasql.RevokeSubjectLoginSessionBatchSize
	hydrasql.RevokeSubjectLoginSessionBatchSize = 3
	t.Cleanup(func() { hydrasql.RevokeSubjectLoginSessionBatchSize = original })

	// The two networks let us assert the batched DELETE never crosses the nid
	// boundary. TestContextualizer reads the nid from the context, so the same
	// manager can act on either network.
	nid1, nid2 := uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())
	reg := testhelpers.NewRegistrySQLFromURL(t, dbal.NewSQLiteTestDatabase(t), true, true,
		driver.DisableValidation(),
		driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(&contextx.TestContextualizer{})))

	ctx1 := contextx.SetNIDContext(t.Context(), nid1)
	ctx2 := contextx.SetNIDContext(t.Context(), nid2)
	require.NoError(t, reg.Persister().Connection(ctx1).Create(&networkx.Network{ID: nid1}))
	require.NoError(t, reg.Persister().Connection(ctx1).Create(&networkx.Network{ID: nid2}))

	m := reg.LoginManager()

	seed := func(t *testing.T, ctx context.Context, subject string, n int) []string {
		t.Helper()
		ids := make([]string, 0, n)
		for range n {
			id := uuidx.NewV4().String()
			ids = append(ids, id)
			require.NoError(t, m.ConfirmLoginSession(ctx, &flow.LoginSession{
				ID:              id,
				AuthenticatedAt: sqlxx.NullTime(time.Now()),
				Subject:         subject,
				Remember:        true,
			}))
		}
		return ids
	}

	deleted := func(t *testing.T, ctx context.Context, id string) bool {
		t.Helper()
		_, err := m.GetRememberedLoginSession(ctx, id)
		if err == nil {
			return false
		}
		require.ErrorIs(t, err, x.ErrNotFound)
		return true
	}

	t.Run("case=deletes all sessions across multiple batches", func(t *testing.T) {
		// 7 sessions with batch size 3 forces three iterations (3 + 3 + 1).
		subjectA := uuidx.NewV4().String()
		subjectB := uuidx.NewV4().String()
		aIDs := seed(t, ctx1, subjectA, 7)
		bIDs := seed(t, ctx1, subjectB, 1)

		require.NoError(t, m.RevokeSubjectLoginSession(ctx1, subjectA))

		for _, id := range aIDs {
			assert.True(t, deleted(t, ctx1, id), "session %s should be deleted", id)
		}

		// The unrelated subject's session must survive.
		assert.False(t, deleted(t, ctx1, bIDs[0]), "unrelated subject must survive")
	})

	t.Run("case=exact batch-size boundary drains and exits", func(t *testing.T) {
		// n == batchSize forces a second query that deletes zero rows and then
		// exits. This guards the loop-exit predicate against an off-by-one
		// (deleted < batchSize vs deleted <= batchSize).
		subject := uuidx.NewV4().String()
		ids := seed(t, ctx1, subject, 3)

		require.NoError(t, m.RevokeSubjectLoginSession(ctx1, subject))

		for _, id := range ids {
			assert.True(t, deleted(t, ctx1, id), "session %s should be deleted", id)
		}
	})

	t.Run("case=subject without sessions is a no-op", func(t *testing.T) {
		assert.NoError(t, m.RevokeSubjectLoginSession(ctx1, uuidx.NewV4().String()))
	})

	t.Run("case=does not cross the network boundary", func(t *testing.T) {
		// The same subject has sessions in both networks. Revoking in nid1 must
		// leave the nid2 sessions untouched.
		subject := uuidx.NewV4().String()
		ids1 := seed(t, ctx1, subject, 4)
		ids2 := seed(t, ctx2, subject, 2)

		require.NoError(t, m.RevokeSubjectLoginSession(ctx1, subject))

		for _, id := range ids1 {
			assert.True(t, deleted(t, ctx1, id), "nid1 session %s should be deleted", id)
		}
		for _, id := range ids2 {
			assert.False(t, deleted(t, ctx2, id), "nid2 session %s must survive", id)
		}
	})

	t.Run("case=cancelled before the first batch deletes nothing", func(t *testing.T) {
		subject := uuidx.NewV4().String()
		ids := seed(t, ctx1, subject, 2)

		ctx, cancel := context.WithCancel(ctx1)
		cancel()

		err := m.RevokeSubjectLoginSession(ctx, subject)
		assert.ErrorIs(t, err, context.Canceled)

		// Nothing should have been deleted because cancellation is checked first.
		for _, id := range ids {
			assert.False(t, deleted(t, ctx1, id), "session %s should still exist", id)
		}
	})

	t.Run("case=cancelled between batches keeps committed batches and can retry", func(t *testing.T) {
		// Cancel after the first batch commits. The function returns the
		// cancellation error, the already-deleted batch stays gone, and a retry
		// with a live context finishes the rest. 7 sessions with batch size 3
		// leaves four behind after one batch.
		subject := uuidx.NewV4().String()
		ids := seed(t, ctx1, subject, 7)

		ctx := &cancelAfterLoopGuard{Context: ctx1, done: make(chan struct{})}
		err := m.RevokeSubjectLoginSession(ctx, subject)
		assert.ErrorIs(t, err, context.Canceled)

		remaining := 0
		for _, id := range ids {
			if !deleted(t, ctx1, id) {
				remaining++
			}
		}
		// One full batch committed before cancellation; the rest survive.
		assert.Equal(t, 4, remaining, "one committed batch should stay deleted")

		// The caller can retry to finish draining the subject.
		require.NoError(t, m.RevokeSubjectLoginSession(ctx1, subject))
		for _, id := range ids {
			assert.True(t, deleted(t, ctx1, id), "session %s should be deleted after retry", id)
		}
	})
}

// cancelAfterLoopGuard reports the context as live on the first Err() check and
// cancelled on every subsequent one. RevokeSubjectLoginSession calls Err() once
// at the top of each batch iteration, so the first batch runs and the second is
// rejected. Done() never fires, so the database driver never aborts the
// in-flight query mid-batch — only the loop guard observes the cancellation.
type cancelAfterLoopGuard struct {
	context.Context
	done chan struct{}

	mu    sync.Mutex
	calls int
}

func (c *cancelAfterLoopGuard) Done() <-chan struct{} { return c.done }

func (c *cancelAfterLoopGuard) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.calls++
	if c.calls > 1 {
		return context.Canceled
	}
	return nil
}
