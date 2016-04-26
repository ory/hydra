package oauth2

type OAuth2Client struct {
	ID           string   `json:"id"`
	Secret       []byte   `json:"secret"`
	RedirectURIs []string `json:"redirectURIs"`
	Owner        string   `json:"owner"`
}

func (c *OAuth2Client) GetID() string {
	return c.ID
}

func (c *OAuth2Client) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *OAuth2Client) GetHashedSecret() []byte {
	return c.Secret
}
