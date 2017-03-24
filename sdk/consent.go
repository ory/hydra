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

// Consent is a helper for singing and verifying consent challenges. For an exemplary reference implementation, check
// https://github.com/ory/hydra-consent-app-go
type Consent struct {
	KeyManager jwk.Manager
}

// ResponseRequest is being used by the consent response singing helper.
type ResponseRequest struct {
	// Challenge is the original consent challenge.
	Challenge        string

	// Subject will be the sub claim of the access token. Usually this is a resource owner (user).
	Subject          string

	// Scopes are the scopes the resource owner granted to the application requesting the access token.
	Scopes           []string

	// AccessTokenExtra is arbitrary data that will be available when performing token introspection or warden requests.
	AccessTokenExtra interface{}

	// IDTokenExtra is arbitrary data that will included as a claim in the ID Token, if requested.
	IDTokenExtra     interface{}
}

// ChallengeClaims are the decoded claims of a consent challenge.
type ChallengeClaims struct {
	// RequestedScopes are the scopes the application requested. Each scope should be explicitly granted by
	// the user.
	RequestedScopes []string `json:"scp"`

	// The ID of the application that initiated the OAuth2 flow.
	Audience        string   `json:"aud"`

	// RedirectURL is the url where the consent app will send the user after the consent flow has been completed.
	RedirectURL     string   `json:"redir"`

	// ExpiresAt is a unix timestamp of the expiry time.
	ExpiresAt       float64  `json:"exp"`

	// ID is the tokens' ID which will be automatically echoed in the consent response.
	ID              string   `json:"jti"`
}

// Valid tests if the challenge's claims are valid.
func (c *ChallengeClaims) Valid() error {
	if time.Now().After(ejwt.ToTime(c.ExpiresAt)) {
		return errors.Errorf("Consent challenge expired")
	}
	return nil
}

// VerifyChallenge verifies a consent challenge and either returns the challenge's claims if it is valid, or an
// error if it is not.
//
//  claims, err := c.VerifyChallenge(challenge)
//  if err != nil {
//    // The challenge is invalid, or the signing key could not be retrieved
//  }
//  // ...
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

// DenyConsent can be used to indicate that the user denied consent. Returns a redirect url or an error
// if the challenge is invalid.
//
//  redirectUrl, _ := c.DenyConsent(challenge)
//  http.Redirect(w, r, redirectUrl, http.StatusFound)
func (c *Consent) DenyConsent(challenge string) (string, error) {
	claims, err := c.VerifyChallenge(challenge)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s&consent=denied", claims.RedirectURL), nil
}

// GenerateResponse generates a consent response and returns the consent response token, or an error if it is invalid.
//
//  redirectUrl, _ := c.GenerateResponse(challenge)
//  http.Redirect(w, r, redirectUrl, http.StatusFound)
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

	return fmt.Sprintf("%s&consent=%s.%s", challenge.RedirectURL, encoded, signature), nil
}
