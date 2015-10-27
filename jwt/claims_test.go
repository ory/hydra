package jwt_test

import (
	. "github.com/ory-am/hydra/jwt"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClaimsCarrier(t *testing.T) {
	type set struct {
		id        string
		issuer    string
		subject   string
		audience  string
		notBefore time.Time
		issuedAt  time.Time
	}
	for k, c := range []*set{
		&set{uuid.New(), "hydra", "peter", "app", time.Now(), time.Now()},
	} {
		carrier := NewClaimsCarrier(c.id, c.issuer, c.subject, c.audience, c.notBefore, c.issuedAt)
		assert.Equal(t, c.id, carrier.ID(), "Case %d", k)
		assert.Equal(t, c.issuer, carrier.Issuer(), "Case %d", k)
		assert.Equal(t, c.subject, carrier.Subject(), "Case %d", k)
		assert.Equal(t, c.audience, carrier.Audience(), "Case %d", k)
		assert.Equal(t, c.notBefore, carrier.NotBefore(), "Case %d", k)
		assert.Equal(t, c.issuedAt, carrier.IssuedAt(), "Case %d", k)
	}
}
