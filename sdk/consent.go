package sdk

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	ejwt "github.com/ory-am/fosite/token/jwt"
	"github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/oauth2"
	"github.com/pkg/errors"
	"time"
)

type Consent struct {
	KeyManager jwk.Manager
}

type ResponseRequest struct {
	Challenge        string
	Subject          string
	Scopes           []string
	AccessTokenExtra interface{}
	IDTokenExtra     interface{}
}

type ChallengeClaims struct {
	RequestedScopes []string `json:"scp"`
	Audience        string   `json:"aud"`
	RedirectURL     string   `json:"redir"`
	ExpiresAt       float64  `json:"exp"`
	ID              string   `json:"jti"`
}

func (c *ChallengeClaims) Valid() error {
	if time.Now().After(ejwt.ToTime(c.ExpiresAt)) {
		return errors.Errorf("Consent challenge expired")
	}
	return nil
}

func (c *Consent) VerifyChallenge(challenge string) (*ChallengeClaims, error) {
	var claims ChallengeClaims
	t, err := jwt.ParseWithClaims(challenge, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		pk, err := c.KeyManager.GetKey(oauth2.ConsentChallengeKey, "public")
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
		return nil, errors.Wrap(err, "The consent chalnge is not a valid JSON Web Token")
	}

	if !t.Valid {
		return nil, errors.Errorf("Consent challenge is invalid")
	} else if err := claims.Valid(); err != nil {
		return nil, errors.Wrap(err, "The consent challenge claims could not be verified")
	}

	return &claims, err
}

func (c *Consent) GenerateResponse(r *ResponseRequest) (string, error) {
	challenge, err := c.VerifyChallenge(r.Challenge)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = jwt.MapClaims{
		"jti":    challenge.ID,
		"scp":    r.Scopes,
		"aud":    challenge.Audience,
		"exp":    challenge.ExpiresAt,
		"sub":    r.Subject,
		"at_ext": r.AccessTokenExtra,
		"id_ext": r.IDTokenExtra,
	}

	ks, err := c.KeyManager.GetKey(oauth2.ConsentEndpointKey, "private")
	if err != nil {
		return "", errors.WithStack(err)
	}

	rsaKey, ok := jwk.First(ks.Keys).Key.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("Could not convert to RSA Private Key")
	}

	var signature, encoded string
	if encoded, err = token.SigningString(); err != nil {
		return "", errors.WithStack(err)
	} else if signature, err = token.Method.Sign(encoded, rsaKey); err != nil {
		return "", errors.WithStack(err)
	}

	return fmt.Sprintf("%s.%s", encoded, signature), nil
}
