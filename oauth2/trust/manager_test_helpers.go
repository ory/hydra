// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package trust

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/ory/x/josex"

	"github.com/go-jose/go-jose/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/jwk"
)

func TestHelperGrantManagerCreateGetDeleteGrant(t1 GrantManager, km jwk.Manager, parallel bool) func(t *testing.T) {
	tokenServicePubKey1 := jose.JSONWebKey{}
	tokenServicePubKey2 := jose.JSONWebKey{}
	mikePubKey := jose.JSONWebKey{}

	return func(t *testing.T) {
		if parallel {
			t.Parallel()
		}
		kid1, kid2 := uuid.NewString(), uuid.NewString()
		kid3 := uuid.NewString()
		set := uuid.NewString()

		keySet, err := km.GenerateAndPersistKeySet(context.Background(), set, kid1, string(jose.RS256), "sig")
		require.NoError(t, err)
		tokenServicePubKey1 = josex.ToPublicKey(&keySet.Keys[0])

		keySet, err = km.GenerateAndPersistKeySet(context.Background(), set, kid2, string(jose.RS256), "sig")
		require.NoError(t, err)
		tokenServicePubKey2 = josex.ToPublicKey(&keySet.Keys[0])

		keySet, err = km.GenerateAndPersistKeySet(context.Background(), "https://mike.example.com", kid3, string(jose.RS256), "sig")
		require.NoError(t, err)
		mikePubKey = josex.ToPublicKey(&keySet.Keys[0])

		storedGrants, err := t1.GetGrants(context.TODO(), 100, 0, "")
		require.NoError(t, err)
		assert.Len(t, storedGrants, 0)

		count, err := t1.CountGrants(context.TODO())
		require.NoError(t, err)
		assert.Equal(t, 0, count)

		createdAt := time.Now().UTC().Round(time.Second)
		expiresAt := createdAt.AddDate(1, 0, 0)
		grant := Grant{
			ID:      uuid.New().String(),
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

		err = t1.CreateGrant(context.TODO(), grant, tokenServicePubKey1)
		require.NoError(t, err)

		storedGrant, err := t1.GetConcreteGrant(context.TODO(), grant.ID)
		require.NoError(t, err)
		assert.Equal(t, grant.ID, storedGrant.ID)
		assert.Equal(t, grant.Issuer, storedGrant.Issuer)
		assert.Equal(t, grant.Subject, storedGrant.Subject)
		assert.Equal(t, grant.Scope, storedGrant.Scope)
		assert.Equal(t, grant.PublicKey, storedGrant.PublicKey)
		assert.Equal(t, grant.CreatedAt.Format(time.RFC3339), storedGrant.CreatedAt.Format(time.RFC3339))
		assert.Equal(t, grant.ExpiresAt.Format(time.RFC3339), storedGrant.ExpiresAt.Format(time.RFC3339))

		grant2 := Grant{
			ID:      uuid.New().String(),
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
		err = t1.CreateGrant(context.TODO(), grant2, tokenServicePubKey2)
		require.NoError(t, err)

		grant3 := Grant{
			ID:      uuid.New().String(),
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

		err = t1.CreateGrant(context.TODO(), grant3, mikePubKey)
		require.NoError(t, err)

		count, err = t1.CountGrants(context.TODO())
		require.NoError(t, err)
		assert.Equal(t, 3, count)

		storedGrants, err = t1.GetGrants(context.TODO(), 100, 0, "")
		sort.Slice(storedGrants, func(i, j int) bool {
			return storedGrants[i].CreatedAt.Before(storedGrants[j].CreatedAt)
		})
		require.NoError(t, err)
		require.Len(t, storedGrants, 3)
		assert.Equal(t, grant.ID, storedGrants[0].ID)
		assert.Equal(t, grant2.ID, storedGrants[1].ID)
		assert.Equal(t, grant3.ID, storedGrants[2].ID)

		storedGrants, err = t1.GetGrants(context.TODO(), 100, 0, set)
		sort.Slice(storedGrants, func(i, j int) bool {
			return storedGrants[i].CreatedAt.Before(storedGrants[j].CreatedAt)
		})
		require.NoError(t, err)
		require.Len(t, storedGrants, 2)
		assert.Equal(t, grant.ID, storedGrants[0].ID)
		assert.Equal(t, grant2.ID, storedGrants[1].ID)

		err = t1.DeleteGrant(context.TODO(), grant.ID)
		require.NoError(t, err)

		_, err = t1.GetConcreteGrant(context.TODO(), grant.ID)
		require.Error(t, err)

		count, err = t1.CountGrants(context.TODO())
		require.NoError(t, err)
		assert.Equal(t, 2, count)

		err = t1.FlushInactiveGrants(context.TODO(), grant2.ExpiresAt, 1000, 100)
		require.NoError(t, err)

		count, err = t1.CountGrants(context.TODO())
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		_, err = t1.GetConcreteGrant(context.TODO(), grant2.ID)
		assert.NoError(t, err)
	}
}

func TestHelperGrantManagerErrors(m GrantManager, km jwk.Manager, parallel bool) func(t *testing.T) {
	pubKey1 := jose.JSONWebKey{}
	pubKey2 := jose.JSONWebKey{}

	return func(t *testing.T) {
		set := uuid.NewString()
		kid1, kid2 := uuid.NewString(), uuid.NewString()

		t.Parallel()
		keySet, err := km.GenerateAndPersistKeySet(context.Background(), set, kid1, string(jose.RS256), "sig")
		require.NoError(t, err)
		pubKey1 = josex.ToPublicKey(&keySet.Keys[0])

		keySet, err = km.GenerateAndPersistKeySet(context.Background(), set, kid2, string(jose.RS256), "sig")
		require.NoError(t, err)
		pubKey2 = josex.ToPublicKey(&keySet.Keys[0])

		createdAt := time.Now()
		expiresAt := createdAt.AddDate(1, 0, 0)
		grant := Grant{
			ID:      uuid.New().String(),
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

		err = m.CreateGrant(context.TODO(), grant, pubKey1)
		require.NoError(t, err)

		grant.ID = uuid.New().String()
		err = m.CreateGrant(context.TODO(), grant, pubKey1)
		require.Error(t, err, "error expected, because combination of issuer + subject + key_id must be unique")

		grant2 := grant
		grant2.PublicKey = PublicKey{
			Set:   set,
			KeyID: kid2,
		}
		err = m.CreateGrant(context.TODO(), grant2, pubKey2)
		require.NoError(t, err)

		nonExistingGrantID := uuid.New().String()
		err = m.DeleteGrant(context.TODO(), nonExistingGrantID)
		require.Error(t, err, "expect error, when deleting non-existing grant")

		_, err = m.GetConcreteGrant(context.TODO(), nonExistingGrantID)
		require.Error(t, err, "expect error, when fetching non-existing grant")
	}
}
