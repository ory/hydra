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
	"net/url"
	"path"
)

func CopyURL(u *url.URL) *url.URL {
	a := new(url.URL)
	*a = *u
	return a
}

func JoinURL(u *url.URL, args ...string) (ep *url.URL) {
	ep = CopyURL(u)
	ep.Path = path.Join(append([]string{ep.Path}, args...)...)
	return ep
}
