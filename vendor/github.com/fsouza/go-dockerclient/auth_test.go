// Copyright 2015 go-dockerclient authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docker

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
)

func TestAuthConfigurationSearchPath(t *testing.T) {
	var testData = []struct {
		dockerConfigEnv string
		homeEnv         string
		expectedPaths   []string
	}{
		{"", "", []string{}},
		{"", "home", []string{path.Join("home", ".docker", "config.json"), path.Join("home", ".dockercfg")}},
		{"docker_config", "", []string{path.Join("docker_config", "config.json")}},
		{"a", "b", []string{path.Join("a", "config.json"), path.Join("b", ".docker", "config.json"), path.Join("b", ".dockercfg")}},
	}
	for _, tt := range testData {
		paths := cfgPaths(tt.dockerConfigEnv, tt.homeEnv)
		if got, want := strings.Join(paths, ","), strings.Join(tt.expectedPaths, ","); got != want {
			t.Errorf("cfgPaths: wrong result. Want: %s. Got: %s", want, got)
		}
	}
}

func TestAuthConfigurationsFromFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "go-dockerclient-auth-test")
	if err != nil {
		t.Errorf("Unable to create temporary directory for TestAuthConfigurationsFromFile: %s", err)
	}
	defer os.RemoveAll(tmpDir)
	authString := base64.StdEncoding.EncodeToString([]byte("user:pass"))
	content := fmt.Sprintf("{\"auths\":{\"foo\": {\"auth\": \"%s\"}}}", authString)
	configFile := path.Join(tmpDir, "docker_config")
	if err = ioutil.WriteFile(configFile, []byte(content), 0600); err != nil {
		t.Errorf("Error writing auth config for TestAuthConfigurationsFromFile: %s", err)
	}
	auths, err := NewAuthConfigurationsFromFile(configFile)
	if err != nil {
		t.Errorf("Error calling NewAuthConfigurationsFromFile: %s", err)
	}
	if _, hasKey := auths.Configs["foo"]; !hasKey {
		t.Errorf("Returned auths did not include expected auth key foo")
	}
}

func TestAuthLegacyConfig(t *testing.T) {
	auth := base64.StdEncoding.EncodeToString([]byte("user:pa:ss"))
	read := strings.NewReader(fmt.Sprintf(`{"docker.io":{"auth":"%s","email":"user@example.com"}}`, auth))
	ac, err := NewAuthConfigurations(read)
	if err != nil {
		t.Error(err)
	}
	c, ok := ac.Configs["docker.io"]
	if !ok {
		t.Error("NewAuthConfigurations: Expected Configs to contain docker.io")
	}
	if got, want := c.Email, "user@example.com"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].Email: wrong result. Want %q. Got %q`, want, got)
	}
	if got, want := c.Username, "user"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].Username: wrong result. Want %q. Got %q`, want, got)
	}
	if got, want := c.Password, "pa:ss"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].Password: wrong result. Want %q. Got %q`, want, got)
	}
	if got, want := c.ServerAddress, "docker.io"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].ServerAddress: wrong result. Want %q. Got %q`, want, got)
	}
}

func TestAuthBadConfig(t *testing.T) {
	auth := base64.StdEncoding.EncodeToString([]byte("userpass"))
	read := strings.NewReader(fmt.Sprintf(`{"docker.io":{"auth":"%s","email":"user@example.com"}}`, auth))
	ac, err := NewAuthConfigurations(read)
	if err != ErrCannotParseDockercfg {
		t.Errorf("Incorrect error returned %v\n", err)
	}
	if ac != nil {
		t.Errorf("Invalid auth configuration returned, should be nil %v\n", ac)
	}
}

func TestAuthAndOtherFields(t *testing.T) {
	auth := base64.StdEncoding.EncodeToString([]byte("user:pass"))
	read := strings.NewReader(fmt.Sprintf(`{
		"auths":{"docker.io":{"auth":"%s","email":"user@example.com"}},
		"detachKeys": "ctrl-e,e",
		"HttpHeaders": { "MyHeader": "MyValue" }}`, auth))

	ac, err := NewAuthConfigurations(read)
	if err != nil {
		t.Error(err)
	}
	c, ok := ac.Configs["docker.io"]
	if !ok {
		t.Error("NewAuthConfigurations: Expected Configs to contain docker.io")
	}
	if got, want := c.Email, "user@example.com"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].Email: wrong result. Want %q. Got %q`, want, got)
	}
	if got, want := c.Username, "user"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].Username: wrong result. Want %q. Got %q`, want, got)
	}
	if got, want := c.Password, "pass"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].Password: wrong result. Want %q. Got %q`, want, got)
	}
	if got, want := c.ServerAddress, "docker.io"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].ServerAddress: wrong result. Want %q. Got %q`, want, got)
	}
}
func TestAuthConfig(t *testing.T) {
	auth := base64.StdEncoding.EncodeToString([]byte("user:pass"))
	read := strings.NewReader(fmt.Sprintf(`{"auths":{"docker.io":{"auth":"%s","email":"user@example.com"}}}`, auth))
	ac, err := NewAuthConfigurations(read)
	if err != nil {
		t.Error(err)
	}
	c, ok := ac.Configs["docker.io"]
	if !ok {
		t.Error("NewAuthConfigurations: Expected Configs to contain docker.io")
	}
	if got, want := c.Email, "user@example.com"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].Email: wrong result. Want %q. Got %q`, want, got)
	}
	if got, want := c.Username, "user"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].Username: wrong result. Want %q. Got %q`, want, got)
	}
	if got, want := c.Password, "pass"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].Password: wrong result. Want %q. Got %q`, want, got)
	}
	if got, want := c.ServerAddress, "docker.io"; got != want {
		t.Errorf(`AuthConfigurations.Configs["docker.io"].ServerAddress: wrong result. Want %q. Got %q`, want, got)
	}
}

func TestAuthCheck(t *testing.T) {
	fakeRT := &FakeRoundTripper{status: http.StatusOK}
	client := newTestClient(fakeRT)
	if _, err := client.AuthCheck(nil); err == nil {
		t.Fatalf("expected error on nil auth config")
	}
	// test good auth
	if _, err := client.AuthCheck(&AuthConfiguration{}); err != nil {
		t.Fatal(err)
	}
	*fakeRT = FakeRoundTripper{status: http.StatusUnauthorized}
	if _, err := client.AuthCheck(&AuthConfiguration{}); err == nil {
		t.Fatal("expected failure from unauthorized auth")
	}
}
