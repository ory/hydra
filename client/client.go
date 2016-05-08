package client

type Client struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Contacts          []string `json:"contacts"`
	Secret            []byte   `json:"secret,omitempty"`
	RedirectURIs      []string `json:"redirect_uris"`
	Owner             string   `json:"owner"`
	PolicyURI         string   `json:"policy_uri"`
	TermsOfServiceURI string   `json:"tos_uri"`
	ClientURI         string   `json:"client_uri"`
	LogoURI           string   `json:"logo_uri"`
	// GrantTypes        []string `json:"grant_types"`
	// ResponseTypes     []string `json:"response_types"`
}

func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *Client) GetHashedSecret() []byte {
	return c.Secret
}
