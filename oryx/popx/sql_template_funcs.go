// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package popx

import (
	"regexp"

	"github.com/pkg/errors"
)

var SQLTemplateFuncs = map[string]interface{}{
	"identifier": Identifier,
}

var identifierPattern = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")

func Identifier(i string) (string, error) {
	if !identifierPattern.MatchString(i) {
		return "", errors.Errorf("invalid SQL identifier '%s'", i)
	}
	return i, nil
}
