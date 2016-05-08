package cmd

import (
	"fmt"
	"net/http"

	"github.com/ory-am/common/pkg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func authenticate() *http.Client {
	oauthConfig := clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     pkg.JoinURL(config.ClusterURL, "oauth2/token"),
		Scopes:       []string{"core", "hydra"},
	}
	token, err := oauthConfig.Token(oauth2.NoContext)
	if err != nil {
		fatal("Unable to retrieve access token. Did you run `hydra connect`?")
	}
	fmt.Printf("Got access token %v", token)

	return oauthConfig.Client(oauth2.NoContext)
}
