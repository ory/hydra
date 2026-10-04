// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build hsm

package hsm

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net/http"
	"sync"

	"github.com/ThalesGroup/crypto11"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/cryptosigner"
	"github.com/gofrs/uuid"
	"github.com/miekg/pkcs11"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/otelx"
)

const tracingComponent = "github.com/ory/hydra/hsm"

type KeyManager struct {
	jwk.Manager
	Context
	c    config.DefaultProvider
	lock *sync.RWMutex
}

var ErrPreGeneratedKeys = &fosite.RFC6749Error{
	CodeField:        http.StatusBadRequest,
	ErrorField:       http.StatusText(http.StatusBadRequest),
	DescriptionField: "Cannot add/update pre generated keys on Hardware Security Module",
}

// tokenLocks serializes mutating PKCS#11 operations per HSM token within a
// process. Several KeyManager instances can target the same token (for example
// one registry per test). crypto11 gives each its own session pool, but SoftHSM
// is not reliably thread-safe for concurrent key generation against a single
// token. The lock is keyed by token identity so independent tokens never block
// each other. The map is never pruned, but it holds one entry per distinct
// token configured in this process — at most a handful — so it does not grow
// with load.
var tokenLocks sync.Map // tokenKey -> *sync.RWMutex

// lockForToken returns the process-wide lock shared by all KeyManagers that
// target the same HSM token.
func lockForToken(c *config.DefaultProvider) *sync.RWMutex {
	l, _ := tokenLocks.LoadOrStore(tokenKey(c), &sync.RWMutex{})
	return l.(*sync.RWMutex)
}

// tokenKey identifies a token by the field crypto11 actually uses to select it:
// the token label when set, otherwise the slot number. NewContext applies the
// same precedence, so two configs that resolve to the same token always share a
// lock even if the unused selector differs.
func tokenKey(c *config.DefaultProvider) string {
	if label := c.HSMTokenLabel(); label != "" {
		return fmt.Sprintf("label\x00%s\x00%s", c.HSMLibraryPath(), label)
	}
	slot := -1
	if s := c.HSMSlotNumber(); s != nil {
		slot = *s
	}
	return fmt.Sprintf("slot\x00%s\x00%d", c.HSMLibraryPath(), slot)
}

func NewKeyManager(hsm Context, config *config.DefaultProvider) *KeyManager {
	return &KeyManager{
		Context: hsm,
		c:       *config,
		lock:    lockForToken(config),
	}
}

func (m *KeyManager) GenerateAndPersistKeySet(ctx context.Context, set, kid, alg, use string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "hsm.GenerateAndPersistKeySet",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", kid),
			attribute.String("alg", alg),
			attribute.String("use", use)))
	defer otelx.End(span, &err)

	m.lock.Lock()
	defer m.lock.Unlock()

	set = m.prefixKeySet(set)

	err = m.deleteExistingKeySet(set)
	if err != nil {
		return nil, err
	}

	if kid == "" {
		kid = uuid.Must(uuid.NewV4()).String()
	}

	privateAttrSet, publicAttrSet, err := getKeyPairAttributes(kid, set, use)
	if err != nil {
		return nil, err
	}

	switch alg {
	case "RS256":
		key, err := m.GenerateRSAKeyPairWithAttributes(publicAttrSet, privateAttrSet, 4096)
		if err != nil {
			return nil, err
		}
		return createKeySet(key, kid, alg, use)
	case "ES256":
		key, err := m.GenerateECDSAKeyPairWithAttributes(publicAttrSet, privateAttrSet, elliptic.P256())
		if err != nil {
			return nil, err
		}
		return createKeySet(key, kid, alg, use)
	case "ES512":
		key, err := m.GenerateECDSAKeyPairWithAttributes(publicAttrSet, privateAttrSet, elliptic.P521())
		if err != nil {
			return nil, err
		}
		return createKeySet(key, kid, alg, use)

	// NOTE:
	//	- HS256, HS512 not supported. Makes sense only if shared HSM is used between Hydra and authenticating client.
	//	- EdDSA not supported. As of now PKCS#11 v2.4 doesn't support EdDSA keys using curve Ed25519. However,
	//	  PKCS#11 3.0 (https://docs.oasis-open.org/pkcs11/pkcs11-curr/v3.0/pkcs11-curr-v3.0.html)
	//	  contains support for EdDSA.

	default:
		return nil, errors.WithStack(jwk.ErrUnsupportedKeyAlgorithm)
	}
}

