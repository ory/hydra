// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"strings"
	"time"

	"github.com/ory/hydra/v2/jwk"

	"github.com/pkg/errors"

	"github.com/go-jose/go-jose/v3"
	"github.com/gobuffalo/pop/v6"

	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/x/otelx"
	"github.com/ory/x/stringsx"

	"github.com/ory/x/sqlcon"
)

var _ trust.GrantManager = &Persister{}

func (p *Persister) CreateGrant(ctx context.Context, g trust.Grant, publicKey jose.JSONWebKey) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateGrant")
	defer otelx.End(span, &err)

	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		// add key, if it doesn't exist
		if _, err := p.GetKey(ctx, g.PublicKey.Set, g.PublicKey.KeyID); err != nil {
			if !errors.Is(err, sqlcon.ErrNoRows) {
				return sqlcon.HandleError(err)
			}

			if err = p.AddKey(ctx, g.PublicKey.Set, &publicKey); err != nil {
				return sqlcon.HandleError(err)
			}
		}

		data := p.sqlDataFromJWTGrant(g)
		return sqlcon.HandleError(p.CreateWithNetwork(ctx, &data))
	})
}

func (p *Persister) GetConcreteGrant(ctx context.Context, id string) (_ trust.Grant, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetConcreteGrant")
	defer otelx.End(span, &err)

	var data trust.SQLData
	if err := p.QueryWithNetwork(ctx).Where("id = ?", id).First(&data); err != nil {
		return trust.Grant{}, sqlcon.HandleError(err)
	}

	return p.jwtGrantFromSQlData(data), nil
}

func (p *Persister) DeleteGrant(ctx context.Context, id string) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteGrant")
	defer otelx.End(span, &err)

	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		grant, err := p.GetConcreteGrant(ctx, id)
		if err != nil {
			return sqlcon.HandleError(err)
		}

		if err := p.QueryWithNetwork(ctx).Where("id = ?", grant.ID).Delete(&trust.SQLData{}); err != nil {
			return sqlcon.HandleError(err)
		}

		return p.DeleteKey(ctx, grant.PublicKey.Set, grant.PublicKey.KeyID)
	})
}

func (p *Persister) GetGrants(ctx context.Context, limit, offset int, optionalIssuer string) (_ []trust.Grant, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetGrants")
	defer otelx.End(span, &err)

	grantsData := make([]trust.SQLData, 0)

	query := p.QueryWithNetwork(ctx).
		Paginate(offset/limit+1, limit).
		Order("id")
	if optionalIssuer != "" {
		query = query.Where("issuer = ?", optionalIssuer)
	}

	if err := query.All(&grantsData); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	grants := make([]trust.Grant, 0, len(grantsData))
	for _, data := range grantsData {
		grants = append(grants, p.jwtGrantFromSQlData(data))
	}

	return grants, nil
}

func (p *Persister) CountGrants(ctx context.Context) (n int, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CountGrants")
	defer otelx.End(span, &err)

	n, err = p.QueryWithNetwork(ctx).
		Count(&trust.SQLData{})
	return n, sqlcon.HandleError(err)
}

func (p *Persister) GetPublicKey(ctx context.Context, issuer string, subject string, keyId string) (_ *jose.JSONWebKey, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetPublicKey")
	defer otelx.End(span, &err)

	var data trust.SQLData
	query := p.QueryWithNetwork(ctx).
		Where("issuer = ?", issuer).
		Where("(subject = ? OR allow_any_subject IS TRUE)", subject).
		Where("key_id = ?", keyId).
		Where("nid = ?", p.NetworkID(ctx))
	if err := query.First(&data); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	keySet, err := p.GetKey(ctx, data.KeySet, keyId)
	if err != nil {
		return nil, err
	}

	return &keySet.Keys[0], nil
}

