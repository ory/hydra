package config

import (
	"github.com/ory/x/stringsx"
	"strings"
)

// AccessTokenStrategyType is the type of access token strategy.
type AccessTokenStrategyType string

const (
	// AccessTokenJWTStrategy is the JWT access token strategy.
	AccessTokenJWTStrategy AccessTokenStrategyType = "jwt"
	// AccessTokenDefaultStrategy is the default access token strategy using HMAC-SHA pass-by-reference tokens.
	AccessTokenDefaultStrategy AccessTokenStrategyType = "opaque"
)

// ToAccessTokenStrategyType converts a string to an AccessTokenStrategyType
func ToAccessTokenStrategyType(strategy string) (AccessTokenStrategyType, error) {
	switch f := stringsx.SwitchExact(strings.ToLower(strategy)); {
	case f.AddCase("jwt"):
		return AccessTokenJWTStrategy, nil
	case f.AddCase("opaque"):
		return AccessTokenDefaultStrategy, nil
	default:
		return "", f.ToUnknownCaseErr()
	}
}

// ToAccessTokenStrategyTypeOrDefault converts a string to an AccessTokenStrategyType or the default
// if argument strategy is not a valid AccessTokenStrategyType.
func ToAccessTokenStrategyTypeOrDefault(strategy string) AccessTokenStrategyType {
	if s, err := ToAccessTokenStrategyType(strategy); err == nil {
		return s
	}
	return AccessTokenDefaultStrategy
}
