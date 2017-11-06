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

package policy

import (
	"net/http/httptest"
	"testing"

	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/herodot"
	"github.com/ory/hydra/compose"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/ladon"
	"github.com/ory/ladon/manager/memory"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockPolicy(t *testing.T) hydra.Policy {
	originalPolicy := &ladon.DefaultPolicy{
		ID:          uuid.New(),
		Description: "description",
		Subjects:    []string{"<peter>"},
		Effect:      ladon.AllowAccess,
		Resources:   []string{"<article|user>"},
		Actions:     []string{"view"},
		Conditions: ladon.Conditions{
			"ip": &ladon.CIDRCondition{
				CIDR: "1234",
			},
			"owner": &ladon.EqualsSubjectCondition{},
		},
	}
	out, err := json.Marshal(originalPolicy)
	require.NoError(t, err)

	var apiPolicy hydra.Policy
	require.NoError(t, json.Unmarshal(out, &apiPolicy))
	out, err = json.Marshal(&apiPolicy)
	require.NoError(t, err)

	var checkPolicy ladon.DefaultPolicy
	require.NoError(t, json.Unmarshal(out, &checkPolicy))
	require.EqualValues(t, checkPolicy.Conditions["ip"], originalPolicy.Conditions["ip"])
	require.EqualValues(t, checkPolicy.Conditions["owner"], originalPolicy.Conditions["owner"])

	return apiPolicy
}

func TestPolicySDK(t *testing.T) {
	localWarden, httpClient := compose.NewMockFirewall("hydra", "alice", fosite.Arguments{scope},
		&ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"alice"},
			Resources: []string{"rn:hydra:policies<.*>"},
			Actions:   []string{"create", "get", "delete", "list", "update"},
			Effect:    ladon.AllowAccess,
		},
	)

	handler := &Handler{
		Manager: &memory.MemoryManager{Policies: map[string]ladon.Policy{}},
		W:       localWarden,
		H:       herodot.NewJSONWriter(nil),
	}

	router := httprouter.New()
	handler.SetRoutes(router)
	server := httptest.NewServer(router)

	client := hydra.NewPolicyApiWithBasePath(server.URL)
	client.Configuration.Transport = httpClient.Transport

	p := mockPolicy(t)

	t.Run("TestPolicyManagement", func(t *testing.T) {
		_, response, err := client.GetPolicy(p.Id)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)

		result, response, err := client.CreatePolicy(p)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.EqualValues(t, p, *result)

		result, response, err = client.GetPolicy(p.Id)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.EqualValues(t, p, *result)

		p.Subjects = []string{"stan"}
		result, response, err = client.UpdatePolicy(p.Id, p)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.EqualValues(t, p, *result)

		results, response, err := client.ListPolicies(10, 0)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Len(t, results, 1)

		result, response, err = client.GetPolicy(p.Id)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.EqualValues(t, p, *result)

		response, err = client.DeletePolicy(p.Id)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, response.StatusCode)

		_, response, err = client.GetPolicy(p.Id)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}
