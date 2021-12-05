/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package trust

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"

	"github.com/ory/hydra/jwk"
)

func TestHelperGrantManagerCreateGetDeleteGrant(m GrantManager) func(t *testing.T) {
	testGenerator := &jwk.RS256Generator{}
	tokenServicePubKey1 := jose.JSONWebKey{}
	tokenServicePubKey2 := jose.JSONWebKey{}
	mikePubKey := jose.JSONWebKey{}

	return func(t *testing.T) {
		keySet, err := testGenerator.Generate("tokenServicePubKey1", "sig")
		require.NoError(t, err)
		tokenServicePubKey1 = keySet.Keys[1]

		keySet, err = testGenerator.Generate("tokenServicePubKey2", "sig")
		require.NoError(t, err)
		tokenServicePubKey2 = keySet.Keys[1]

		keySet, err = testGenerator.Generate("mikePubKey", "sig")
		require.NoError(t, err)
		mikePubKey = keySet.Keys[1]

		storedGrants, err := m.GetGrants(context.TODO(), 100, 0, "")
		require.NoError(t, err)
		assert.Len(t, storedGrants, 0)

		count, err := m.CountGrants(context.TODO())
		require.NoError(t, err)
		assert.Equal(t, 0, count)

		createdAt := time.Now().UTC().Round(time.Second)
		expiresAt := createdAt.AddDate(1, 0, 0)
		grant := Grant{
			ID:      uuid.New().String(),
			Issuer:  "token-service",
			Subject: "bob@example.com",
			Scope:   []string{"openid", "offline"},
			PublicKey: PublicKey{
				Set:   "token-service",
				KeyID: "public:tokenServicePubKey1",
			},
			CreatedAt: createdAt,
			ExpiresAt: expiresAt,
		}
		err = m.CreateGrant(context.TODO(), grant, tokenServicePubKey1)
		require.NoError(t, err)

		storedGrant, err := m.GetConcreteGrant(context.TODO(), grant.ID)
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
			Issuer:  "token-service",
			Subject: "maria@example.com",
			Scope:   []string{"openid"},
			PublicKey: PublicKey{
				Set:   "token-service",
				KeyID: "public:tokenServicePubKey2",
			},
			CreatedAt: createdAt.Add(time.Minute * 5),
			ExpiresAt: createdAt.Add(-time.Minute * 5),
		}
		err = m.CreateGrant(context.TODO(), grant2, tokenServicePubKey2)
		require.NoError(t, err)

		grant3 := Grant{
			ID:      uuid.New().String(),
			Issuer:  "https://mike.example.com",
			Subject: "mike@example.com",
			Scope:   []string{"permissions", "openid", "offline"},
			PublicKey: PublicKey{
				Set:   "https://mike.example.com",
				KeyID: "public:mikePubKey",
			},
			CreatedAt: createdAt.Add(time.Hour),
			ExpiresAt: createdAt.Add(-time.Hour * 24),
		}
		err = m.CreateGrant(context.TODO(), grant3, mikePubKey)
		require.NoError(t, err)

		count, err = m.CountGrants(context.TODO())
		require.NoError(t, err)
		assert.Equal(t, 3, count)

		storedGrants, err = m.GetGrants(context.TODO(), 100, 0, "")
		sort.Slice(storedGrants, func(i, j int) bool {
			return storedGrants[i].CreatedAt.Before(storedGrants[j].CreatedAt)
		})
		require.NoError(t, err)
		require.Len(t, storedGrants, 3)
		assert.Equal(t, grant.ID, storedGrants[0].ID)
		assert.Equal(t, grant2.ID, storedGrants[1].ID)
		assert.Equal(t, grant3.ID, storedGrants[2].ID)

		storedGrants, err = m.GetGrants(context.TODO(), 100, 0, "token-service")
		sort.Slice(storedGrants, func(i, j int) bool {
			return storedGrants[i].CreatedAt.Before(storedGrants[j].CreatedAt)
		})
		require.NoError(t, err)
		require.Len(t, storedGrants, 2)
		assert.Equal(t, grant.ID, storedGrants[0].ID)
		assert.Equal(t, grant2.ID, storedGrants[1].ID)

		err = m.DeleteGrant(context.TODO(), grant.ID)
		require.NoError(t, err)

		_, err = m.GetConcreteGrant(context.TODO(), grant.ID)
		require.Error(t, err)

		count, err = m.CountGrants(context.TODO())
		require.NoError(t, err)
		assert.Equal(t, 2, count)

		err = m.FlushInactiveGrants(context.TODO(), grant2.ExpiresAt, 1000, 100)
		require.NoError(t, err)

		count, err = m.CountGrants(context.TODO())
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		_, err = m.GetConcreteGrant(context.TODO(), grant2.ID)
		assert.NoError(t, err)
	}
}

func TestHelperGrantManagerErrors(m GrantManager) func(t *testing.T) {
	testGenerator := &jwk.RS256Generator{}
	pubKey1 := jose.JSONWebKey{}
	pubKey2 := jose.JSONWebKey{}

	return func(t *testing.T) {
		keySet, err := testGenerator.Generate("pubKey1", "sig")
		require.NoError(t, err)
		pubKey1 = keySet.Keys[1]

		keySet, err = testGenerator.Generate("pubKey2", "sig")
		require.NoError(t, err)
		pubKey2 = keySet.Keys[1]

		createdAt := time.Now()
		expiresAt := createdAt.AddDate(1, 0, 0)
		grant := Grant{
			ID:      uuid.New().String(),
			Issuer:  "issuer",
			Subject: "subject",
			Scope:   []string{"openid", "offline"},
			PublicKey: PublicKey{
				Set:   "set",
				KeyID: "public:pubKey1",
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
			Set:   "set",
			KeyID: "public:pubKey2",
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
