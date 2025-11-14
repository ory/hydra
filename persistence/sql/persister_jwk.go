// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"encoding/json"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ory/hydra/v2/aead"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/pop/v6"
	"github.com/ory/x/otelx"
	"github.com/ory/x/sqlcon"
)

var _ jwk.Manager = (*JWKPersister)(nil)

type JWKPersister struct {
	*BasePersister
}

// GenerateAndPersistKeySet implements jwk.Manager.
func (p *JWKPersister) GenerateAndPersistKeySet(ctx context.Context, set, kid, alg, use string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GenerateAndPersistKeySet",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", kid),
			attribute.String("alg", alg)))
	defer otelx.End(span, &err)

	if kid == "" {
		kid = uuid.Must(uuid.NewV4()).String()
	}

	keys, err := jwk.GenerateJWK(jose.SignatureAlgorithm(alg), kid, use)
	if err != nil {
		return nil, errors.Wrapf(jwk.ErrUnsupportedKeyAlgorithm, "%s", err)
	}

	err = p.AddKeySet(ctx, set, keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

// AddKey implements jwk.Manager.
func (p *JWKPersister) AddKey(ctx context.Context, set string, key *jose.JSONWebKey) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.AddKey",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", key.KeyID)))
	defer otelx.End(span, &err)

	out, err := json.Marshal(key)
	if err != nil {
		return errors.WithStack(err)
	}

	encrypted, err := aead.NewAESGCM(p.d.Config()).Encrypt(ctx, out, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	return sqlcon.HandleError(p.CreateWithNetwork(ctx, &jwk.SQLData{
		Set:     set,
		KID:     key.KeyID,
		Version: 0,
		Key:     encrypted,
	}))
}

// AddKeySet implements jwk.Manager.
func (p *JWKPersister) AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.AddKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		for _, key := range keys.Keys {
			out, err := json.Marshal(key)
			if err != nil {
				return errors.WithStack(err)
			}

			encrypted, err := aead.NewAESGCM(p.d.Config()).Encrypt(ctx, out, nil)
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

// UpdateKey updates or creates the key. Implements jwk.Manager.
func (p *JWKPersister) UpdateKey(ctx context.Context, set string, key *jose.JSONWebKey) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.UpdateKey",
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

// UpdateKeySet updates or creates the key set. Implements jwk.Manager.
func (p *JWKPersister) UpdateKeySet(ctx context.Context, set string, keySet *jose.JSONWebKeySet) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.UpdateKeySet", trace.WithAttributes(attribute.String("set", set)))
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

// GetKey implements jwk.Manager.
func (p *JWKPersister) GetKey(ctx context.Context, set, kid string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetKey",
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

	key, err := aead.NewAESGCM(p.d.Config()).Decrypt(ctx, j.Key, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var c jose.JSONWebKey
	if err := json.Unmarshal(key, &c); err != nil {
		return nil, errors.WithStack(err)
	}

	return &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{c},
	}, nil
}

// GetKeySet implements jwk.Manager.
func (p *JWKPersister) GetKeySet(ctx context.Context, set string) (keys *jose.JSONWebKeySet, err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	var js jwk.SQLDataRows
	if err := p.QueryWithNetwork(ctx).
		Where("sid = ?", set).
		Order("created_at DESC").
		All(&js); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return js.ToJWK(ctx, aead.NewAESGCM(p.d.Config()))
}

// DeleteKey implements jwk.Manager.
func (p *JWKPersister) DeleteKey(ctx context.Context, set, kid string) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteKey",
		trace.WithAttributes(
			attribute.String("set", set),
			attribute.String("kid", kid)))
	defer otelx.End(span, &err)

	err = p.QueryWithNetwork(ctx).Where("sid=? AND kid=?", set, kid).Delete(&jwk.SQLData{})
	return sqlcon.HandleError(err)
}

// DeleteKeySet implements jwk.Manager.
func (p *JWKPersister) DeleteKeySet(ctx context.Context, set string) (err error) {
	ctx, span := p.d.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteKeySet", trace.WithAttributes(attribute.String("set", set)))
	defer otelx.End(span, &err)

	err = p.QueryWithNetwork(ctx).Where("sid=?", set).Delete(&jwk.SQLData{})
	return sqlcon.HandleError(err)
}