func (m *KeyManager) GetKey(ctx context.Context, set, kid string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "hsm.GetKey",
		trace.WithAttributes(attribute.String("set", set), attribute.String("kid", kid)))
	defer otelx.End(span, &err)

	m.lock.RLock()
	defer m.lock.RUnlock()

	set = m.prefixKeySet(set)

	keyPair, err := m.FindKeyPair([]byte(kid), []byte(set))
	if err != nil {
		return nil, err
	}

	if keyPair == nil {
		return nil, errors.WithStack(x.ErrNotFound)
	}

	id, alg, use, err := m.getKeySetAttributes(ctx, keyPair, []byte(kid))
	if err != nil {
		return nil, err
	}

	return createKeySet(keyPair, id, alg, use)
}

func (m *KeyManager) GetKeySet(ctx context.Context, set string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "hsm.GetKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	m.lock.RLock()
	defer m.lock.RUnlock()

	set = m.prefixKeySet(set)

	keyPairs, err := m.FindKeyPairs(nil, []byte(set))
	if err != nil {
		return nil, err
	}

	if keyPairs == nil {
		return nil, errors.WithStack(x.ErrNotFound)
	}

	var keys []jose.JSONWebKey
	for _, keyPair := range keyPairs {
		kid, alg, use, err := m.getKeySetAttributes(ctx, keyPair, nil)
		if err != nil {
			return nil, err
		}
		keys = append(keys, createKeys(keyPair, kid, alg, use)...)
	}

	return &jose.JSONWebKeySet{
		Keys: keys,
	}, nil
}

func (m *KeyManager) DeleteKey(ctx context.Context, set, kid string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "hsm.DeleteKey",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", kid)))
	defer otelx.End(span, &err)

	m.lock.Lock()
	defer m.lock.Unlock()

	set = m.prefixKeySet(set)

	keyPair, err := m.FindKeyPair([]byte(kid), []byte(set))
	if err != nil {
		return err
	}

	if keyPair != nil {
		err = keyPair.Delete()
		if err != nil {
			return err
		}
	} else {
		return errors.WithStack(x.ErrNotFound)
	}
	return nil
}

