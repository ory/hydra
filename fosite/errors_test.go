// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/ory/hydra/v2/fosite/i18n"
)

func TestRFC6749Error(t *testing.T) {
	t.Run("case=wrap", func(t *testing.T) {
		orig := errors.New("hi")
		wrap := new(RFC6749Error)
		wrap.Wrap(orig)

		assert.EqualValues(t, orig.(stackTracer).StackTrace(), wrap.StackTrace())
	})

	t.Run("case=wrap_self", func(t *testing.T) {
		wrap := new(RFC6749Error)
		wrap.Wrap(wrap)

		assert.Empty(t, wrap.StackTrace())
	})
}

func TestErrorI18N(t *testing.T) {
	catalog := i18n.NewDefaultMessageCatalog([]*i18n.DefaultLocaleBundle{
		{
			LangTag: "en",
			Messages: []*i18n.DefaultMessage{
				{
					ID:               "access_denied",
					FormattedMessage: "The resource owner or authorization server denied the request.",
				},
				{
					ID:               "badRequestMethod",
					FormattedMessage: "HTTP method is '%s', expected 'POST'.",
				},
			},
		},
		{
			LangTag: "es",
			Messages: []*i18n.DefaultMessage{
				{
					ID:               "access_denied",
					FormattedMessage: "El propietario del recurso o el servidor de autorización denegó la solicitud.",
				},
				{
					ID:               "HTTP method is '%s', expected 'POST'.",
					FormattedMessage: "El método HTTP es '%s', esperado 'POST'.",
				},
				{
					ID:               "Unable to parse HTTP body, make sure to send a properly formatted form request body.",
					FormattedMessage: "No se puede analizar el cuerpo HTTP, asegúrese de enviar un cuerpo de solicitud de formulario con el formato adecuado.",
				},
				{
					ID:               "badRequestMethod",
					FormattedMessage: "El método HTTP es '%s', esperado 'POST'.",
				},
			},
		},
	})

	t.Run("case=legacy", func(t *testing.T) {
		err := ErrAccessDenied.WithLocalizer(catalog, language.Spanish).WithHintf("HTTP method is '%s', expected 'POST'.", "GET")
		assert.EqualValues(t, "El propietario del recurso o el servidor de autorización denegó la solicitud. El método HTTP es 'GET', esperado 'POST'.", err.GetDescription())
	})

	t.Run("case=unsupported_locale_legacy", func(t *testing.T) {
		err := ErrAccessDenied.WithLocalizer(catalog, language.Afrikaans).WithHintf("HTTP method is '%s', expected 'POST'.", "GET")
		assert.EqualValues(t, "The resource owner or authorization server denied the request. HTTP method is 'GET', expected 'POST'.", err.GetDescription())
	})

	t.Run("case=simple", func(t *testing.T) {
		err := ErrAccessDenied.WithLocalizer(catalog, language.Spanish).WithHintIDOrDefaultf("badRequestMethod", "HTTP method is '%s', expected 'POST'.", "GET")
		assert.EqualValues(t, "El propietario del recurso o el servidor de autorización denegó la solicitud. El método HTTP es 'GET', esperado 'POST'.", err.GetDescription())
	})

	t.Run("case=unsupported_locale", func(t *testing.T) {
		err := ErrAccessDenied.WithLocalizer(catalog, language.Afrikaans).WithHintIDOrDefaultf("badRequestMethod", "HTTP method is '%s', expected 'POST'.", "GET")
		assert.EqualValues(t, "The resource owner or authorization server denied the request. HTTP method is 'GET', expected 'POST'.", err.GetDescription())
	})
}
