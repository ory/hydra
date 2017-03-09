package oauth2

import (
	"bytes"
	"encoding/gob"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/token/jwt"
	"time"
)

type JWTSessionContainer interface {
	// GetJWTClaims returns the claims.
	GetJWTClaims() *jwt.JWTClaims

	// GetJWTHeader returns the header.
	GetJWTHeader() *jwt.Headers

	fosite.Session
}

// JWTSession Container for the JWT session.
type JWTSession struct {
	JWTClaims *jwt.JWTClaims
	JWTHeader *jwt.Headers
	ExpiresAt map[fosite.TokenType]time.Time
	Username  string
	Subject   string
}

func (j *JWTSession) GetJWTClaims() *jwt.JWTClaims {
	if j.JWTClaims == nil {
		j.JWTClaims = &jwt.JWTClaims{}
	}
	return j.JWTClaims
}

func (j *JWTSession) GetJWTHeader() *jwt.Headers {
	if j.JWTHeader == nil {
		j.JWTHeader = &jwt.Headers{}
	}
	return j.JWTHeader
}

func (s *JWTSession) SetExpiresAt(key fosite.TokenType, exp time.Time) {
	if s.ExpiresAt == nil {
		s.ExpiresAt = make(map[fosite.TokenType]time.Time)
	}
	s.ExpiresAt[key] = exp
}

func (s *JWTSession) GetExpiresAt(key fosite.TokenType) time.Time {
	if s.ExpiresAt == nil {
		s.ExpiresAt = make(map[fosite.TokenType]time.Time)
	}

	if _, ok := s.ExpiresAt[key]; !ok {
		return time.Time{}
	}
	return s.ExpiresAt[key]
}

func (s *JWTSession) GetUsername() string {
	if s == nil {
		return ""
	}
	return s.Username
}

func (s *JWTSession) GetSubject() string {
	if s == nil {
		return ""
	}

	return s.Subject
}

func (s *JWTSession) Clone() fosite.Session {
	if s == nil {
		return nil
	}

	var clone JWTSession
	var mod bytes.Buffer
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)
	_ = enc.Encode(s)
	_ = dec.Decode(&clone)
	return &clone
}
