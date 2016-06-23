package slack

import (
	"errors"
	"net/url"
)

type emojiResponseFull struct {
	Emoji map[string]string `json:"emoji"`
	SlackResponse
}

// GetEmoji retrieves all the emojis
func (api *Client) GetEmoji() (map[string]string, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	response := &emojiResponseFull{}
	err := post("emoji.list", values, response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.Emoji, nil
}