func (p *Persister) GetPublicKeys(ctx context.Context, issuer string, subject string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetPublicKeys")
	defer otelx.End(span, &err)

	expiresAt := "expires_at > NOW()"
	if p.conn.Dialect.Name() == "sqlite3" {
		expiresAt = "expires_at > datetime('now')"
	}

	grantsData := make([]trust.SQLData, 0)
	query := p.QueryWithNetwork(ctx).
		Select("key_id").
		Where(expiresAt).
		Where("issuer = ?", issuer).
		Where("(subject = ? OR allow_any_subject IS TRUE)", subject).
		Order("created_at DESC").
		Limit(100) // Load maximum of 100 keys

	if err := query.All(&grantsData); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if len(grantsData) == 0 {
		return &jose.JSONWebKeySet{}, nil
	}

	keyIDs := make([]interface{}, len(grantsData))
	for k, d := range grantsData {
		keyIDs[k] = d.KeyID
	}

	var js jwk.SQLDataRows
	if err := p.QueryWithNetwork(ctx).
		// key_set and issuer are set to the same value on creation:
		//
		//	grant := Grant{
		//		ID:              uuid.New().String(),
		//		Issuer:          grantRequest.Issuer,
		//		Subject:         grantRequest.Subject,
		//		AllowAnySubject: grantRequest.AllowAnySubject,
		//		Scope:           grantRequest.Scope,
		//		PublicKey: PublicKey{
		//			Set:   grantRequest.Issuer, // group all keys by issuer, so set=issuer
		//			KeyID: grantRequest.PublicKeyJWK.KeyID,
		//		},
		//		CreatedAt: time.Now().UTC().Round(time.Second),
		//		ExpiresAt: grantRequest.ExpiresAt.UTC().Round(time.Second),
		//	}
		//
		// Therefore it is fine if we only look for the issuer here instead of the key set id.
		Where("sid = ?", issuer).
		Where("kid IN (?)", keyIDs).
		Order("created_at DESC").
		All(&js); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return js.ToJWK(ctx, p.r)
}

func (p *Persister) GetPublicKeyScopes(ctx context.Context, issuer string, subject string, keyId string) (_ []string, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetPublicKeyScopes")
	defer otelx.End(span, &err)

	var data trust.SQLData
	query := p.QueryWithNetwork(ctx).
		Where("issuer = ?", issuer).
		Where("(subject = ? OR allow_any_subject IS TRUE)", subject).
		Where("key_id = ?", keyId).
		Where("nid = ?", p.NetworkID(ctx))

	if err := query.First(&data); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return p.jwtGrantFromSQlData(data).Scope, nil
}

func (p *Persister) IsJWTUsed(ctx context.Context, jti string) (ok bool, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.IsJWTUsed")
	defer otelx.End(span, &err)

	err = p.ClientAssertionJWTValid(ctx, jti)
	if err != nil {
		return true, nil
	}

	return false, nil
}

func (p *Persister) MarkJWTUsedForTime(ctx context.Context, jti string, exp time.Time) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.MarkJWTUsedForTime")
	defer otelx.End(span, &err)

	return p.SetClientAssertionJWT(ctx, jti, exp)
}

func (p *Persister) sqlDataFromJWTGrant(g trust.Grant) trust.SQLData {
	return trust.SQLData{
		ID:              g.ID,
		Issuer:          g.Issuer,
		Subject:         g.Subject,
		AllowAnySubject: g.AllowAnySubject,
		Scope:           strings.Join(g.Scope, "|"),
		KeySet:          g.PublicKey.Set,
		KeyID:           g.PublicKey.KeyID,
		CreatedAt:       g.CreatedAt,
		ExpiresAt:       g.ExpiresAt,
	}
}

func (p *Persister) jwtGrantFromSQlData(data trust.SQLData) trust.Grant {
	return trust.Grant{
		ID:              data.ID,
		Issuer:          data.Issuer,
		Subject:         data.Subject,
		AllowAnySubject: data.AllowAnySubject,
		Scope:           stringsx.Splitx(data.Scope, "|"),
		PublicKey: trust.PublicKey{
			Set:   data.KeySet,
			KeyID: data.KeyID,
		},
		CreatedAt: data.CreatedAt,
		ExpiresAt: data.ExpiresAt,
	}
}

func (p *Persister) FlushInactiveGrants(ctx context.Context, notAfter time.Time, _ int, _ int) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FlushInactiveGrants")
	defer otelx.End(span, &err)

	deleteUntil := time.Now().UTC()
	if deleteUntil.After(notAfter) {
		deleteUntil = notAfter
	}
	return sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("expires_at < ?", deleteUntil).Delete(&trust.SQLData{}))
}
