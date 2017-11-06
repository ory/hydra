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
	"fmt"
	"net/url"
	"path"
)

func JoinURLStrings(host string, parts ...string) string {
	var trailing string

	last := parts[len(parts)-1]
	if last[len(last)-1:] == "/" {
		trailing = "/"
	}

	u, err := url.Parse(host)
	if err != nil {
		return fmt.Sprintf("%s%s%s", path.Join(append([]string{u.Path}, parts...)...), trailing)
	}

	if u.Path == "" {
		u.Path = "/"
	}
	return fmt.Sprintf("%s://%s%s%s", u.Scheme, u.Host, path.Join(append([]string{u.Path}, parts...)...), trailing)
}
