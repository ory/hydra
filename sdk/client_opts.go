package sdk

import (
	"io/ioutil"
	"net/url"

	"gopkg.in/yaml.v1"
)

// ClusterURL sets Hydra service URL
func ClusterURL(urlStr string) option {
	return func(c *Client) error {
		var err error
		c.clusterURL, err = url.Parse(urlStr)
		return err
	}
}

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

// ClientID sets OAuth client ID
func ClientID(id string) option {
	return func(c *Client) error {
		c.clientID = id
		return nil
	}
}

// ClientSecret sets OAuth client secret
func ClientSecret(secret string) option {
	return func(c *Client) error {
		c.clientSecret = secret
		return nil
	}
}

// SkipTLSVerify skips TLS verification
func SkipTLSVerify() option {
	return func(c *Client) error {
		c.skipTLSVerify = true
		return nil
	}
}

// Scopes sets client scopes granted by Hydra
func Scopes(scopes ...string) option {
	return func(c *Client) error {
		c.scopes = scopes
		return nil
	}
}
