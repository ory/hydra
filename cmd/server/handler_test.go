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
	airID := os.Getenv("AIRBRAKE_PROJECT_ID")
	airKey := os.Getenv("AIRBRAKE_PROJECT_KEY")

	t.Run("Airbrake1", testUseAirbrakeMiddleware)
	t.Run("Airbrake2", testUseAirbrakeMiddlewareWithoutValues)
	t.Run("Airbrake3", testUseAirbrakeMiddlewareWithoutProperID)

	os.Setenv("AIRBRAKE_PROJECT_ID", airID)
	os.Setenv("AIRBRAKE_PROJECT_KEY", airKey)
	logrus.SetOutput(os.Stderr)
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
