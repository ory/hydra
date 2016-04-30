package oauth2

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-errors/errors"
	"github.com/ory-am/fosite"
	ejwt "github.com/ory-am/fosite/enigma/jwt"
	"github.com/ory-am/fosite/handler/oidc/strategy"
	"github.com/ory-am/hydra/key"
	"github.com/pborman/uuid"
)

const (
	ConsentChallengeKey = "consentChallenge"
	ConsentEndpointKey  = "consentEndpoint"
)

type DefaultConsentStrategy struct {
	Issuer string

	KeyManager key.Manager
}

func (s *DefaultConsentStrategy) ValidateResponse(a fosite.AuthorizeRequester, token string) (claims *Session, err error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		pk, err := s.KeyManager.GetAsymmetricKey(ConsentEndpointKey)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPublicKeyFromPEM(pk.Public)
	})

	if err != nil {
		return nil, errors.Errorf("Couldn't parse token: %v", err)
	} else if !t.Valid {
		return nil, errors.Errorf("Token is invalid")
	}

	if time.Now().After(ejwt.ToTime(t.Claims["exp"])) {
		return nil, errors.Errorf("Token expired")
	}

	if ejwt.ToString(t.Claims["aud"]) != a.GetClient().GetID() {
		return nil, errors.Errorf("Audience mismatch")
	}

	subject := ejwt.ToString(t.Claims["sub"])
	delete(t.Claims, "sub")
	return &Session{
		Subject: subject,
		IDTokenSession: &strategy.IDTokenSession{
			IDTokenClaims: &strategy.IDTokenClaims{
				Audience:  a.GetClient().GetID(),
				Subject:   subject,
				Issuer:    s.Issuer,
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now(),
				Extra:     t.Claims,
			},
			Header: &ejwt.Header{},
		},
	}, err

}

func (s *DefaultConsentStrategy) IssueChallenge(authorizeRequest fosite.AuthorizeRequester, redirectURL string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = map[string]interface{}{
		"nonce": uuid.New(),
		"scp":   authorizeRequest.GetScopes(),
		"aud":   authorizeRequest.GetClient().GetID(),
		"exp":   time.Now().Add(time.Hour),
		"redir": redirectURL,
	}

	key, err := s.KeyManager.GetAsymmetricKey(ConsentChallengeKey)
	if err != nil {
		return "", errors.New(err)
	}

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(key.Private)
	if err != nil {
		return "", errors.New(err)
	}

	var signature, encoded string
	if encoded, err = token.SigningString(); err != nil {
		return "", errors.New(err)
	} else if signature, err = token.Method.Sign(encoded, rsaKey); err != nil {
		return "", errors.New(err)
	}

	return fmt.Sprintf("%s.%s", encoded, signature), nil

}
