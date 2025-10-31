// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"time"

	"github.com/mohae/deepcopy"
)

// Session is an interface that is used to store session data between OAuth2 requests. It can be used to look up
// when a session expires or what the subject's name was.
type Session interface {
	// SetExpiresAt sets the expiration time of a token.
	//
	//  session.SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(time.Hour))
	SetExpiresAt(key TokenType, exp time.Time)

	// GetExpiresAt returns the expiration time of a token if set, or time.IsZero() if not.
	//
	//  session.GetExpiresAt(fosite.AccessToken)
	GetExpiresAt(key TokenType) time.Time

	// GetUsername returns the username, if set. This is optional and only used during token introspection.
	GetUsername() string

	// GetSubject returns the subject, if set. This is optional and only used during token introspection.
	GetSubject() string

	// Clone clones the session.
	Clone() Session
}

// DefaultSession is a default implementation of the session interface.
type DefaultSession struct {
	ExpiresAt map[TokenType]time.Time `json:"expires_at"`
	Username  string                  `json:"username"`
	Subject   string                  `json:"subject"`
	Extra     map[string]interface{}  `json:"extra"`
}

func (s *DefaultSession) SetExpiresAt(key TokenType, exp time.Time) {
	if s.ExpiresAt == nil {
		s.ExpiresAt = make(map[TokenType]time.Time)
	}
	s.ExpiresAt[key] = exp
}

func (s *DefaultSession) GetExpiresAt(key TokenType) time.Time {
	if s.ExpiresAt == nil {
		s.ExpiresAt = make(map[TokenType]time.Time)
	}

	return s.ExpiresAt[key]
}

func (s *DefaultSession) GetUsername() string {
	if s == nil {
		return ""
	}
	return s.Username
}

func (s *DefaultSession) SetSubject(subject string) {
	s.Subject = subject
}

func (s *DefaultSession) GetSubject() string {
	if s == nil {
		return ""
	}

	return s.Subject
}

func (s *DefaultSession) Clone() Session {
	if s == nil {
		return nil
	}

	return deepcopy.Copy(s).(Session)
}

// ExtraClaimsSession provides an interface for session to store any extra claims.
type ExtraClaimsSession interface {
	// GetExtraClaims returns a map to store extra claims.
	// The returned value can be modified in-place.
	GetExtraClaims() map[string]interface{}
}

// GetExtraClaims implements ExtraClaimsSession for DefaultSession.
// The returned value can be modified in-place.
func (s *DefaultSession) GetExtraClaims() map[string]interface{} {
	if s == nil {
		return nil
	}

	if s.Extra == nil {
		s.Extra = make(map[string]interface{})
	}

	return s.Extra
}
