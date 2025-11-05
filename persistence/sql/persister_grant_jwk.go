// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/hydra/v2/fosite/handler/rfc7523"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/oauth2/trust"
	"github.com/ory/pop/v6"
	"github.com/ory/x/otelx"
	keysetpagination "github.com/ory/x/pagination/keysetpagination_v2"
	"github.com/ory/x/sqlcon"
	"github.com/ory/x/sqlxx"
)

var _ trust.GrantManager = &Persister{}

type SQLGrant struct {
	ID              uuid.UUID                      `db:"id"`
	NID             uuid.UUID                      `db:"nid"`
	Issuer          string                         `db:"issuer"`
	Subject         string                         `db:"subject"`
	AllowAnySubject bool                           `db:"allow_any_subject"`
	Scope           sqlxx.StringSlicePipeDelimiter `db:"scope"`
	KeySet          string                         `db:"key_set"`
	KeyID           string                         `db:"key_id"`
	CreatedAt       time.Time                      `db:"created_at"`
	ExpiresAt       time.Time                      `db:"expires_at"`
}

func (SQLGrant) TableName() string {
	return "hydra_oauth2_trusted_jwt_bearer_issuer"
}

func (p *Persister) RFC7523KeyStorage() rfc7523.RFC7523KeyStorage {
	return p
}

func (p *Persister) CreateGrant(ctx context.Context, g trust.Grant, publicKey jose.JSONWebKey) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.CreateGrant")
	defer otelx.End(span, &err)

	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		// add key, if it doesn't exist
		if _, err := p.d.KeyManager().GetKey(ctx, g.PublicKey.Set, g.PublicKey.KeyID); err != nil {
			if !errors.Is(err, sqlcon.ErrNoRows) {
				return sqlcon.HandleError(err)
			}

			if err = p.d.KeyManager().AddKey(ctx, g.PublicKey.Set, &publicKey); err != nil {
				return sqlcon.HandleError(err)
			}
		}

		data := SQLGrant{}.fromGrant(g)
		return sqlcon.HandleError(p.CreateWithNetwork(ctx, &data))
	})
}

func (p *Persister) GetConcreteGrant(ctx context.Context, id uuid.UUID) (_ trust.Grant, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetConcreteGrant")
	defer otelx.End(span, &err)

	var data SQLGrant
	if err := p.QueryWithNetwork(ctx).Where("id = ?", id).First(&data); err != nil {
		return trust.Grant{}, sqlcon.HandleError(err)
	}

	return data.toGrant(), nil
}

func (p *Persister) DeleteGrant(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.DeleteGrant")
	defer otelx.End(span, &err)

	return p.Transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		grant, err := p.GetConcreteGrant(ctx, id)
		if err != nil {
			return sqlcon.HandleError(err)
		}

		if err := p.QueryWithNetwork(ctx).Where("id = ?", grant.ID).Delete(&SQLGrant{}); err != nil {
			return sqlcon.HandleError(err)
		}

		return p.d.KeyManager().DeleteKey(ctx, grant.PublicKey.Set, grant.PublicKey.KeyID)
	})
}

func (p *Persister) GetGrants(ctx context.Context, optionalIssuer string, pageOpts ...keysetpagination.Option) (_ []trust.Grant, _ *keysetpagination.Paginator, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetGrants")
	defer otelx.End(span, &err)

	paginator := keysetpagination.NewPaginator(append(pageOpts,
		keysetpagination.WithDefaultToken(keysetpagination.NewPageToken(keysetpagination.Column{Name: "id", Value: uuid.Nil})),
	)...)

	var grantsData []SQLGrant
	query := p.QueryWithNetwork(ctx).Scope(keysetpagination.Paginate[SQLGrant](paginator))
	if optionalIssuer != "" {
		query = query.Where("issuer = ?", optionalIssuer)
	}

	if err := query.All(&grantsData); err != nil {
		return nil, nil, sqlcon.HandleError(err)
	}
	grantsData, nextPage := keysetpagination.Result(grantsData, paginator)

	grants := make([]trust.Grant, len(grantsData))
	for i := range grantsData {
		grants[i] = grantsData[i].toGrant()
	}

	return grants, nextPage, nil
}

