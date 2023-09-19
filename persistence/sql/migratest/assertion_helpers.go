// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package migratest

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/flow"
	testhelpersuuid "github.com/ory/hydra/v2/internal/testhelpers/uuid"
	"github.com/ory/x/sqlxx"
)

func fixturizeFlow(t *testing.T, f *flow.Flow) {
	testhelpersuuid.AssertUUID(t, f.NID)
	f.NID = uuid.Nil
	require.NotZero(t, f.ClientID)
	f.ClientID = ""
	require.NotNil(t, f.Client)
	f.Client = nil
	recently := time.Now().Add(-2 * time.Minute)
	require.Greater(t, time.Time(f.LoginInitializedAt).UnixNano(), recently.UnixNano())
	f.LoginInitializedAt = sqlxx.NullTime{}
	require.True(t, f.RequestedAt.After(recently))
	f.RequestedAt = time.Time{}
	require.True(t, time.Time(f.LoginAuthenticatedAt).After(recently))
	f.LoginAuthenticatedAt = sqlxx.NullTime{}
	f.ConsentHandledAt = sqlxx.NullTime{}
}
