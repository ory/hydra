package openid

import (
	"encoding/gob"
	"net/http"
	"time"

	"bytes"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/token/jwt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const defaultExpiryTime = time.Hour

type Session interface {
	IDTokenClaims() *jwt.IDTokenClaims
	IDTokenHeaders() *jwt.Headers

	fosite.Session
}

// IDTokenSession is a session container for the id token
type DefaultSession struct {
	Claims    *jwt.IDTokenClaims
	Headers   *jwt.Headers
	ExpiresAt map[fosite.TokenType]time.Time
	Username  string
	Subject   string
}

func NewDefaultSession() *DefaultSession {
	return &DefaultSession{
		Claims:  &jwt.IDTokenClaims{},
		Headers: &jwt.Headers{},
	}
}

func (s *DefaultSession) Clone() fosite.Session {
	if s == nil {
		return nil
	}

	var clone DefaultSession
	var mod bytes.Buffer
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)
	_ = enc.Encode(s)
	_ = dec.Decode(&clone)
	return &clone
}

func (s *DefaultSession) SetExpiresAt(key fosite.TokenType, exp time.Time) {
	if s.ExpiresAt == nil {
		s.ExpiresAt = make(map[fosite.TokenType]time.Time)
	}
	s.ExpiresAt[key] = exp
}

func (s *DefaultSession) GetExpiresAt(key fosite.TokenType) time.Time {
	if s.ExpiresAt == nil {
		s.ExpiresAt = make(map[fosite.TokenType]time.Time)
	}

	if _, ok := s.ExpiresAt[key]; !ok {
		return time.Time{}
	}
	return s.ExpiresAt[key]
}

func (s *DefaultSession) GetUsername() string {
	if s == nil {
		return ""
	}
	return s.Username
}

func (s *DefaultSession) GetSubject() string {
	if s == nil {
		return ""
	}

	return s.Subject
}

func (s *DefaultSession) IDTokenHeaders() *jwt.Headers {
	if s.Headers == nil {
		s.Headers = &jwt.Headers{}
	}
	return s.Headers
}

func (s *DefaultSession) IDTokenClaims() *jwt.IDTokenClaims {
	if s.Claims == nil {
		s.Claims = &jwt.IDTokenClaims{}
	}
	return s.Claims
}

type DefaultStrategy struct {
	*jwt.RS256JWTStrategy

	Expiry time.Duration
	Issuer string
}

func (h DefaultStrategy) GenerateIDToken(_ context.Context, _ *http.Request, requester fosite.Requester) (token string, err error) {
	if h.Expiry == 0 {
		h.Expiry = defaultExpiryTime
	}

	sess, ok := requester.GetSession().(Session)
	if !ok {
		return "", errors.New("Session must be of type strategy.Session")
	}

	claims := sess.IDTokenClaims()
	if requester.GetRequestForm().Get("max_age") != "" && (claims.AuthTime.IsZero() || claims.AuthTime.After(time.Now())) {
		return "", errors.New("Authentication time claim is required when max_age is set and can not be in the future")
	}

	if claims.Subject == "" {
		return "", errors.New("Subject claim can not be empty")
	}

	if claims.ExpiresAt.IsZero() {
		claims.ExpiresAt = time.Now().Add(h.Expiry)
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return "", errors.New("Expiry claim can not be in the past")
	}

	if claims.AuthTime.IsZero() {
		claims.AuthTime = time.Now()
	}

	if claims.Issuer == "" {
		claims.Issuer = h.Issuer
	}

	nonce := requester.GetRequestForm().Get("nonce")
	// OPTIONAL. String value used to associate a Client session with an ID Token, and to mitigate replay attacks.
	// Although optional, this is considered good practice and therefore enforced.
	if len(nonce) < fosite.MinParameterEntropy {
		// We're assuming that using less then 8 characters for the state can not be considered "unguessable"
		return "", errors.WithStack(fosite.ErrInsufficientEntropy)
	}

	claims.Nonce = nonce
	claims.Audience = requester.GetClient().GetID()
	claims.IssuedAt = time.Now()

	token, _, err = h.RS256JWTStrategy.Generate(claims.ToMapClaims(), sess.IDTokenHeaders())
	return token, err
}
