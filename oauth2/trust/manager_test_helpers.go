// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package trust

import (
	"sort"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/josex"
	"github.com/ory/x/sqlcon"

	"github.com/ory/hydra/v2/jwk"
)

func TestHelperGrantManagerCreateGetDeleteGrant(t1 GrantManager, km jwk.Manager, parallel bool) func(t *testing.T) {
	return func(t *testing.T) {
		if parallel {
			t.Parallel()
		}
		kid1, kid2 := uuid.Must(uuid.NewV4()).String(), uuid.Must(uuid.NewV4()).String()
		kid3 := uuid.Must(uuid.NewV4()).String()
		set := uuid.Must(uuid.NewV4()).String()

		key1, err := jwk.GenerateJWK(jose.RS256, kid1, "sig")
		require.NoError(t, err)
		tokenServicePubKey1 := josex.ToPublicKey(&key1.Keys[0])

		key2, err := jwk.GenerateJWK(jose.RS256, kid2, "sig")
		require.NoError(t, err)
		tokenServicePubKey2 := josex.ToPublicKey(&key2.Keys[0])

		key3, err := jwk.GenerateJWK(jose.RS256, kid3, "sig")
		require.NoError(t, err)
		mikePubKey := josex.ToPublicKey(&key3.Keys[0])

		storedGrants, nextPage, err := t1.GetGrants(t.Context(), "")
		require.NoError(t, err)
		assert.Len(t, storedGrants, 0)
		assert.True(t, nextPage.IsLast())

		createdAt := time.Now().UTC().Round(time.Second)
		expiresAt := createdAt.AddDate(1, 0, 0)
		grant := Grant{
			ID:      uuid.Must(uuid.NewV4()),
			Issuer:  set,
			Subject: "bob@example.com",
			Scope:   []string{"openid", "offline"},
			PublicKey: PublicKey{
				Set:   set,
				KeyID: kid1,
			},
			CreatedAt: createdAt,
			ExpiresAt: expiresAt,
		}

		require.NoError(t, t1.CreateGrant(t.Context(), grant, tokenServicePubKey1))

		storedGrant, err := t1.GetConcreteGrant(t.Context(), grant.ID)
		require.NoError(t, err)
		assert.Equal(t, grant.ID, storedGrant.ID)
		assert.Equal(t, grant.Issuer, storedGrant.Issuer)
		assert.Equal(t, grant.Subject, storedGrant.Subject)
		assert.Equal(t, grant.Scope, storedGrant.Scope)
		assert.Equal(t, grant.PublicKey, storedGrant.PublicKey)
		assert.Equal(t, grant.CreatedAt.Format(time.RFC3339), storedGrant.CreatedAt.Format(time.RFC3339))
		assert.Equal(t, grant.ExpiresAt.Format(time.RFC3339), storedGrant.ExpiresAt.Format(time.RFC3339))

		grant2 := Grant{
			ID:      uuid.Must(uuid.NewV4()),
			Issuer:  set,
			Subject: "maria@example.com",
			Scope:   []string{"openid"},
			PublicKey: PublicKey{
				Set:   set,
				KeyID: kid2,
			},
			CreatedAt: createdAt.Add(time.Minute * 5),
			ExpiresAt: createdAt.Add(-time.Minute * 5),
		}
		require.NoError(t, t1.CreateGrant(t.Context(), grant2, tokenServicePubKey2))

		grant3 := Grant{
			ID:      uuid.Must(uuid.NewV4()),
			Issuer:  "https://mike.example.com",
			Subject: "mike@example.com",
			Scope:   []string{"permissions", "openid", "offline"},
			PublicKey: PublicKey{
				Set:   "https://mike.example.com",
				KeyID: kid3,
			},
			CreatedAt: createdAt.Add(time.Hour),
			ExpiresAt: createdAt.Add(-time.Hour * 24),
		}

		require.NoError(t, t1.CreateGrant(t.Context(), grant3, mikePubKey))

		storedGrants, nextPage, err = t1.GetGrants(t.Context(), "")
		sort.Slice(storedGrants, func(i, j int) bool {
			return storedGrants[i].CreatedAt.Before(storedGrants[j].CreatedAt)
		})
		require.NoError(t, err)
		require.Len(t, storedGrants, 3)
		assert.Equal(t, grant.ID, storedGrants[0].ID)
		assert.Equal(t, grant2.ID, storedGrants[1].ID)
		assert.Equal(t, grant3.ID, storedGrants[2].ID)
		assert.True(t, nextPage.IsLast())

		storedGrants, nextPage, err = t1.GetGrants(t.Context(), set)
		sort.Slice(storedGrants, func(i, j int) bool {
			return storedGrants[i].CreatedAt.Before(storedGrants[j].CreatedAt)
		})
		require.NoError(t, err)
		require.Len(t, storedGrants, 2)
		assert.Equal(t, grant.ID, storedGrants[0].ID)
		assert.Equal(t, grant2.ID, storedGrants[1].ID)
		assert.True(t, nextPage.IsLast())

		require.NoError(t, t1.DeleteGrant(t.Context(), grant.ID))

		_, err = t1.GetConcreteGrant(t.Context(), grant.ID)
		require.ErrorIs(t, err, sqlcon.ErrNoRows)

		require.NoError(t, t1.FlushInactiveGrants(t.Context(), grant2.ExpiresAt, 1000, 100))

		_, err = t1.GetConcreteGrant(t.Context(), grant2.ID)
		assert.NoError(t, err)

		_, err = t1.GetConcreteGrant(t.Context(), grant3.ID)
		assert.ErrorIs(t, err, sqlcon.ErrNoRows)
	}
}

