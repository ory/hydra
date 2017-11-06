// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package group_test

import (
	"net/http/httptest"
	"testing"

	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/compose"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	. "github.com/ory/hydra/warden/group"
	"github.com/ory/ladon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupSDK(t *testing.T) {
	clientManagers["memory"] = &MemoryManager{
		Groups: map[string]Group{},
	}

	localWarden, httpClient := compose.NewMockFirewall("foo", "alice", fosite.Arguments{Scope}, &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"alice"},
		Resources: []string{"rn:hydra:warden<.*>"},
		Actions:   []string{"list", "create", "get", "delete", "update", "members.add", "members.remove"},
		Effect:    ladon.AllowAccess,
	})

	handler := &Handler{
		Manager: &MemoryManager{
			Groups: map[string]Group{},
		},
		H: herodot.NewJSONWriter(nil),
		W: localWarden,
	}

	router := httprouter.New()
	handler.SetRoutes(router)
	server := httptest.NewServer(router)

	client := hydra.NewWardenApiWithBasePath(server.URL)
	client.Configuration.Transport = httpClient.Transport

	t.Run("flows", func(*testing.T) {
		_, response, err := client.GetGroup("4321")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		firstGroup := hydra.Group{Id: "1", Members: []string{"bar", "foo"}}
		result, response, err := client.CreateGroup(firstGroup)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.EqualValues(t, firstGroup, *result)

		secondGroup := hydra.Group{Id: "2", Members: []string{"foo"}}
		result, response, err = client.CreateGroup(secondGroup)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.EqualValues(t, secondGroup, *result)

		result, response, err = client.GetGroup("1")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.EqualValues(t, firstGroup, *result)

		results, response, err := client.FindGroupsByMember("foo")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Len(t, results, 2)

		client.AddMembersToGroup("1", hydra.GroupMembers{Members: []string{"baz"}})

		results, response, err = client.FindGroupsByMember("baz")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Len(t, results, 1)

		response, err = client.RemoveMembersFromGroup("1", hydra.GroupMembers{Members: []string{"baz"}})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, response.StatusCode)

		results, response, err = client.FindGroupsByMember("baz")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Len(t, results, 0)

		response, err = client.DeleteGroup("1")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, response.StatusCode)

		_, response, err = client.GetGroup("4321")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}
