// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"encoding/json"

	"github.com/go-jose/go-jose/v3"
	"github.com/gobuffalo/pop/v6"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/sqlcon"
)

var _ jwk.Manager = &Persister{}

func (p *Persister) GenerateAndPersistKeySet(ctx context.Context, set, kid, alg, use string) (*jose.JSONWebKeySet, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GenerateAndPersistKey")
	defer span.End()

	keys, err := jwk.GenerateJWK(ctx, jose.SignatureAlgorithm(alg), kid, use)
	if err != nil {
		return nil, errors.Wrapf(jwk.ErrUnsupportedKeyAlgorithm, "%s", err)
	}

	err = p.AddKeySet(ctx, set, keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (p *Persister) AddKey(ctx context.Context, set string, key *jose.JSONWebKey) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.AddKey")
	defer span.End()

	out, err := json.Marshal(key)
	if err != nil {
		return errorsx.WithStack(err)
	}

	encrypted, err := p.r.KeyCipher().Encrypt(ctx, out, nil)
	if err != nil {
		return errorsx.WithStack(err)
	}

	return sqlcon.HandleError(p.CreateWithNetwork(ctx, &jwk.SQLData{
		Set:     set,
		KID:     key.KeyID,
		Version: 0,
		Key:     encrypted,
	}))
}

func (p *Persister) AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.AddKey")
	defer span.End()

	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		for _, key := range keys.Keys {
			out, err := json.Marshal(key)
			if err != nil {
				return errorsx.WithStack(err)
			}

			encrypted, err := p.r.KeyCipher().Encrypt(ctx, out, nil)
			if err != nil {
				return err
			}

			if err := p.CreateWithNetwork(ctx, &jwk.SQLData{
				Set:     set,
				KID:     key.KeyID,
				Version: 0,
				Key:     encrypted,
			}); err != nil {
				return sqlcon.HandleError(err)
			}
		}
		return nil
	})
}

// UpdateKey updates or creates the key.
func (p *Persister) UpdateKey(ctx context.Context, set string, key *jose.JSONWebKey) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.UpdateKey")
	defer span.End()

	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := p.DeleteKey(ctx, set, key.KeyID); err != nil {
			return err
		}
		if err := p.AddKey(ctx, set, key); err != nil {
			return err
		}
		return nil
	})
}

// UpdateKeySet updates or creates the key set.
func (p *Persister) UpdateKeySet(ctx context.Context, set string, keySet *jose.JSONWebKeySet) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.UpdateKeySet")
	defer span.End()

	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := p.DeleteKeySet(ctx, set); err != nil {
			return err
		}
		if err := p.AddKeySet(ctx, set, keySet); err != nil {
			return err
		}
		return nil
	})
}

func (p *Persister) GetKey(ctx context.Context, set, kid string) (*jose.JSONWebKeySet, error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetKey")
	defer span.End()

	var j jwk.SQLData
	if err := p.QueryWithNetwork(ctx).
		Where("sid = ? AND kid = ?", set, kid).
		Order("created_at DESC").
		First(&j); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	key, err := p.r.KeyCipher().Decrypt(ctx, j.Key, nil)
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	var c jose.JSONWebKey
	if err := json.Unmarshal(key, &c); err != nil {
		return nil, errorsx.WithStack(err)
	}

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{c},
	}, nil
}

func (p *Persister) GetKeySet(ctx context.Context, set string) (keys *jose.JSONWebKeySet, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetKeySet")
	defer span.End()

	var js []jwk.SQLData
	if err := p.QueryWithNetwork(ctx).
		Where("sid = ?", set).
		Order("created_at DESC").
		All(&js); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if len(js) == 0 {
		return nil, errors.Wrap(x.ErrNotFound, "")
	}

	keys = &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
	for _, d := range js {
		key, err := p.r.KeyCipher().Decrypt(ctx, d.Key, nil)
		if err != nil {
			return nil, errorsx.WithStack(err)
		}

		var c jose.JSONWebKey
		if err := json.Unmarshal(key, &c); err != nil {
			return nil, errorsx.WithStack(err)
		}
		keys.Keys = append(keys.Keys, c)
	}

	if len(keys.Keys) == 0 {
		return nil, errorsx.WithStack(x.ErrNotFound)
	}

	return keys, nil
}

func (p *Persister) DeleteKey(ctx context.Context, set, kid string) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteKey")
	defer span.End()

	err := p.QueryWithNetwork(ctx).Where("sid=? AND kid=?", set, kid).Delete(&jwk.SQLData{})
	return sqlcon.HandleError(err)
}

func (p *Persister) DeleteKeySet(ctx context.Context, set string) error {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteKeySet")
	defer span.End()

	err := p.QueryWithNetwork(ctx).Where("sid=?", set).Delete(&jwk.SQLData{})
	return sqlcon.HandleError(err)
}
