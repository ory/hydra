package jwt

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"time"
)

type ClaimsCarrier map[string]interface{}

// TODO NewClaimsCarrier should require a request object instead
func NewClaimsCarrier(id, issuer, subject, audience string, expiresAt, notBefore, issuedAt time.Time) ClaimsCarrier {
	return ClaimsCarrier{
		"sub": subject,
		"jid": id,
		"iat": issuedAt.Unix(),
		"iss": issuer,
		"nbf": notBefore.Unix(),
		"aud": audience,
		"exp": expiresAt.Unix(),
	}
}

func (c ClaimsCarrier) AssertExpired() bool {
	return c.GetExpiresAt().Before(time.Now())
}

func (c ClaimsCarrier) AssertNotYetValid() bool {
	return c.GetNotBefore().After(time.Now())
}

func (c ClaimsCarrier) GetSubject() string {
	return c.getAsString("sub")
}

func (c ClaimsCarrier) GetID() string {
	return c.getAsString("jid")
}

func (c ClaimsCarrier) GetIssuedAt() time.Time {
	return c.getAsTime("iat")
}

func (c ClaimsCarrier) GetNotBefore() time.Time {
	return c.getAsTime("nbf")
}

func (c ClaimsCarrier) GetAudience() string {
	return c.getAsString("aud")
}

func (c ClaimsCarrier) GetExpiresAt() time.Time {
	return c.getAsTime("exp")
}

func (c ClaimsCarrier) GetIssuer() string {
	return c.getAsString("iss")
}

func (c ClaimsCarrier) String() string {
	result, err := json.Marshal(c)
	if err != nil {
		log.Warnf(`Could not marshal ClaimsCarrier "%v": "%v".`, c, err)
		return ""
	}
	return string(result)
}

func (c ClaimsCarrier) getAsString(key string) string {
	if s, ok := c[key]; ok {
		if r, ok := s.(string); ok {
			return r
		}
	}
	return ""
}

func (c ClaimsCarrier) getAsTime(key string) time.Time {
	ret := &time.Time{}
	if s, ok := c[key]; ok {
		if r, ok := s.(time.Time); ok {
			return r
		} else if p, ok := s.(int64); ok {
			return time.Unix(p, 0)
		} else if p, ok := s.(float64); ok {
			return time.Unix(int64(p), 0)
		}
	}
	return *ret
}
