// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import "github.com/ory/hydra/v2/client"

type SubjectIdentifierAlgorithmPublic struct{}

func NewSubjectIdentifierAlgorithmPublic() *SubjectIdentifierAlgorithmPublic {
	return &SubjectIdentifierAlgorithmPublic{}
}

func (g *SubjectIdentifierAlgorithmPublic) Obfuscate(subject string, client *client.Client) (string, error) {
	return subject, nil
}
