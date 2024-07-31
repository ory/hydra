// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

//go:build conformity
// +build conformity

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	backoff "github.com/cenkalti/backoff/v3"

	hydrac "github.com/ory/hydra-client-go/v2"

	"github.com/ory/x/httpx"

	"github.com/ory/x/stringslice"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/x/urlx"
)

type status int

const (
	statusFailed status = iota
	statusRetry
	statusRunning
	statusSuccess
)

var (
	skipWhenShort = []string{"oidcc-test-plan"}

	plans = []url.Values{
		{"planName": {"oidcc-formpost-implicit-certification-test-plan"}, "variant": {"{\"server_metadata\":\"discovery\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-formpost-basic-certification-test-plan"}, "variant": {"{\"server_metadata\":\"discovery\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-formpost-hybrid-certification-test-plan"}, "variant": {"{\"server_metadata\":\"discovery\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-hybrid-certification-test-plan"}, "variant": {"{\"server_metadata\":\"discovery\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-implicit-certification-test-plan"}, "variant": {"{\"server_metadata\":\"discovery\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-dynamic-certification-test-plan"}, "variant": {"{\"response_type\":\"code\"}"}},
		{"planName": {"oidcc-dynamic-certification-test-plan"}, "variant": {"{\"response_type\":\"id_token\"}"}},
		{"planName": {"oidcc-dynamic-certification-test-plan"}, "variant": {"{\"response_type\":\"id_token token\"}"}},
		{"planName": {"oidcc-dynamic-certification-test-plan"}, "variant": {"{\"response_type\":\"code id_token\"}"}},
		{"planName": {"oidcc-dynamic-certification-test-plan"}, "variant": {"{\"response_type\":\"code token\"}"}},
		{"planName": {"oidcc-dynamic-certification-test-plan"}, "variant": {"{\"response_type\":\"code id_token token\"}"}},
		{"planName": {"oidcc-config-certification-test-plan"}},

		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"code\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"id_token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"id_token token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"code id_token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"code token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"code id_token token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"code\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"id_token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"id_token token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"code id_token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"code token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"client_secret_basic\",\"response_type\":\"code id_token token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},

		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"code\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"id_token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"id_token token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"code id_token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"code token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"code id_token token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"code\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"id_token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"id_token token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"code id_token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"code token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
		{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"private_key_jwt\",\"response_type\":\"code id_token token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},

		/*
			See https://gitlab.com/openid/conformance-suite/-/issues/856

			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"code\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"id_token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"id_token token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"code id_token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"code token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"code id_token token\",\"response_mode\":\"default\",\"client_registration\":\"dynamic_client\"}"}},
			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"code\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"id_token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"id_token token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"code id_token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"code token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
			{"planName": {"oidcc-test-plan"}, "variant": {"{\"client_auth_type\":\"none\",\"response_type\":\"code id_token token\",\"response_mode\":\"form_post\",\"client_registration\":\"dynamic_client\"}"}},
		*/

		{"planName": {"oidcc-formpost-basic-certification-test-plan"}, "variant": {"{\"server_metadata\":\"discovery\",\"client_registration\":\"dynamic_client\"}"}},
	}
	server     = urlx.ParseOrPanic("https://127.0.0.1:8443")
	config, _  = os.ReadFile("./config.json")
	httpClient = httpx.NewResilientClient(httpx.ResilientClientWithMinxRetryWait(time.Second * 5))
	workdir    string

	hydra = hydrac.NewAPIClient(hydrac.NewConfiguration())
)

func init() {
	httpClient.HTTPClient.Timeout = 5 * time.Second
	httpClient.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	hydra.GetConfig().HTTPClient = httpClient.HTTPClient
	hydra.GetConfig().Servers = hydrac.ServerConfigurations{{URL: "https://127.0.0.1:4445"}}
}

func waitForServices(t *testing.T) {
	var conformOk, hydraOk bool
	start := time.Now()
	for {
		server := server.String()
		res, err := httpClient.Get(server)
		conformOk = err == nil && res.StatusCode == 200
		t.Logf("Checking %s (%v): %s (%+v)", server, conformOk, err, res)

		server = "https://127.0.0.1:4444/health/ready"
		res, err = httpClient.Get(server)
		hydraOk = err == nil && res.StatusCode == 200
		t.Logf("Checking %s (%v): %s (%+v)", server, hydraOk, err, res)

		if conformOk && hydraOk {
			break
		}

		if time.Since(start).Minutes() > 2 {
			require.FailNow(t, "Waiting for service exceeded timeout of two minutes.")
		}

		t.Logf("Waiting for deployments to come alive...")
		time.Sleep(time.Second)
	}
}

func TestPlans(t *testing.T) {
	waitForServices(t)

	var err error
	workdir, err = filepath.Abs("../../")
	require.NoError(t, err)

	t.Run("parallel=true", func(t *testing.T) {
		for k := range plans {
			plan := plans[k]
			t.Run(fmt.Sprintf("plan=%s", plan), func(t *testing.T) {
				t.Parallel()
				createPlan(t, plan, true)
			})
		}
	})

	t.Run("parallel=false", func(t *testing.T) {
		// Run remaining tests which do not work when parallelism is active
		for _, plan := range plans {
			t.Run(fmt.Sprintf("plan=%s", plan), func(t *testing.T) {
				createPlan(t, plan, false)
			})
		}
	})
}

func makePost(t *testing.T, href string, payload io.Reader, esc int) []byte {
	res, err := httpClient.Post(href, "application/json", payload)
	require.NoError(t, err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, esc, res.StatusCode, "%s\n%s", href, body)
	return body
}

func createPlan(t *testing.T, extra url.Values, isParallel bool) {
	planName := extra.Get("planName")
	if stringslice.Has(skipWhenShort, planName) && testing.Short() {
		t.Skipf("Skipping test plan '%s' because short tests", planName)
		return
	}

	// https://localhost:8443/api/plan?planName=oidcc-formpost-basic-certification-test-plan&variant={"server_metadata":"discovery","client_registration":"dynamic_client"}&variant={"server_metadata":"discovery","client_registration":"dynamic_client"}
	//planConfig, err := sjson.SetBytes(config, "alias", uuid.New())
	//require.NoError(t, err)
	body := makePost(t, urlx.CopyWithQuery(urlx.AppendPaths(server, "/api/plan"), extra).String(),
		bytes.NewReader(config),
		201)

	plan := gjson.GetBytes(body, "id").String()
	require.NotEmpty(t, plan)

	t.Logf("Created plan: %s", plan)
	gjson.GetBytes(body, "modules").ForEach(func(_, v gjson.Result) bool {
		module := v.Get("testModule").String()

		t.Logf("Running testModule %s for plan %s", module, plan)
		t.Run("testModule="+module, func(t *testing.T) {
			if isParallel {
				t.Parallel()
			}

			if module == "oidcc-server-rotate-keys" && isParallel {
				t.Skipf("Test module 'oidcc-server-rotate-keys' can not run in parallel tests and was skipped...")
				return
			} else if module != "oidcc-server-rotate-keys" && !isParallel {
				t.Skipf("Without parallelism only test module 'oidcc-server-rotate-keys' will be executed.")
				return
			}

			params := url.Values{"test": {module}, "plan": {plan}, "variant": {v.Get("variant").Raw}}

			const maxRetries = 5
			for retry := 1; retry <= maxRetries; retry++ {
				time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)

				t.Logf("Creating retry %d/%d testModule %s for plan %s with params: %+v", retry, maxRetries, module, plan, params)
				body := makePost(t, urlx.CopyWithQuery(urlx.AppendPaths(server, "/api/runner"), params).String(),
					nil, 201)

				conf := backoff.NewExponentialBackOff()
				conf.MaxElapsedTime = time.Minute * 5
				conf.MaxInterval = time.Second * 5
				conf.InitialInterval = time.Second

				for {
					nb := conf.NextBackOff()
					if nb == backoff.Stop {
						t.Logf("Waited %.2f minutes for a status change for testModule %s for plan %s but received none. Retrying with a fresh test...", conf.MaxElapsedTime.Minutes(), module, plan)
						break
					}
					time.Sleep(nb)

					state, passed := checkStatus(t, gjson.GetBytes(body, "id").String())
					switch passed {
					case statusRetry:
						t.Logf("Status from testModule %s for plan %s with params marked the test for retry. Retrying with a fresh test...", module, plan)
						break
					case statusFailed:
						panic("This statement should never be reached")
					case statusSuccess:
						return
					}

					switch module {
					case "oidcc-server-rotate-keys":
						if state == "CONFIGURED" {
							t.Logf("Rotating ID Token keys....")

							conf := backoff.NewExponentialBackOff()
							conf.MaxElapsedTime = time.Minute * 5
							conf.MaxInterval = time.Second * 5
							conf.InitialInterval = time.Second
							var err error

							for {
								bo := conf.NextBackOff()
								require.NotEqual(t, backoff.Stop, bo, "%+v", err)

								_, _, err = hydra.JwkAPI.CreateJsonWebKeySet(context.Background(), "hydra.openid.id-token").CreateJsonWebKeySet(hydrac.CreateJsonWebKeySet{
									Alg: "RS256",
								}).Execute()
								if err == nil {
									break
								}

								time.Sleep(bo)
							}

							makePost(t, urlx.AppendPaths(server, "/api/runner/", gjson.GetBytes(body, "id").String()).String(), nil, 200)
						}
					}
				}
			}
			require.FailNowf(t, "Retries exceeded", "Exceeded maximum retries %d for test %s in plan %s", maxRetries, module, plan)
		})

		return true
	})
}

