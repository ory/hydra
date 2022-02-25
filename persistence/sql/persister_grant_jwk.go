package sql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"gopkg.in/square/go-jose.v2"

	"github.com/ory/hydra/oauth2/trust"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/stringsx"

	"github.com/ory/x/sqlcon"
)

var _ trust.GrantManager = &Persister{}

func (p *Persister) CreateGrant(ctx context.Context, g trust.Grant, publicKey jose.JSONWebKey) error {
	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		// add key, if it doesn't exist
		if _, err := p.GetKey(ctx, g.PublicKey.Set, g.PublicKey.KeyID); err != nil {
			if errorsx.Cause(err) != sqlcon.ErrNoRows {
				return sqlcon.HandleError(err)
			}

			if err = p.AddKey(ctx, g.PublicKey.Set, &publicKey); err != nil {
				return sqlcon.HandleError(err)
			}
		}

		data := p.sqlDataFromJWTGrant(g)

		return sqlcon.HandleError(p.Connection(ctx).Create(&data))
	})
}

func (p *Persister) GetConcreteGrant(ctx context.Context, id string) (trust.Grant, error) {
	var data trust.SQLData
	if err := p.Connection(ctx).Where("id = ?", id).First(&data); err != nil {
		return trust.Grant{}, sqlcon.HandleError(err)
	}

	return p.jwtGrantFromSQlData(data), nil
}

func (p *Persister) DeleteGrant(ctx context.Context, id string) error {
	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		grant, err := p.GetConcreteGrant(ctx, id)
		if err != nil {
			return sqlcon.HandleError(err)
		}

		if err := p.Connection(ctx).Destroy(&trust.SQLData{ID: grant.ID}); err != nil {
			return sqlcon.HandleError(err)
		}

		return p.DeleteKey(ctx, grant.PublicKey.Set, grant.PublicKey.KeyID)
	})
}

func (p *Persister) GetGrants(ctx context.Context, limit, offset int, optionalIssuer string) ([]trust.Grant, error) {
	grantsData := make([]trust.SQLData, 0)

	query := p.Connection(ctx).Paginate(offset/limit+1, limit).Order("id")
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

func (p *Persister) CountGrants(ctx context.Context) (int, error) {
	n, err := p.Connection(ctx).Count(&trust.SQLData{})
	return n, sqlcon.HandleError(err)
}

func (p *Persister) GetPublicKey(ctx context.Context, issuer string, subject string, keyId string) (*jose.JSONWebKey, error) {
	var data trust.SQLData
	query := p.Connection(ctx).
		Where("issuer = ?", issuer).
		Where("subject = ? OR allow_any_subject IS TRUE", subject).
		Where("key_id = ?", keyId)
	if err := query.First(&data); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	keySet, err := p.GetKey(ctx, data.KeySet, keyId)
	if err != nil {
		return nil, err
	}

	return &keySet.Keys[0], nil
}

func (p *Persister) GetPublicKeys(ctx context.Context, issuer string, subject string) (*jose.JSONWebKeySet, error) {
	grantsData := make([]trust.SQLData, 0)
	query := p.Connection(ctx).
		Where("issuer = ?", issuer).
		Where("subject = ? OR allow_any_subject IS TRUE", subject)
	if err := query.All(&grantsData); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	if len(grantsData) == 0 {
		return &jose.JSONWebKeySet{}, nil
	}

	// because keys must be grouped by issuer, we can retrieve set name from first grant
	keySet, err := p.GetKeySet(ctx, grantsData[0].KeySet)
	if err != nil {
		return nil, err
	}

	// find keys, that belong to grants
	filteredKeySet := &jose.JSONWebKeySet{}
	for _, data := range grantsData {
		if keys := keySet.Key(data.KeyID); len(keys) > 0 {
			filteredKeySet.Keys = append(filteredKeySet.Keys, keys...)
		}
	}

	return filteredKeySet, nil
}

func (p *Persister) GetPublicKeyScopes(ctx context.Context, issuer string, subject string, keyId string) ([]string, error) {
	var data trust.SQLData
	query := p.Connection(ctx).
		Where("issuer = ?", issuer).
		Where("subject = ? OR allow_any_subject IS TRUE", subject).
		Where("key_id = ?", keyId)
	if err := query.First(&data); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return p.jwtGrantFromSQlData(data).Scope, nil
}

func (p *Persister) IsJWTUsed(ctx context.Context, jti string) (bool, error) {
	err := p.ClientAssertionJWTValid(ctx, jti)
	if err != nil {
		return true, nil
	}

	return false, nil
}

func (p *Persister) MarkJWTUsedForTime(ctx context.Context, jti string, exp time.Time) error {
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

func (p *Persister) FlushInactiveGrants(ctx context.Context, notAfter time.Time, limit int, batchSize int) error {
	deleteUntil := time.Now().UTC()
	if deleteUntil.After(notAfter) {
		deleteUntil = notAfter
	}
	return sqlcon.HandleError(p.Connection(ctx).RawQuery(
		fmt.Sprintf("DELETE FROM %s WHERE expires_at < ?", trust.SQLData{}.TableName()),
		deleteUntil,
	).Exec())
}
