package jwk

import (
	"context"

	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"

	"github.com/ory/hydra/x"
)

type ManagerStrategy struct {
	hardwareKeyManager Manager
	softwareKeyManager Manager
}

func NewManagerStrategy(hardwareKeyManager Manager, softwareKeyManager Manager) *ManagerStrategy {
	return &ManagerStrategy{
		hardwareKeyManager: hardwareKeyManager,
		softwareKeyManager: softwareKeyManager,
	}
}

func (m ManagerStrategy) GenerateAndPersistKeySet(ctx context.Context, set, kid, alg, use string) (*jose.JSONWebKeySet, error) {
	return m.hardwareKeyManager.GenerateAndPersistKeySet(ctx, set, kid, alg, use)
}

func (m ManagerStrategy) AddKey(ctx context.Context, set string, key *jose.JSONWebKey) error {
	return m.softwareKeyManager.AddKey(ctx, set, key)
}

func (m ManagerStrategy) AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) error {
	return m.softwareKeyManager.AddKeySet(ctx, set, keys)
}

func (m ManagerStrategy) UpdateKey(ctx context.Context, set string, key *jose.JSONWebKey) error {
	return m.softwareKeyManager.UpdateKey(ctx, set, key)
}

func (m ManagerStrategy) UpdateKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) error {
	return m.softwareKeyManager.UpdateKeySet(ctx, set, keys)
}

func (m ManagerStrategy) GetKey(ctx context.Context, set, kid string) (*jose.JSONWebKeySet, error) {
	keySet, err := m.hardwareKeyManager.GetKey(ctx, set, kid)
	if err != nil && !errors.Is(err, x.ErrNotFound) {
		return nil, err
	} else if keySet != nil {
		return keySet, nil
	} else {
		return m.softwareKeyManager.GetKey(ctx, set, kid)
	}
}

func (m ManagerStrategy) GetKeySet(ctx context.Context, set string) (*jose.JSONWebKeySet, error) {
	keySet, err := m.hardwareKeyManager.GetKeySet(ctx, set)
	if err != nil && !errors.Is(err, x.ErrNotFound) {
		return nil, err
	} else if keySet != nil {
		return keySet, nil
	} else {
		return m.softwareKeyManager.GetKeySet(ctx, set)
	}
}

func (m ManagerStrategy) DeleteKey(ctx context.Context, set, kid string) error {
	err := m.hardwareKeyManager.DeleteKey(ctx, set, kid)
	if err != nil && !errors.Is(err, x.ErrNotFound) {
		return err
	} else if errors.Is(err, x.ErrNotFound) {
		return m.softwareKeyManager.DeleteKey(ctx, set, kid)
	} else {
		return nil
	}
}

func (m ManagerStrategy) DeleteKeySet(ctx context.Context, set string) error {
	err := m.hardwareKeyManager.DeleteKeySet(ctx, set)
	if err != nil && !errors.Is(err, x.ErrNotFound) {
		return err
	} else if errors.Is(err, x.ErrNotFound) {
		return m.softwareKeyManager.DeleteKeySet(ctx, set)
	} else {
		return nil
	}
}
