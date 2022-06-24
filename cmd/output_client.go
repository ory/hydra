package cmd

import (
	"fmt"
	"strings"

	hydra "github.com/ory/hydra-client-go"
)

type (
	outputOAuth2Client           hydra.OAuth2Client
	outputOAuth2ClientCollection struct {
		clients []hydra.OAuth2Client
	}
)

func (_ outputOAuth2Client) Header() []string {
	return []string{"CLIENT ID", "CLIENT SECRET", "GRANT TYPES", "RESPONSE TYPES", "SCOPE", "AUDIENCE", "REDIRECT URIS"}
}

func (i outputOAuth2Client) Columns() []string {
	data := [7]string{
		fmt.Sprintf("%v", i.ClientId),
		fmt.Sprintf("%v", i.ClientSecret),
		strings.Join(i.GrantTypes, ", "),
		strings.Join(i.ResponseTypes, ", "),
		fmt.Sprintf("%v", i.Scope),
		strings.Join(i.Audience, ","),
		strings.Join(i.RedirectUris, ", "),
	}
	return data[:]
}

func (i outputOAuth2Client) Interface() interface{} {
	return i
}

func (_ outputOAuth2ClientCollection) Header() []string {
	return outputOAuth2Client{}.Header()
}

func (c outputOAuth2ClientCollection) Table() [][]string {
	rows := make([][]string, len(c.clients))
	for i, client := range c.clients {
		rows[i] = outputOAuth2Client(client).Columns()
	}
	return rows
}

func (c outputOAuth2ClientCollection) Interface() interface{} {
	return c.clients
}

func (c outputOAuth2ClientCollection) Len() int {
	return len(c.clients)
}
