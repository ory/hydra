package oauth2

import (
	"fmt"
	"time"

	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-errors/errors"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/oidc/strategy"
	ejwt "github.com/ory-am/fosite/token/jwt"
	"github.com/ory-am/hydra/jwk"
	"github.com/pborman/uuid"
)

const (
	ConsentChallengeKey = "consent.challenge"
	ConsentEndpointKey  = "consent.endpoint"
)

type DefaultConsentStrategy struct {
	Issuer string

	KeyManager jwk.Manager
}

func (s *DefaultConsentStrategy) ValidateResponse(a fosite.AuthorizeRequester, token string) (claims *Session, err error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		pk, err := s.KeyManager.GetKey(ConsentEndpointKey, "public")
		if err != nil {
			return nil, err
		}

		rsaKey, ok := jwk.First(pk.Keys).Key.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("Could not convert to RSA Private Key")
		}
		return rsaKey, nil
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
	return &Session{
		Subject: subject,
		Session: &strategy.DefaultSession{
			Claims: &ejwt.IDTokenClaims{
				Audience:  a.GetClient().GetID(),
				Subject:   subject,
				Issuer:    s.Issuer,
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now(),
				Extra:     t.Claims,
			},
			Headers: &ejwt.Headers{},
		},
	}, err

}

func (s *DefaultConsentStrategy) IssueChallenge(authorizeRequest fosite.AuthorizeRequester, redirectURL string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = map[string]interface{}{
		"nonce": uuid.New(),
		"scp":   authorizeRequest.GetScopes(),
		"aud":   authorizeRequest.GetClient().GetID(),
		"exp":   time.Now().Add(time.Hour).Unix(),
		"redir": redirectURL,
	}

	ks, err := s.KeyManager.GetKey(ConsentChallengeKey, "private")
	if err != nil {
		return "", errors.New(err)
	}

	rsaKey, ok := jwk.First(ks.Keys).Key.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("Could not convert to RSA Private Key")
	}

	var signature, encoded string
	if encoded, err = token.SigningString(); err != nil {
		return "", errors.New(err)
	} else if signature, err = token.Method.Sign(encoded, rsaKey); err != nil {
		return "", errors.New(err)
	}

	return fmt.Sprintf("%s.%s", encoded, signature), nil

}
