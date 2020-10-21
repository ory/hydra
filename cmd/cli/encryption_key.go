package cli

import (
	"github.com/pkg/errors"
	"github.com/sawadashota/encrypta"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"net/http"
)

const (
	FlagPGPKey    = "pgp-key"
	FlagPGPKeyURL = "pgp-key-url"
	FlagKeybase   = "keybase"
)

func noEncrypt(secret string) (string, error) {
	return secret, nil
}

func encryptWithKey(key encrypta.EncryptionKey) func(string) (string, error) {
	return func(secret string) (string, error) {
		enc, err := key.Encrypt([]byte(secret))
		if err != nil {
			return "", errors.WithStack(err)
		}
		return enc.Base64Encode(), nil
	}
}

// NewEncryptionFunc for client secret
func NewEncryptionFunc(cmd *cobra.Command, client *http.Client) (func(string) (string, error), error) {
	if client == nil {
		client = http.DefaultClient
	}

	pgpKey, err := cmd.Flags().GetString(FlagPGPKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	pgpKeyURL, err := cmd.Flags().GetString(FlagPGPKeyURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	keybaseUsername, err := cmd.Flags().GetString(FlagKeybase)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if pgpKey != "" {
		ek, err := encrypta.NewPublicKeyFromBase64Encoded(pgpKey)
		return encryptWithKey(ek), errors.WithStack(err)
	}

	if pgpKeyURL != "" {
		ek, err := encrypta.NewPublicKeyFromURL(pgpKeyURL, encrypta.HTTPClientOption(client))
		return encryptWithKey(ek), errors.WithStack(err)
	}

	if keybaseUsername != "" {
		ek, err := encrypta.NewPublicKeyFromKeybase(keybaseUsername, encrypta.HTTPClientOption(client))
		return encryptWithKey(ek), errors.WithStack(err)
	}

	return noEncrypt, nil
}

func RegisterSecretEncryptionFlags(flags *pflag.FlagSet) {
	flags.String(FlagPGPKey, "", "Base64 encoded PGP encryption key for encrypting client secret")
	flags.String(FlagPGPKeyURL, "", "PGP encryption key URL for encrypting client secret")
	flags.String(FlagKeybase, "", "Keybase username for encrypting client secret")
}
