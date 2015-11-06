package provider

import (
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/RangelReale/osin"
	"github.com/ory-am/hydra/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/ory-am/hydra/Godeps/_workspace/src/golang.org/x/oauth2"
	"testing"
)

func TestGetAuthCodeURL(t *testing.T) {
	conf := oauth2.Config{
		ClientID:     "remote-client",
		ClientSecret: "secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://third-party/oauth2/auth",
			TokenURL: "http://third-party/oauth2/token",
		},
		RedirectURL: "http://hydra/oauth2/auth",
		Scopes:      []string{"scope"},
	}
	ar := &osin.AuthorizeRequest{
		Type:        osin.CODE,
		Client:      &osin.DefaultClient{Id: "client", Secret: "secret"},
		Scope:       "scope",
		RedirectUri: "http://remote/callback",
		State:       "state",
	}
	url := GetAuthCodeURL(conf, ar, "foo")
	assert.Equal(t, "http://third-party/oauth2/auth?client_id=remote-client&redirect_uri=http%3A%2F%2Fhydra%2Foauth2%2Fauth%3Focl%3Dclient%26opr%3Dfoo%26ord%3Dhttp%253A%252F%252Fremote%252Fcallback%26osc%3Dscope%26ost%3Dstate%26otp%3Dcode&response_type=code&scope=scope&state=state", url)
}
