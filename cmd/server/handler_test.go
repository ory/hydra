package server

import (
	"bytes"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/hydra/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

func TestStart(t *testing.T) {
	router := httprouter.New()
	h := &Handler{
		Config: &config.Config{
			DatabaseURL: "memory",
		},
	}
	h.registerRoutes(router)
}

func TestMiddleware(t *testing.T) {
	nrKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	nrName := os.Getenv("NEW_RELIC_APP_NAME")
	airID := os.Getenv("AIRBRAKE_PROJECT_ID")
	airKey := os.Getenv("AIRBRAKE_PROJECT_KEY")

	t.Run("NewRelic1", testUseNewRelicMiddlewareWithWrongValues)
	t.Run("NewRelic2", testUseNewRelicMiddlewareWithoutValues)
	t.Run("Airbrake1", testUseAirbrakeMiddleware)
	t.Run("Airbrake2", testUseAirbrakeMiddlewareWithoutValues)
	t.Run("Airbrake3", testUseAirbrakeMiddlewareWithoutProperID)

	os.Setenv("NEW_RELIC_LICENSE_KEY", nrKey)
	os.Setenv("NEW_RELIC_APP_NAME", nrName)
	os.Setenv("AIRBRAKE_PROJECT_ID", airID)
	os.Setenv("AIRBRAKE_PROJECT_KEY", airKey)
	logrus.SetOutput(os.Stderr)
}

func TestStatsdRegex(t *testing.T) {
	regx = newStatsdTagsSanitizerRegex()
	resource := "rn:hydra.warden-token_allowed"
	sanitized_resource := regx.ReplaceAllString(resource, "_")
	assert.NotEqual(t, resource, sanitized_resource, "regex does not behave as expected")
	assert.Equal(t, "rn_hydra.warden-token_allowed", sanitized_resource, "regex does not behave as expected")
}

func testUseNewRelicMiddlewareWithWrongValues(t *testing.T) {
	os.Setenv("NEW_RELIC_LICENSE_KEY", "1")
	os.Setenv("NEW_RELIC_APP_NAME", "1")
	var buffer bytes.Buffer
	logrus.SetOutput(&buffer)
	n := negroni.New()

	useNewRelicMiddleware(n)
	assert.Contains(t, string(buffer.Bytes()), "Error creating New Relic app: license length is not 40")
}

func testUseNewRelicMiddlewareWithoutValues(t *testing.T) {
	os.Setenv("NEW_RELIC_LICENSE_KEY", "")
	var buffer bytes.Buffer
	logrus.SetOutput(&buffer)
	n := negroni.New()

	useNewRelicMiddleware(n)
	assert.Contains(t, string(buffer.Bytes()), "New Relic disabled - configs not found")
}

func testUseAirbrakeMiddleware(t *testing.T) {
	os.Setenv("AIRBRAKE_PROJECT_ID", "1")
	os.Setenv("AIRBRAKE_PROJECT_KEY", "1")
	var buffer bytes.Buffer
	logrus.SetOutput(&buffer)
	n := negroni.New()

	useAirbrakeMiddleware(n)
	assert.Contains(t, string(buffer.Bytes()), "Airbrake enabled!")
}

func testUseAirbrakeMiddlewareWithoutValues(t *testing.T) {
	os.Setenv("AIRBRAKE_PROJECT_ID", "")
	var buffer bytes.Buffer
	logrus.SetOutput(&buffer)
	n := negroni.New()

	useAirbrakeMiddleware(n)
	assert.Contains(t, string(buffer.Bytes()), "Airbrake disabled - configs not found")
}

func testUseAirbrakeMiddlewareWithoutProperID(t *testing.T) {
	os.Setenv("AIRBRAKE_PROJECT_ID", "not a number")
	os.Setenv("AIRBRAKE_PROJECT_KEY", "1")
	var buffer bytes.Buffer
	logrus.SetOutput(&buffer)
	n := negroni.New()

	useAirbrakeMiddleware(n)
	assert.Contains(t, string(buffer.Bytes()), "Airbrake disabled - error parsing airbrake project ID")
}
