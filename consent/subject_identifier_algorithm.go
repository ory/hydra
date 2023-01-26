// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import "github.com/ory/hydra/v2/client"

type SubjectIdentifierAlgorithm interface {
	// Obfuscate derives a pairwise subject identifier from the given string.
	Obfuscate(subject string, client *client.Client) (string, error)
}
