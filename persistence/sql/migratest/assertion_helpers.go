package migratest

import (
	"testing"
	"time"

	"github.com/instana/testify/require"

	"github.com/ory/hydra/flow"
	"github.com/ory/x/sqlxx"
)

func fixturizeFlow(t *testing.T, f *flow.Flow) {
	require.NotZero(t, f.ClientID)
	f.ClientID = ""
	require.NotNil(t, f.Client)
	f.Client = nil
	recently := time.Now().Add(-time.Minute)
	require.Greater(t, time.Time(f.LoginInitializedAt).UnixNano(), recently.UnixNano())
	f.LoginInitializedAt = sqlxx.NullTime{}
	require.True(t, f.RequestedAt.After(recently))
	f.RequestedAt = time.Time{}
	require.True(t, time.Time(f.LoginAuthenticatedAt).After(recently))
	f.LoginAuthenticatedAt = sqlxx.NullTime{}
	f.ConsentHandledAt = sqlxx.NullTime{}
}