func (m *KeyManager) DeleteKeySet(ctx context.Context, set string) (err error) {
	ctx, span := otel.GetTracerProvider().Tracer(tracingComponent).Start(ctx, "hsm.DeleteKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	m.lock.Lock()
	defer m.lock.Unlock()

	set = m.prefixKeySet(set)

	keyPairs, err := m.FindKeyPairs(nil, []byte(set))
	if err != nil {
		return err
	}

	if keyPairs == nil {
		return errors.WithStack(x.ErrNotFound)
	}

	for _, keyPair := range keyPairs {
		err = keyPair.Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *KeyManager) AddKey(_ context.Context, _ string, _ *jose.JSONWebKey) error {
	return errors.WithStack(ErrPreGeneratedKeys)
}

func (m *KeyManager) AddKeySet(_ context.Context, _ string, _ *jose.JSONWebKeySet) error {
	return errors.WithStack(ErrPreGeneratedKeys)
}

func (m *KeyManager) UpdateKey(_ context.Context, _ string, _ *jose.JSONWebKey) error {
	return errors.WithStack(ErrPreGeneratedKeys)
}

func (m *KeyManager) UpdateKeySet(_ context.Context, _ string, _ *jose.JSONWebKeySet) error {
	return errors.WithStack(ErrPreGeneratedKeys)
}

func (m *KeyManager) getKeySetAttributes(ctx context.Context, key crypto11.Signer, kid []byte) (string, string, string, error) {
	if kid == nil {
		ckaId, err := m.GetAttribute(key, crypto11.CkaId)
		if err != nil {
			return "", "", "", err
		}
		kid = ckaId.Value
	}

	var alg string
	switch k := key.Public().(type) {
	case *rsa.PublicKey:
		alg = "RS256"
		if k.N.BitLen() < 4096 && !m.c.IsDevelopmentMode(ctx) {
			return "", "", "", errors.WithStack(jwk.ErrMinimalRsaKeyLength)
		}
	case *ecdsa.PublicKey:
		if k.Curve == elliptic.P521() {
			alg = "ES512"
		} else if k.Curve == elliptic.P256() {
			alg = "ES256"
		} else {
			return "", "", "", errors.WithStack(jwk.ErrUnsupportedEllipticCurve)
		}
	default:
		return "", "", "", errors.WithStack(jwk.ErrUnsupportedKeyAlgorithm)
	}

	use := "sig"
	ckaDecrypt, _ := m.GetAttribute(key, crypto11.CkaDecrypt)
	if ckaDecrypt != nil && len(ckaDecrypt.Value) != 0 && ckaDecrypt.Value[0] == 0x1 {
		use = "enc"
	}
	return string(kid), alg, use, nil
}

func getKeyPairAttributes(kid, set, use string) (crypto11.AttributeSet, crypto11.AttributeSet, error) {
	privateAttrSet, err := crypto11.NewAttributeSetWithIDAndLabel([]byte(kid), []byte(set))
	if err != nil {
		return nil, nil, err
	}

	publicAttrSet, err := crypto11.NewAttributeSetWithIDAndLabel([]byte(kid), []byte(set))
	if err != nil {
		return nil, nil, err
	}

	if len(use) == 0 || use == "sig" {
		publicAttrSet.AddIfNotPresent([]*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
			pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, false),
		})
		privateAttrSet.AddIfNotPresent([]*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
			pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, false),
		})
	} else {
		publicAttrSet.AddIfNotPresent([]*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_VERIFY, false),
			pkcs11.NewAttribute(pkcs11.CKA_ENCRYPT, true),
		})
		privateAttrSet.AddIfNotPresent([]*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_SIGN, false),
			pkcs11.NewAttribute(pkcs11.CKA_DECRYPT, true),
		})
	}

	return privateAttrSet, publicAttrSet, nil
}

func (m *KeyManager) deleteExistingKeySet(set string) error {
	existingKeyPairs, err := m.FindKeyPairs(nil, []byte(set))
	if err != nil {
		return err
	}
	if len(existingKeyPairs) != 0 {
		for _, keyPair := range existingKeyPairs {
			_ = keyPair.Delete()
		}
	}
	return nil
}

func createKeySet(key crypto11.Signer, kid, alg, use string) (*jose.JSONWebKeySet, error) {
	return &jose.JSONWebKeySet{
		Keys: createKeys(key, kid, alg, use),
	}, nil
}

func createKeys(key crypto11.Signer, kid, alg, use string) []jose.JSONWebKey {
	return []jose.JSONWebKey{{
		Algorithm:                   alg,
		Use:                         use,
		Key:                         cryptosigner.Opaque(key),
		KeyID:                       kid,
		Certificates:                []*x509.Certificate{},
		CertificateThumbprintSHA1:   []uint8{},
		CertificateThumbprintSHA256: []uint8{},
	}}
}

func (m *KeyManager) prefixKeySet(set string) string {
	return fmt.Sprintf("%s%s", m.c.HSMKeySetPrefix(), set)
}
