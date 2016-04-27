package client

type Client struct {
	ID           string   `json:"id"`
	Secret       []byte   `json:"secret"`
	RedirectURIs []string `json:"redirectURIs"`
	Owner        string   `json:"owner"`
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
