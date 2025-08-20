// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"crypto/sha256"
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/client"
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
		return "", errors.WithStack(fosite.ErrInvalidRequest.WithHintf("OAuth 2.0 Client %s has multiple redirect_uris but no sector_identifier_uri was set which is not allowed when performing using subject type pairwise. Please reconfigure the OAuth 2.0 client properly.", client.GetID()))
	} else if len(client.SectorIdentifierURI) == 0 && len(client.RedirectURIs) == 0 {
		return "", errors.WithStack(fosite.ErrInvalidRequest.WithHintf("OAuth 2.0 Client %s neither specifies a sector_identifier_uri nor a redirect_uri which is not allowed when performing using subject type pairwise. Please reconfigure the OAuth 2.0 client properly.", client.GetID()))
	} else if len(client.SectorIdentifierURI) > 0 {
		id = client.SectorIdentifierURI
	} else {
		redirectURL, err := url.Parse(client.RedirectURIs[0])
		if err != nil {
			return "", errors.WithStack(err)
		}
		id = redirectURL.Host
	}

	return fmt.Sprintf("%x", sha256.Sum256(append(append([]byte{}, []byte(id+subject)...), g.Salt...))), nil
}
