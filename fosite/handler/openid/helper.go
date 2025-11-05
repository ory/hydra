// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package openid

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/ory/hydra/v2/fosite"
)

type IDTokenHandleHelper struct {
	IDTokenStrategy OpenIDConnectTokenStrategyProvider
}

func (i *IDTokenHandleHelper) GetAccessTokenHash(ctx context.Context, requester fosite.AccessRequester, responder fosite.AccessResponder) string {
	token := responder.GetAccessToken()
	// The session should always be a openid.Session but best to safely cast
	if session, ok := requester.GetSession().(Session); ok {
		val, err := i.ComputeHash(ctx, session, token)
		if err != nil {
			// this should never happen
			panic(err)
		}

		return val
	}

	buffer := bytes.NewBufferString(token)
	hash := sha256.New()
	// sha256.digest.Write() always returns nil for err, the panic should never happen
	_, err := hash.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
	hashBuf := bytes.NewBuffer(hash.Sum([]byte{}))

	return base64.RawURLEncoding.EncodeToString(hashBuf.Bytes()[:hashBuf.Len()/2])
}

func (i *IDTokenHandleHelper) generateIDToken(ctx context.Context, lifespan time.Duration, fosr fosite.Requester) (token string, err error) {
	token, err = i.IDTokenStrategy.OpenIDConnectTokenStrategy().GenerateIDToken(ctx, lifespan, fosr)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (i *IDTokenHandleHelper) IssueImplicitIDToken(ctx context.Context, lifespan time.Duration, ar fosite.Requester, resp fosite.AuthorizeResponder) error {
	token, err := i.generateIDToken(ctx, lifespan, ar)
	if err != nil {
		return err
	}
	resp.AddParameter("id_token", token)
	return nil
}

func (i *IDTokenHandleHelper) IssueExplicitIDToken(ctx context.Context, lifespan time.Duration, ar fosite.Requester, resp fosite.AccessResponder) error {
	token, err := i.generateIDToken(ctx, lifespan, ar)
	if err != nil {
		return err
	}

	resp.SetExtra("id_token", token)
	return nil
}

// ComputeHash computes the hash using the alg defined in the id_token header
func (i *IDTokenHandleHelper) ComputeHash(ctx context.Context, sess Session, token string) (string, error) {
	var err error
	hash := sha256.New()
	if alg, ok := sess.IDTokenHeaders().Get("alg").(string); ok && len(alg) > 2 {
		if hashSize, err := strconv.Atoi(alg[2:]); err == nil {
			if hashSize == 384 {
				hash = sha512.New384()
			} else if hashSize == 512 {
				hash = sha512.New()
			}
		}
	}

	buffer := bytes.NewBufferString(token)
	_, err = hash.Write(buffer.Bytes())
	if err != nil {
		return "", err
	}
	hashBuf := bytes.NewBuffer(hash.Sum([]byte{}))

	return base64.RawURLEncoding.EncodeToString(hashBuf.Bytes()[:hashBuf.Len()/2]), nil
}
