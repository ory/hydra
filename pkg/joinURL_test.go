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

package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinURLStrings(t *testing.T) {
	for k, c := range []struct {
		give []string
		get  string
	}{
		{
			give: []string{"http://localhost/", "/home"},
			get:  "http://localhost/home",
		},
		{
			give: []string{"http://localhost", "/home"},
			get:  "http://localhost/home",
		},
		{
			give: []string{"https://localhost/", "/home"},
			get:  "https://localhost/home",
		},
		{
			give: []string{"http://localhost/", "/home", "home/", "/home/"},
			get:  "http://localhost/home/home/home/",
		},
	} {
		assert.Equal(t, c.get, JoinURLStrings(c.give[0], c.give[1:]...), "Case %d", k)
	}
}
