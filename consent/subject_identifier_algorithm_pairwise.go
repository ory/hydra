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

package consent

import (
	"crypto/sha256"
	"fmt"
	"net/url"

	"github.com/ory/x/errorsx"

	"github.com/ory/fosite"
	"github.com/ory/hydra/client"
)

type SubjectIdentifierAlgorithmPairwise struct {
	Salt []byte
}

func NewSubjectIdentifierAlgorithmPairwise(salt []byte) *SubjectIdentifierAlgorithmPairwise {
	return &SubjectIdentifierAlgorithmPairwise{Salt: salt}
}

func (g *SubjectIdentifierAlgorithmPairwise) Obfuscate(subject string, client *client.Client) (string, error) {
	// sub = SHA-256 ( sector_identifier || local_account_id || salt ).
	var id string
	if len(client.SectorIdentifierURI) == 0 && len(client.RedirectURIs) > 1 {
		return "", errorsx.WithStack(fosite.ErrInvalidRequest.WithHintf("OAuth 2.0 Client %s has multiple redirect_uris but no sector_identifier_uri was set which is not allowed when performing using subject type pairwise. Please reconfigure the OAuth 2.0 client properly.", client.GetID()))
	} else if len(client.SectorIdentifierURI) == 0 && len(client.RedirectURIs) == 0 {
		return "", errorsx.WithStack(fosite.ErrInvalidRequest.WithHintf("OAuth 2.0 Client %s neither specifies a sector_identifier_uri nor a redirect_uri which is not allowed when performing using subject type pairwise. Please reconfigure the OAuth 2.0 client properly.", client.GetID()))
	} else if len(client.SectorIdentifierURI) > 0 {
		id = client.SectorIdentifierURI
	} else {
		redirectURL, err := url.Parse(client.RedirectURIs[0])
		if err != nil {
			return "", errorsx.WithStack(err)
		}
		id = redirectURL.Host
	}

	return fmt.Sprintf("%x", sha256.Sum256(append(append([]byte{}, []byte(id+subject)...), g.Salt...))), nil
}
