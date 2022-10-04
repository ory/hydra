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

package jwk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"

	"github.com/stretchr/testify/assert"
)

func TestMustRSAPrivate(t *testing.T) {
	keys, err := GenerateJWK(context.Background(), jose.RS256, "foo", "sig")
	require.NoError(t, err)

	priv := keys.Key("foo")[0]
	_, err = ToRSAPrivate(&priv)
	assert.Nil(t, err)

	MustRSAPrivate(&priv)

	pub := keys.Key("foo")[0].Public()
	_, err = ToRSAPublic(&pub)
	assert.Nil(t, err)
	MustRSAPublic(&pub)
}
