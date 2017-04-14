package sdk

import (
	"io/ioutil"
	"net/url"

	"gopkg.in/yaml.v2"
)

type hydraConfig struct {
	ClusterURL   string `yaml:"cluster_url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

// FromYAML loads configurations from a YAML file
func FromYAML(file string) option {
	return func(c *Client) error {
		var err error
		var config = hydraConfig{}

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return err
		}

		c.clusterURL, err = url.Parse(config.ClusterURL)
		if err != nil {
			return err
		}

		c.clientID = config.ClientID
		c.clientSecret = config.ClientSecret

		return nil
	}
}

// ClusterURL sets Hydra service URL
//
//  var hydra, err = sdk.Connect(
// 	sdk.ClientID("https://localhost:1234/"),
//  )
func ClusterURL(urlStr string) option {
	return func(c *Client) error {
		var err error
		c.clusterURL, err = url.Parse(urlStr)
		return err
	}
}

// ClientID sets the OAuth2 Client ID.
//
//  var hydra, err = sdk.Connect(
// 	sdk.ClientID("client-id"),
//  )
func ClientID(id string) option {
	return func(c *Client) error {
		c.clientID = id
		return nil
	}
}

// ClientSecret sets OAuth2 Client secret.
//
//  var hydra, err = sdk.Connect(
// 	sdk.ClientSecret("client-secret"),
//  )
func ClientSecret(secret string) option {
	return func(c *Client) error {
		c.clientSecret = secret
		return nil
	}
}

// SkipTLSVerify skips TLS verification for HTTPS connections.
//
//  var hydra, err = sdk.Connect(
// 	sdk.SkipTLSVerify(),
//  )
func SkipTLSVerify() option {
	return func(c *Client) error {
		c.skipTLSVerify = true
		return nil
	}
}

// Scopes is a list of scopes that are requested in the client credentials grant.
//
//  var hydra, err = sdk.Connect(
//  	sdk.Scopes("foo", "bar"),
//  )
func Scopes(scopes ...string) option {
	return func(c *Client) error {
		c.scopes = scopes
		return nil
	}
}
