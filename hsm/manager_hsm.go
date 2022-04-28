//go:build hsm
// +build hsm

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

	"github.com/ory/hydra/driver/config"

	"github.com/pkg/errors"

	"github.com/pborman/uuid"

	"github.com/ory/fosite"
	"github.com/ory/hydra/jwk"

	"github.com/miekg/pkcs11"

	"github.com/ory/hydra/x"

	"github.com/ThalesIgnite/crypto11"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/cryptosigner"
)

type KeyManager struct {
	jwk.Manager
	sync.RWMutex
	Context
	KeySetPrefix string
}

var ErrPreGeneratedKeys = &fosite.RFC6749Error{
	CodeField:        http.StatusBadRequest,
	ErrorField:       http.StatusText(http.StatusBadRequest),
	DescriptionField: "Cannot add/update pre generated keys on Hardware Security Module",
}

func NewKeyManager(hsm Context, config *config.Provider) *KeyManager {
	return &KeyManager{
		Context:      hsm,
		KeySetPrefix: config.HsmKeySetPrefix(),
	}
}

func (m *KeyManager) GenerateAndPersistKeySet(_ context.Context, set, kid, alg, use string) (*jose.JSONWebKeySet, error) {
	m.Lock()
	defer m.Unlock()

	set = m.prefixKeySet(set)

	err := m.deleteExistingKeySet(set)
	if err != nil {
		return nil, err
	}

	if len(kid) == 0 {
		kid = uuid.New()
	}

	privateAttrSet, publicAttrSet, err := getKeyPairAttributes(kid, set, use)
	if err != nil {
		return nil, err
	}

	switch {
	case alg == "RS256":
		key, err := m.GenerateRSAKeyPairWithAttributes(publicAttrSet, privateAttrSet, 4096)
		if err != nil {
			return nil, err
		}
		return createKeySet(key, kid, alg, use)
	case alg == "ES256":
		key, err := m.GenerateECDSAKeyPairWithAttributes(publicAttrSet, privateAttrSet, elliptic.P256())
		if err != nil {
			return nil, err
		}
		return createKeySet(key, kid, alg, use)
	case alg == "ES512":
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

func (m *KeyManager) GetKey(_ context.Context, set, kid string) (*jose.JSONWebKeySet, error) {
	m.RLock()
	defer m.RUnlock()

	set = m.prefixKeySet(set)

	keyPair, err := m.FindKeyPair([]byte(kid), []byte(set))
	if err != nil {
		return nil, err
	}

	if keyPair == nil {
		return nil, errors.WithStack(x.ErrNotFound)
	}

	id, alg, use, err := getKeySetAttributes(m, keyPair, []byte(kid))
	if err != nil {
		return nil, err
	}

	return createKeySet(keyPair, id, alg, use)
}

func (m *KeyManager) GetKeySet(_ context.Context, set string) (*jose.JSONWebKeySet, error) {
	m.RLock()
	defer m.RUnlock()

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
		kid, alg, use, err := getKeySetAttributes(m, keyPair, nil)
		if err != nil {
			return nil, err
		}
		keys = append(keys, createKeys(keyPair, kid, alg, use)...)
	}

	return &jose.JSONWebKeySet{
		Keys: keys,
	}, nil
}

func (m *KeyManager) DeleteKey(_ context.Context, set, kid string) error {
	m.Lock()
	defer m.Unlock()

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

func (m *KeyManager) DeleteKeySet(_ context.Context, set string) error {
	m.Lock()
	defer m.Unlock()

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

func getKeySetAttributes(m *KeyManager, key crypto11.Signer, kid []byte) (string, string, string, error) {
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
		// TODO Should we validate minimal key length by checking CKA_MODULUS_BITS?
		// TODO see https://github.com/ory/hydra/issues/2905
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

func getKeyPairAttributes(kid string, set string, use string) (crypto11.AttributeSet, crypto11.AttributeSet, error) {

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
	}, {
		Algorithm:                   alg,
		Use:                         use,
		Key:                         key.Public(),
		KeyID:                       kid,
		Certificates:                []*x509.Certificate{},
		CertificateThumbprintSHA1:   []uint8{},
		CertificateThumbprintSHA256: []uint8{},
	}}
}

func (m *KeyManager) prefixKeySet(set string) string {
	return fmt.Sprintf("%s%s", m.KeySetPrefix, set)
}
