// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestSimpleTranslation(t *testing.T) {
	catalog := NewDefaultMessageCatalog([]*DefaultLocaleBundle{
		{
			LangTag: "en",
			Messages: []*DefaultMessage{
				{
					ID:               "badRequestMethod",
					FormattedMessage: "HTTP method is '%s', expected 'POST'.",
				},
				{
					ID:               "badRequestBody",
					FormattedMessage: "Unable to parse HTTP body, make sure to send a properly formatted form request body.",
				},
			},
		},
		{
			LangTag: "es",
			Messages: []*DefaultMessage{
				{
					ID:               "badRequestMethod",
					FormattedMessage: "El método HTTP es '%s', esperado 'POST'.",
				},
				{
					ID:               "badRequestBody",
					FormattedMessage: "No se puede analizar el cuerpo HTTP, asegúrese de enviar un cuerpo de solicitud de formulario con el formato adecuado.",
				},
			},
		},
	})

	msg := GetMessage(catalog, "badRequestMethod", language.Spanish, "GET")
	assert.Equal(t, msg, "El método HTTP es 'GET', esperado 'POST'.")

	msg = GetMessage(catalog, "badRequestBody", language.English, "GET")
	assert.Equal(t, msg, "Unable to parse HTTP body, make sure to send a properly formatted form request body.")
}
