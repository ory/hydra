// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build hsm
// +build hsm

package hsm

import (
	"crypto/elliptic"

	"github.com/ThalesGroup/crypto11"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/x/logrusx"
)

type Context interface {
	GenerateRSAKeyPairWithAttributes(public, private crypto11.AttributeSet, bits int) (crypto11.SignerDecrypter, error)
	GenerateECDSAKeyPairWithAttributes(public, private crypto11.AttributeSet, curve elliptic.Curve) (crypto11.Signer, error)
	FindKeyPair(id []byte, label []byte) (crypto11.Signer, error)
	FindKeyPairs(id []byte, label []byte) (signer []crypto11.Signer, err error)
	GetAttribute(key interface{}, attribute crypto11.AttributeType) (a *crypto11.Attribute, err error)
}

func NewContext(c *config.DefaultProvider, l *logrusx.Logger) Context {
	config11 := &crypto11.Config{
		Path: c.HSMLibraryPath(),
		Pin:  c.HSMPin(),
	}

	if c.HSMTokenLabel() != "" {
		config11.TokenLabel = c.HSMTokenLabel()
	} else {
		config11.SlotNumber = c.HSMSlotNumber()
	}

	ctx11, err := crypto11.Configure(config11)
	if err != nil {
		l.WithError(err).Fatalf("Unable to configure Hardware Security Module. Library path: %s, slot: %v, token label: %s",
			c.HSMLibraryPath(), *c.HSMSlotNumber(), c.HSMTokenLabel())
	} else {
		l.Info("Hardware Security Module is configured.")
	}

	var hsmContext Context = ctx11
	return hsmContext
}
