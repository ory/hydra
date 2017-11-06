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

package group

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelperManagers(m Manager) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		_, err := m.GetGroup("4321")
		assert.NotNil(t, err)

		c := &Group{
			ID:      "1",
			Members: []string{"bar", "foo"},
		}
		assert.NoError(t, m.CreateGroup(c))
		assert.NoError(t, m.CreateGroup(&Group{
			ID:      "2",
			Members: []string{"foo"},
		}))

		d, err := m.GetGroup("1")
		require.NoError(t, err)
		assert.EqualValues(t, c.Members, d.Members)
		assert.EqualValues(t, c.ID, d.ID)

		ds, err := m.FindGroupsByMember("foo")
		require.NoError(t, err)
		assert.Len(t, ds, 2)

		assert.NoError(t, m.AddGroupMembers("1", []string{"baz"}))

		ds, err = m.FindGroupsByMember("baz")
		require.NoError(t, err)
		assert.Len(t, ds, 1)

		assert.NoError(t, m.RemoveGroupMembers("1", []string{"baz"}))
		ds, err = m.FindGroupsByMember("baz")
		require.NoError(t, err)
		assert.Len(t, ds, 0)

		assert.NoError(t, m.DeleteGroup("1"))
		_, err = m.GetGroup("1")
		require.NotNil(t, err)
	}
}
