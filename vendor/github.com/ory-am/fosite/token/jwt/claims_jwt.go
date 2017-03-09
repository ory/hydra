package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
)

// JWTClaims represent a token's claims.
type JWTClaims struct {
	Subject   string
	Issuer    string
	Audience  string
	JTI       string
	IssuedAt  time.Time
	NotBefore time.Time
	ExpiresAt time.Time
	Scope     []string
	Extra     map[string]interface{}
}

// ToMap will transform the headers to a map structure
func (c *JWTClaims) ToMap() map[string]interface{} {
	var ret = Copy(c.Extra)

	ret["jti"] = c.JTI
	if c.JTI == "" {
		ret["jti"] = uuid.New()
	}

	ret["sub"] = c.Subject
	ret["iss"] = c.Issuer
	ret["aud"] = c.Audience

	if !c.IssuedAt.IsZero() {
		ret["iat"] = float64(c.IssuedAt.Unix()) // jwt-go does not support int64 as datatype
	}

	if !c.NotBefore.IsZero() {
		ret["nbf"] = float64(c.NotBefore.Unix()) // jwt-go does not support int64 as datatype
	}

	ret["exp"] = float64(c.ExpiresAt.Unix()) // jwt-go does not support int64 as datatype

	if c.Scope != nil {
		ret["scp"] = c.Scope
	}

	return ret
}

// FromMap will set the claims based on a mapping
func (c *JWTClaims) FromMap(m map[string]interface{}) {
	c.Extra = make(map[string]interface{})
	for k, v := range m {
		switch k {
		case "jti":
			if s, ok := v.(string); ok {
				c.JTI = s
			}
		case "sub":
			if s, ok := v.(string); ok {
				c.Subject = s
			}
		case "iss":
			if s, ok := v.(string); ok {
				c.Issuer = s
			}
		case "aud":
			if s, ok := v.(string); ok {
				c.Audience = s
			}
		case "iat":
			switch v.(type) {
			case float64:
				c.IssuedAt = time.Unix(int64(v.(float64)), 0)
			case int64:
				c.IssuedAt = time.Unix(v.(int64), 0)
			}
		case "nbf":
			switch v.(type) {
			case float64:
				c.NotBefore = time.Unix(int64(v.(float64)), 0)
			case int64:
				c.NotBefore = time.Unix(v.(int64), 0)
			}
		case "exp":
			switch v.(type) {
			case float64:
				c.ExpiresAt = time.Unix(int64(v.(float64)), 0)
			case int64:
				c.ExpiresAt = time.Unix(v.(int64), 0)
			}
		case "scp":
			switch v.(type) {
			case []string:
				c.Scope = v.([]string)
			case []interface{}:
				c.Scope = make([]string, len(v.([]interface{})))
				for i, vi := range v.([]interface{}) {
					if s, ok := vi.(string); ok {
						c.Scope[i] = s
					}
				}
			}
		default:
			c.Extra[k] = v
		}
	}
}

// Add will add a key-value pair to the extra field
func (c *JWTClaims) Add(key string, value interface{}) {
	if c.Extra == nil {
		c.Extra = make(map[string]interface{})
	}
	c.Extra[key] = value
}

// Get will get a value from the extra field based on a given key
func (c JWTClaims) Get(key string) interface{} {
	return c.ToMap()[key]
}

// ToMapClaims will return a jwt-go MapClaims representaion
func (c JWTClaims) ToMapClaims() jwt.MapClaims {
	return c.ToMap()
}

// FromMapClaims will populate claims from a jwt-go MapClaims representaion
func (c *JWTClaims) FromMapClaims(mc jwt.MapClaims) {
	c.FromMap(mc)
}
