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

package pkg

import (
	"time"

	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/storage"
	"github.com/ory/fosite/token/hmac"
)

var HMACStrategy = &oauth2.HMACSHAStrategy{
	Enigma: &hmac.HMACStrategy{
		GlobalSecret: []byte("1234567890123456789012345678901234567890"),
	},
	AccessTokenLifespan:   time.Hour,
	AuthorizeCodeLifespan: time.Hour,
}

func FositeStore() *storage.MemoryStore {
	return storage.NewMemoryStore()
}

func Tokens(length int) (res [][]string) {
	for i := 0; i < length; i++ {
		tok, sig, _ := HMACStrategy.Enigma.Generate()
		res = append(res, []string{sig, tok})
	}
	return res
}
