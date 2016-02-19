package jwt

import (
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/pborman/uuid"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClaimsCarrier(t *testing.T) {
	for k, c := range []struct {
		id        string
		issuer    string
		subject   string
		audience  string
		expiresAt time.Time
		notBefore time.Time
		issuedAt  time.Time
	}{
		{uuid.New(), "hydra", "peter", "app", time.Now().Add(time.Hour), time.Now(), time.Now()},
	} {
		carrier := NewClaimsCarrier(c.id, c.issuer, c.subject, c.audience, c.expiresAt, c.notBefore, c.issuedAt)
		assert.Equal(t, c.id, carrier.GetID(), "Case %d", k)
		assert.Equal(t, c.issuer, carrier.GetIssuer(), "Case %d", k)
		assert.Equal(t, c.subject, carrier.GetSubject(), "Case %d", k)
		assert.Equal(t, c.audience, carrier.GetAudience(), "Case %d", k)

		assert.Equal(t, c.notBefore.Day(), carrier.GetNotBefore().Day(), "Case %d", k)
		assert.Equal(t, c.issuedAt.Day(), carrier.GetIssuedAt().Day(), "Case %d", k)
		assert.Equal(t, c.expiresAt.Day(), carrier.GetExpiresAt().Day(), "Case %d", k)

		assert.Empty(t, carrier.getAsString("doesnotexist"), "Case %d", k)
		assert.Equal(t, time.Time{}, carrier.getAsTime("doesnotexist"), "Case %d", k)
		assert.NotEmpty(t, carrier.String(), "Case %d", k)
	}
}
