package sdk

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestClusterURLOption(t *testing.T) {
	c := &Client{}
	expected := "https://localhost:4444"
	err := ClusterURL(expected)(c)

	assert.Nil(t, err)
	assert.Equal(t, expected, c.clusterURL.String())
}

func TestClientIDOption(t *testing.T) {
	c := &Client{}
	expected := "deadbeef-dead-beef-beef"
	err := ClientID(expected)(c)

	assert.Nil(t, err)
	assert.Equal(t, expected, c.clientID)
}

func TestClientSecretOption(t *testing.T) {
	c := &Client{}

	expected := "secret"
	err := ClientSecret(expected)(c)

	assert.Nil(t, err)
	assert.Equal(t, expected, c.clientSecret)
}

func TestSkipSSLOption(t *testing.T) {
	c := &Client{}

	err := SkipTLSVerify()(c)

	assert.Nil(t, err)
	assert.Equal(t, true, c.skipTLSVerify)
}

func TestScopesOption(t *testing.T) {
	c := &Client{}

	expected := []string{"a", "b", "c"}
	err := Scopes(expected...)(c)

	assert.Nil(t, err)
	assert.Equal(t, expected, c.scopes)
}

func TestFromYAMLOption(t *testing.T) {
	c := &Client{}

	conf := &hydraConfig{
		ClusterURL:   "https://localhost:4444",
		ClientID:     "1cfe0a5e-2533-4312-9e74-128b5aab4431",
		ClientSecret: "Q6&u=iwvTPh8r)Ar",
	}

	tmpFile, err := ioutil.TempFile("", "hydra_sdk")
	assert.Nil(t, err)

	fileContent, err := yaml.Marshal(conf)
	assert.Nil(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.Write(fileContent)
	assert.Nil(t, err)

	err = FromYAML(tmpFile.Name())(c)
	assert.Nil(t, err)

	assert.Equal(t, conf.ClusterURL, c.clusterURL.String())
	assert.Equal(t, conf.ClientID, c.clientID)
	assert.Equal(t, conf.ClientSecret, c.clientSecret)
}