func TestHelperGrantManagerErrors(m GrantManager, km jwk.Manager) func(t *testing.T) {
	return func(t *testing.T) {
		set := uuid.Must(uuid.NewV4()).String()
		kid1, kid2 := uuid.Must(uuid.NewV4()).String(), uuid.Must(uuid.NewV4()).String()

		t.Parallel()

		key1, err := jwk.GenerateJWK(jose.RS256, kid1, "sig")
		require.NoError(t, err)
		pubKey1 := josex.ToPublicKey(&key1.Keys[0])

		key2, err := jwk.GenerateJWK(jose.RS256, kid2, "sig")
		require.NoError(t, err)
		pubKey2 := josex.ToPublicKey(&key2.Keys[0])

		createdAt := time.Now()
		expiresAt := createdAt.AddDate(1, 0, 0)
		grant := Grant{
			ID:      uuid.Must(uuid.NewV4()),
			Issuer:  "issuer",
			Subject: "subject",
			Scope:   []string{"openid", "offline"},
			PublicKey: PublicKey{
				Set:   set,
				KeyID: kid1,
			},
			CreatedAt: createdAt,
			ExpiresAt: expiresAt,
		}

		require.NoError(t, m.CreateGrant(t.Context(), grant, pubKey1))

		grant.ID = uuid.Must(uuid.NewV4())
		err = m.CreateGrant(t.Context(), grant, pubKey1)
		require.ErrorIs(t, err, sqlcon.ErrUniqueViolation, "error expected, because combination of issuer + subject + key_id must be unique")

		grant2 := grant
		grant2.PublicKey = PublicKey{
			Set:   set,
			KeyID: kid2,
		}
		require.NoError(t, m.CreateGrant(t.Context(), grant2, pubKey2))

		nonExistingGrantID := uuid.Must(uuid.NewV4())
		err = m.DeleteGrant(t.Context(), nonExistingGrantID)
		require.Error(t, err, "expect error, when deleting non-existing grant")

		_, err = m.GetConcreteGrant(t.Context(), nonExistingGrantID)
		require.Error(t, err, "expect error, when fetching non-existing grant")
	}
}
