/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @Copyright 	2017-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package jwk

import (
	"context"
	"github.com/ory/hydra/driver/config"
	"strings"

	"github.com/pkg/errors"

	"github.com/ory/fosite/token/jwt"
)

type JWTSigner interface {
	GetPublicKeyID(ctx context.Context) (string, error)
	jwt.Signer
}

type DefaultJWTSigner struct {
	*jwt.DefaultSigner
	r     InternalRegistry
	c     *config.DefaultProvider
	setID string
}

func NewDefaultJWTSigner(c config.DefaultProvider, r InternalRegistry, keyID string) (*DefaultJWTSigner, error) {
	j := &DefaultJWTSigner{c: &c, r: r, setID: keyID, DefaultSigner: &jwt.DefaultSigner{}}
	j.DefaultSigner.GetPrivateKey = j.getPrivateKey
	return j, nil
}

func (j *DefaultJWTSigner) GetPublicKeyID(ctx context.Context) (string, error) {
	keys, err := j.r.KeyManager().GetKeySet(ctx, j.setID)
	if err != nil {
		return "", err
	}

	public, err := FindPublicKey(keys)
	if err != nil {
		return "", err
	}

	return public.KeyID, nil
}

func (j *DefaultJWTSigner) getPrivateKey(ctx context.Context) (interface{}, error) {
	keys, err := j.r.KeyManager().GetKeySet(ctx, j.setID)
	if err != nil {
		return nil, err
	}

	public, err := FindPublicKey(keys)
	if err != nil {
		return nil, err
	}

	private, err := FindPrivateKey(keys)
	if err != nil {
		return nil, err
	}

	if strings.Replace(public.KeyID, "public:", "", 1) != strings.Replace(private.KeyID, "private:", "", 1) {
		return nil, errors.New("public and private key pair kids do not match")
	}

	return private, nil
}
