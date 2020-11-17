// +build conformity

package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"testing"
	"time"

	"github.com/ory/x/stringslice"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/ory/x/urlx"
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
	server    = urlx.ParseOrPanic("https://localhost:8443")
	config, _ = ioutil.ReadFile("./config.json")
	client    = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

func waitForServices(t *testing.T) {
	var conformOk, hydraOk bool
	start := time.Now()
	for {
		res, err := client.Get(server.String())
		conformOk = err == nil && res.StatusCode == 200

		res, err = client.Get("https://127.0.0.1:4444/health/ready")
		hydraOk = err == nil && res.StatusCode == 200

		if conformOk && hydraOk {
			break
		}

		if time.Since(start).Minutes() > 2 {
			t.Fatalf("Waiting for service exceeded timeout of two minutes.")
		}

		t.Logf("Waiting for deployments to come alive...")
		time.Sleep(time.Second)
	}
}

func TestPlans(t *testing.T) {
	waitForServices(t)

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
	res, err := client.Post(href, "application/json", payload)
	require.NoError(t, err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
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

		t.Run("testModule="+module, func(t *testing.T) {
			if isParallel {
				t.Parallel()
			}

			if module == "oidcc-server-rotate-keys" && isParallel {
				t.Skipf("Test module 'oidcc-server-rotate-keys' can not run in parallel tests and was skipped...")
				return
			} else if module != "oidcc-server-rotate-keys" && !isParallel {
				t.Skipf("Without paralleism only test module 'oidcc-server-rotate-keys' will be executed.")
				return
			}

			body = makePost(t, urlx.CopyWithQuery(
				urlx.AppendPaths(server, "/api/runner"),
				url.Values{"test": {module}, "plan": {plan}, "variant": {v.Get("variant").Raw}}).String(), nil, 201)

			for {
				time.Sleep(time.Millisecond * 100)

				state, passed := checkStatus(t, gjson.GetBytes(body, "id").String())
				if passed {
					break
				}

				switch module {
				case "oidcc-server-rotate-keys":
					if state == "CONFIGURED" {
						t.Logf("Rotating ID Token keys....")
						cmd := exec.Command("./hydra/hydra", "keys", "create", "--endpoint=https://127.0.0.1:4445/", "hydra.openid.id-token", "-a", "RS256", "--skip-tls-verify")
						var buf bytes.Buffer
						cmd.Stderr = &buf
						cmd.Stdout = &buf
						require.NoError(t, cmd.Run(), "%s", buf.String())

						makePost(t, urlx.AppendPaths(server, "/api/runner/", gjson.GetBytes(body, "id").String()).String(), nil, 200)
					}
				}
			}
		})

		return true
	})
}

func checkStatus(t *testing.T, testID string) (state string, passed bool) {
	res, err := client.Get(urlx.AppendPaths(server, "/api/info", testID).String())
	require.NoError(t, err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode, "%s", body)

	state = gjson.GetBytes(body, "status").String()
	switch state {
	case "INTERRUPTED":
		t.Fatalf("Status returned was INTERRUPTED: %s", body)
		return
	case "FINISHED":
		status := gjson.GetBytes(body, "result").String()

		if status == "PASSED" || status == "WARNING" || status == "SKIPPED" || status == "REVIEW" {
			return state, true
		} else if status == "FAILED" {
			t.Fatalf("Expected status not to be FAILED got: %s", body)
			return
		}

		t.Fatalf("Unexpected status: %s", body)
		return
	case "CONFIGURED":
		fallthrough
	case "CREATED":
		fallthrough
	case "RUNNING":
		fallthrough
	case "WAITING":
		return state, false
	}

	t.Fatalf("Unexpected state: %s", body)
	return
}
