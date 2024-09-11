// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"encoding/json"

	"github.com/go-jose/go-jose/v3"
	"github.com/gobuffalo/pop/v6"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ory/x/errorsx"
	"github.com/ory/x/otelx"

	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/x/sqlcon"
)

var _ jwk.Manager = &Persister{}

func (p *Persister) GenerateAndPersistKeySet(ctx context.Context, set, kid, alg, use string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GenerateAndPersistKeySet",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", kid),
			attribute.String("alg", alg)))
	defer otelx.End(span, &err)

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

func (p *Persister) AddKey(ctx context.Context, set string, key *jose.JSONWebKey) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.AddKey",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", key.KeyID)))

	defer otelx.End(span, &err)

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

func (p *Persister) AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.AddKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
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
func (p *Persister) UpdateKey(ctx context.Context, set string, key *jose.JSONWebKey) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.UpdateKey",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", key.KeyID)))
	defer otelx.End(span, &err)

	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
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
func (p *Persister) UpdateKeySet(ctx context.Context, set string, keySet *jose.JSONWebKeySet) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.UpdateKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := p.DeleteKeySet(ctx, set); err != nil {
			return err
		}
		if err := p.AddKeySet(ctx, set, keySet); err != nil {
			return err
		}
		return nil
	})
}

func (p *Persister) GetKey(ctx context.Context, set, kid string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetKey",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", kid)))
	defer otelx.End(span, &err)

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
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	var js jwk.SQLDataRows
	if err := p.QueryWithNetwork(ctx).
		Where("sid = ?", set).
		Order("created_at DESC").
		All(&js); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return js.ToJWK(ctx, p.r)
}

func (p *Persister) DeleteKey(ctx context.Context, set, kid string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteKey",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", kid)))
	defer otelx.End(span, &err)

	err = p.QueryWithNetwork(ctx).Where("sid=? AND kid=?", set, kid).Delete(&jwk.SQLData{})
	return sqlcon.HandleError(err)
}

func (p *Persister) DeleteKeySet(ctx context.Context, set string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	err = p.QueryWithNetwork(ctx).Where("sid=?", set).Delete(&jwk.SQLData{})
	return sqlcon.HandleError(err)
}
