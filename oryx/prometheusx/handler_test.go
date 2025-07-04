// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package prometheusx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"

	"github.com/ory/herodot"
	"github.com/ory/x/logrusx"
	prometheus "github.com/ory/x/prometheusx"

	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	router := httprouter.New()
	l := logrusx.New("Ory X", "test")
	writer := herodot.NewJSONWriter(l)
	metricsHandler := prometheus.NewHandler(writer, "test")
	metricsHandler.SetRoutes(router)
	ts := httptest.NewServer(router)
	defer ts.Close()

	c := http.DefaultClient

	response, err := c.Get(ts.URL + prometheus.MetricsPrometheusPath)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, response.StatusCode)

	textParser := expfmt.TextParser{}
	text, err := textParser.TextToMetricFamilies(response.Body)
	require.NoError(t, err)
	require.EqualValues(t, "go_info", *text["go_info"].Name)
}
