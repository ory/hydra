package sql

import (
	"context"
	"encoding/json"

	"github.com/gobuffalo/pop/v6"
	"gopkg.in/square/go-jose.v2"

	"github.com/ory/x/errorsx"

	"github.com/pkg/errors"

	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon"
)

var _ jwk.Manager = &Persister{}

func (p *Persister) GenerateAndPersistKeySet(ctx context.Context, set, kid, alg, use string) (*jose.JSONWebKeySet, error) {
	generator, found := p.r.KeyGenerators()[alg]
	if !found {
		return nil, errorsx.WithStack(jwk.ErrUnsupportedKeyAlgorithm)
	}

	keys, err := generator.Generate(kid, use)
	if err != nil {
		return nil, err
	}

	err = p.AddKeySet(ctx, set, keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (p *Persister) AddKey(ctx context.Context, set string, key *jose.JSONWebKey) error {
	out, err := json.Marshal(key)
	if err != nil {
		return errorsx.WithStack(err)
	}

	encrypted, err := p.r.KeyCipher().Encrypt(out)
	if err != nil {
		return errorsx.WithStack(err)
	}

	return sqlcon.HandleError(p.Connection(ctx).Create(&jwk.SQLData{
		Set:     set,
		KID:     key.KeyID,
		Version: 0,
		Key:     encrypted,
	}))
}

func (p *Persister) AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) error {
	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		for _, key := range keys.Keys {
			out, err := json.Marshal(key)
			if err != nil {
				return errorsx.WithStack(err)
			}

			encrypted, err := p.r.KeyCipher().Encrypt(out)
			if err != nil {
				return err
			}

			if err := c.Create(&jwk.SQLData{
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

func (p *Persister) UpdateKey(ctx context.Context, set string, key *jose.JSONWebKey) error {
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

func (p *Persister) UpdateKeySet(ctx context.Context, set string, keySet *jose.JSONWebKeySet) error {
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
	var j jwk.SQLData
	if err := p.Connection(ctx).
		Where("sid = ? AND kid = ?", set, kid).
		Order("created_at DESC").
		First(&j); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	key, err := p.r.KeyCipher().Decrypt(j.Key)
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

func (p *Persister) GetKeySet(ctx context.Context, set string) (*jose.JSONWebKeySet, error) {
	var js []jwk.SQLData
	if err := p.Connection(ctx).
		Where("sid = ?", set).
		Order("created_at DESC").
		All(&js); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if len(js) == 0 {
		return nil, errors.Wrap(x.ErrNotFound, "")
	}

	keys := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
	for _, d := range js {
		key, err := p.r.KeyCipher().Decrypt(d.Key)
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
	return sqlcon.HandleError(p.Connection(ctx).RawQuery("DELETE FROM hydra_jwk WHERE sid=? AND kid=?", set, kid).Exec())
}

func (p *Persister) DeleteKeySet(ctx context.Context, set string) error {
	return sqlcon.HandleError(p.Connection(ctx).RawQuery("DELETE FROM hydra_jwk WHERE sid=?", set).Exec())
}