func checkStatus(t *testing.T, testID string) (string, status) {
	res, err := httpClient.Get(urlx.AppendPaths(server, "/api/info", testID).String())
	require.NoError(t, err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode, "%s", body)

	state := gjson.GetBytes(body, "status").String()
	t.Logf("Got status %s for %s", state, testID)
	switch state {
	case "INTERRUPTED":
		t.Logf("Test was INTERRUPTED: %s", body)
		return state, statusRetry
	case "FINISHED":
		result := gjson.GetBytes(body, "result").String()
		t.Logf("Got result %s for %s", result, testID)

		if result == "PASSED" || result == "WARNING" || result == "SKIPPED" || result == "REVIEW" {
			return state, statusSuccess
		} else if result == "FAILED" {
			require.FailNowf(t, "Test was FAILED", "Expected status not to be FAILED got: %s", body)
			return state, statusFailed
		}

		require.FailNowf(t, "Test failed with another error", "Unexpected status: %s", body)
		return state, statusFailed
	case "CONFIGURED":
		fallthrough
	case "CREATED":
		fallthrough
	case "RUNNING":
		fallthrough
	case "WAITING":
		return state, statusRunning
	}

	require.FailNowf(t, "Unexpected state", "Unexpected state: %s", body)
	return state, statusFailed
}
