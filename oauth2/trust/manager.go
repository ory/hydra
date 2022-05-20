package trust

import (
	"context"
	"time"

	"gopkg.in/square/go-jose.v2"
)

type GrantManager interface {
	CreateGrant(ctx context.Context, g Grant, publicKey jose.JSONWebKey) error
	GetConcreteGrant(ctx context.Context, id string) (Grant, error)
	DeleteGrant(ctx context.Context, id string) error
	GetGrants(ctx context.Context, limit, offset int, optionalIssuer string) ([]Grant, error)
	CountGrants(ctx context.Context) (int, error)
	FlushInactiveGrants(ctx context.Context, notAfter time.Time, limit int, batchSize int) error
}

type SQLData struct {
	ID              string    `db:"id"`
	Issuer          string    `db:"issuer"`
	Subject         string    `db:"subject"`
	AllowAnySubject bool      `db:"allow_any_subject"`
	Scope           string    `db:"scope"`
	KeySet          string    `db:"key_set"`
	KeyID           string    `db:"key_id"`
	CreatedAt       time.Time `db:"created_at"`
	ExpiresAt       time.Time `db:"expires_at"`
}

func (SQLData) TableName() string {
	return "hydra_oauth2_trusted_jwt_bearer_issuer"
}
