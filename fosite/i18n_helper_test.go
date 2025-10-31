// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/ory/hydra/v2/fosite/i18n"
)

func TestErrorTranslation(t *testing.T) {
	catalog := i18n.NewDefaultMessageCatalog([]*i18n.DefaultLocaleBundle{
		{
			LangTag: "en",
			Messages: []*i18n.DefaultMessage{
				{
					ID:               "badRequestMethod",
					FormattedMessage: "HTTP method is '%s', expected 'POST'.",
				},
				{
					ID:               "invalid_request",
					FormattedMessage: "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed.",
				},
			},
		},
		{
			LangTag: "es",
			Messages: []*i18n.DefaultMessage{
				{
					ID:               "badRequestMethod",
					FormattedMessage: "El método HTTP es '%s', esperado 'POST'.",
				},
				{
					ID:               "invalid_request",
					FormattedMessage: "A la solicitud le falta un parámetro obligatorio, incluye un valor de parámetro no válido, incluye un parámetro más de una vez o tiene un formato incorrecto.",
				},
			},
		},
	})

	errWithNoCatalog := ErrInvalidRequest.WithHintIDOrDefaultf("badRequestMethod", "HTTP method is '%s', expected 'POST'.", "GET")
	errWithCatalog := errWithNoCatalog.WithLocalizer(catalog, language.Spanish)

	assert.Equal(t, "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. HTTP method is 'GET', expected 'POST'.",
		errWithNoCatalog.GetDescription(), "Message does not match when no catalog is specified")
	assert.Equal(t, "A la solicitud le falta un parámetro obligatorio, incluye un valor de parámetro no válido, incluye un parámetro más de una vez o tiene un formato incorrecto. El método HTTP es 'GET', esperado 'POST'.",
		errWithCatalog.GetDescription(), "Message does not match when catalog is specified")
}
