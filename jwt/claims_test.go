package jwt_test

import (
	. "github.com/ory-am/hydra/jwt"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClaimsCarrier(t *testing.T) {
	for k, c := range []struct {
		id        string
		issuer    string
		subject   string
		audience  string
		notBefore time.Time
		issuedAt  time.Time
	}{
		{uuid.New(), "hydra", "peter", "app", time.Now(), time.Now()},
	} {
		carrier := NewClaimsCarrier(c.id, c.issuer, c.subject, c.audience, c.notBefore, c.issuedAt)
		assert.Equal(t, c.id, carrier.GetID(), "Case %d", k)
		assert.Equal(t, c.issuer, carrier.GetIssuer(), "Case %d", k)
		assert.Equal(t, c.subject, carrier.GetSubject(), "Case %d", k)
		assert.Equal(t, c.audience, carrier.GetAudience(), "Case %d", k)
		assert.Equal(t, c.notBefore, carrier.GetNotBefore(), "Case %d", k)
		assert.Equal(t, c.issuedAt, carrier.GetIssuedAt(), "Case %d", k)
	}
}
