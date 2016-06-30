package sdk

import (
	"io/ioutil"
	"net/url"

	"gopkg.in/yaml.v1"
)

func ClusterURL(urlStr string) option {
	return func(c *Client) error {
		var err error
		c.clusterURL, err = url.Parse(urlStr)
		return err
	}
}

func FromYAML(file string) option {
	return func(c *Client) error {
		var err error
		var config = struct {
			ClusterURL   string `yaml:"cluster_url"`
			ClientID     string `yaml:"client_id"`
			ClientSecret string `yaml:"client_secret"`
		}{}

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

func ClientID(id string) option {
	return func(c *Client) error {
		c.clientID = id
		return nil
	}
}

func ClientSecret(secret string) option {
	return func(c *Client) error {
		c.clientSecret = secret
		return nil
	}
}

func SkipSSL() option {
	return func(c *Client) error {
		c.skipSSL = true
		return nil
	}
}

func Scopes(scopes ...string) option {
	return func(c *Client) error {
		c.scopes = scopes
		return nil
	}
}
