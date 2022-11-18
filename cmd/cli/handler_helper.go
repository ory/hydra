// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

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
