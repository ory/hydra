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

package cli

import (
	"net/http"

	"github.com/sawadashota/encrypta"
	"github.com/spf13/cobra"

	"github.com/ory/x/flagx"
)

const (
	FlagEncryptionPGPKey    = "pgp-key"
	FlagEncryptionPGPKeyURL = "pgp-key-url"
	FlagEncryptionKeybase   = "keybase"
)

// NewEncryptionKey for client secret
func NewEncryptionKey(cmd *cobra.Command, client *http.Client) (ek encrypta.EncryptionKey, encryptSecret bool, err error) {
	if client == nil {
		client = http.DefaultClient
	}

	pgpKey := flagx.MustGetString(cmd, FlagEncryptionPGPKey)
	pgpKeyURL := flagx.MustGetString(cmd, FlagEncryptionPGPKeyURL)
	keybaseUsername := flagx.MustGetString(cmd, FlagEncryptionKeybase)

	if pgpKey != "" {
		ek, err = encrypta.NewPublicKeyFromBase64Encoded(pgpKey)
		encryptSecret = true
		return
	}

	if pgpKeyURL != "" {
		ek, err = encrypta.NewPublicKeyFromURL(pgpKeyURL, encrypta.HTTPClientOption(client))
		encryptSecret = true
		return
	}

	if keybaseUsername != "" {
		ek, err = encrypta.NewPublicKeyFromKeybase(keybaseUsername, encrypta.HTTPClientOption(client))
		encryptSecret = true
		return
	}

	return nil, false, nil
}
