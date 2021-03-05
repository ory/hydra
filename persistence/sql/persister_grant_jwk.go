package sql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"
	"gopkg.in/square/go-jose.v2"

	"github.com/ory/hydra/grant/jwtbearer"
	"github.com/ory/x/errorsx"

	"github.com/ory/x/sqlcon"
)

var _ jwtbearer.GrantManager = &Persister{}

const scopeSeparator = " "

func (p *Persister) CreateGrant(ctx context.Context, g jwtbearer.Grant, publicKey jose.JSONWebKey) error {
	// add key, if it doesn't exist
	if _, err := p.GetKey(ctx, g.PublicKey.Set, g.PublicKey.KeyID); err != nil {
		if errorsx.Cause(err) != sqlcon.ErrNoRows {
			return err
		}

		if err = p.AddKey(ctx, g.PublicKey.Set, &publicKey); err != nil {
			return err
		}
	}

	data := p.sqlDataFromJWTGrant(g)

	return sqlcon.HandleError(p.Connection(ctx).Create(&data))
}

func (p *Persister) GetConcreteGrant(ctx context.Context, id string) (jwtbearer.Grant, error) {
	var data jwtbearer.SQLData
	if err := p.Connection(ctx).Where("id = ?", id).First(&data); err != nil {
		return jwtbearer.Grant{}, sqlcon.HandleError(err)
	}

	return p.jwtGrantFromSQlData(data), nil
}

func (p *Persister) DeleteGrant(ctx context.Context, id string) error {
	grant, err := p.GetConcreteGrant(ctx, id)
	if err != nil {
		return err
	}

	return p.transaction(ctx, func(ctx context.Context, c *pop.Connection) error {
		if err := p.Connection(ctx).Destroy(&jwtbearer.SQLData{ID: grant.ID}); err != nil {
			return sqlcon.HandleError(err)
		}

		return p.DeleteKey(ctx, grant.PublicKey.Set, grant.PublicKey.KeyID)
	})
}

func (p *Persister) GetGrants(ctx context.Context, limit, offset int, optionalIssuer string) ([]jwtbearer.Grant, error) {
	grantsData := make([]jwtbearer.SQLData, 0)

	query := p.Connection(ctx).Paginate(offset/limit+1, limit).Order("id")
	if optionalIssuer != "" {
		query = query.Where("issuer = ?", optionalIssuer)
	}

	if err := query.All(&grantsData); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	grants := make([]jwtbearer.Grant, 0, len(grantsData))
	for _, data := range grantsData {
		grants = append(grants, p.jwtGrantFromSQlData(data))
	}

	return grants, nil
}

func (p *Persister) CountGrants(ctx context.Context) (int, error) {
	n, err := p.Connection(ctx).Count(&jwtbearer.SQLData{})
	return n, sqlcon.HandleError(err)
}

func (p *Persister) GetPublicKey(ctx context.Context, issuer string, subject string, keyId string) (*jose.JSONWebKey, error) {
	var data jwtbearer.SQLData
	query := p.Connection(ctx).
		Where("issuer = ?", issuer).
		Where("subject = ?", subject).
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
	grantsData := make([]jwtbearer.SQLData, 0)
	query := p.Connection(ctx).
		Where("issuer = ?", issuer).
		Where("subject = ?", subject)
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
	var data jwtbearer.SQLData
	query := p.Connection(ctx).
		Where("issuer = ?", issuer).
		Where("subject = ?", subject).
		Where("key_id = ?", keyId)
	if err := query.First(&data); err != nil {
		return nil, sqlcon.HandleError(err)
	}

	return strings.Split(data.Scope, scopeSeparator), nil
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

func (p *Persister) sqlDataFromJWTGrant(g jwtbearer.Grant) jwtbearer.SQLData {
	return jwtbearer.SQLData{
		ID:        g.ID,
		Issuer:    g.Issuer,
		Subject:   g.Subject,
		Scope:     strings.Join(g.Scope, " "),
		KeySet:    g.PublicKey.Set,
		KeyID:     g.PublicKey.KeyID,
		CreatedAt: g.CreatedAt,
		ExpiresAt: g.ExpiresAt,
	}
}

func (p *Persister) jwtGrantFromSQlData(data jwtbearer.SQLData) jwtbearer.Grant {
	return jwtbearer.Grant{
		ID:      data.ID,
		Issuer:  data.Issuer,
		Subject: data.Subject,
		Scope:   strings.Split(data.Scope, scopeSeparator),
		PublicKey: jwtbearer.PublicKey{
			Set:   data.KeySet,
			KeyID: data.KeyID,
		},
		CreatedAt: data.CreatedAt,
		ExpiresAt: data.ExpiresAt,
	}
}

func (p *Persister) FlushInactiveGrants(ctx context.Context, notAfter time.Time) error {
	return sqlcon.HandleError(p.Connection(ctx).RawQuery(
		fmt.Sprintf("DELETE FROM %s WHERE expires_at < ? AND expires_at < ?", jwtbearer.SQLData{}.TableName()),
		time.Now().UTC(),
		notAfter,
	).Exec())
}
