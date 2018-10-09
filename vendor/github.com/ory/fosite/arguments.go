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
 *
 */

package fosite

import "strings"

type Arguments []string

func (r Arguments) Matches(items ...string) bool {
	found := make(map[string]bool)
	for _, item := range items {
		if !StringInSlice(item, r) {
			return false
		}
		found[item] = true
	}

	return len(found) == len(r) && len(r) == len(items)
}

func (r Arguments) Has(items ...string) bool {
	for _, item := range items {
		if !StringInSlice(item, r) {
			return false
		}
	}

	return true
}

func (r Arguments) HasOneOf(items ...string) bool {
	for _, item := range items {
		if StringInSlice(item, r) {
			return true
		}
	}

	return false
}

func (r Arguments) Exact(name string) bool {
	return name == strings.Join(r, " ")
}