func (p *Persister) GetPublicKey(ctx context.Context, issuer string, subject string, keyId string) (_ *jose.JSONWebKey, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetPublicKey")
	defer otelx.End(span, &err)

	tableName := SQLGrant{}.TableName()
	// Index hint.
	if p.Connection(ctx).Dialect.Name() == "cockroach" {
		tableName += "@hydra_oauth2_trusted_jwt_bearer_issuer_nid_uq_idx"
	}

	sql := fmt.Sprintf(`SELECT key_set FROM %s WHERE key_id = ? AND nid = ? AND issuer = ? AND (subject = ? OR allow_any_subject IS TRUE) LIMIT 1`, tableName)
	query := p.Connection(ctx).RawQuery(sql,
		keyId, p.NetworkID(ctx), issuer, subject,
	)
	var keySetID string
	if err := query.First(&keySetID); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	// TODO: Consider merging this query with the one above using a `JOIN`.
	keySet, err := p.d.KeyManager().GetKey(ctx, keySetID, keyId)
	if err != nil {
		return nil, err
	}

	return &keySet.Keys[0], nil
}

func (p *Persister) GetPublicKeys(ctx context.Context, issuer string, subject string) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetPublicKeys")
	defer otelx.End(span, &err)

	q := p.QueryWithNetwork(ctx)
	expiresAt := "expires_at > NOW()"
	if q.Connection.Dialect.Name() == "sqlite3" {
		expiresAt = "expires_at > datetime('now')"
	}

	grantsData := make([]SQLGrant, 0)
	query := q.
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

	return js.ToJWK(ctx, p.r.KeyCipher())
}

func (p *Persister) GetPublicKeyScopes(ctx context.Context, issuer string, subject string, keyId string) (_ []string, err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.GetPublicKeyScopes")
	defer otelx.End(span, &err)

	tableName := SQLGrant{}.TableName()
	// Index hint.
	if p.Connection(ctx).Dialect.Name() == "cockroach" {
		tableName += "@hydra_oauth2_trusted_jwt_bearer_issuer_nid_uq_idx"
	}

	sql := fmt.Sprintf(`SELECT scope FROM %s WHERE key_id = ? AND nid = ? AND issuer = ? AND (subject = ? OR allow_any_subject IS TRUE) LIMIT 1`, tableName)
	query := p.Connection(ctx).RawQuery(sql,
		keyId, p.NetworkID(ctx), issuer, subject,
	)
	var scopes sqlxx.StringSlicePipeDelimiter
	if err := query.First(&scopes); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return scopes, nil
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

func (SQLGrant) fromGrant(g trust.Grant) SQLGrant {
	return SQLGrant{
		ID:              g.ID,
		Issuer:          g.Issuer,
		Subject:         g.Subject,
		AllowAnySubject: g.AllowAnySubject,
		Scope:           g.Scope,
		KeySet:          g.PublicKey.Set,
		KeyID:           g.PublicKey.KeyID,
		CreatedAt:       g.CreatedAt,
		ExpiresAt:       g.ExpiresAt,
	}
}

func (d SQLGrant) toGrant() trust.Grant {
	return trust.Grant{
		ID:              d.ID,
		Issuer:          d.Issuer,
		Subject:         d.Subject,
		AllowAnySubject: d.AllowAnySubject,
		Scope:           d.Scope,
		PublicKey: trust.PublicKey{
			Set:   d.KeySet,
			KeyID: d.KeyID,
		},
		CreatedAt: d.CreatedAt,
		ExpiresAt: d.ExpiresAt,
	}
}

func (p *Persister) FlushInactiveGrants(ctx context.Context, notAfter time.Time, _ int, _ int) (err error) {
	ctx, span := p.r.Tracer(ctx).Tracer().Start(ctx, "persistence.sql.FlushInactiveGrants")
	defer otelx.End(span, &err)

	deleteUntil := time.Now().UTC()
	if deleteUntil.After(notAfter) {
		deleteUntil = notAfter
	}
	return sqlcon.HandleError(p.QueryWithNetwork(ctx).Where("expires_at < ?", deleteUntil).Delete(&SQLGrant{}))
}
