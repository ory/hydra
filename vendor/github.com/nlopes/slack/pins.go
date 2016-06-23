package slack

import (
	"errors"
	"net/url"
)

type listPinsResponseFull struct {
	Items  []Item
	Paging `json:"paging"`
	SlackResponse
}

// AddPin pins an item in a channel
func (api *Client) AddPin(channel string, item ItemRef) error {
	values := url.Values{
		"channel": {channel},
		"token":   {api.config.token},
	}
	if item.Timestamp != "" {
		values.Set("timestamp", string(item.Timestamp))
	}
	if item.File != "" {
		values.Set("file", string(item.File))
	}
	if item.Comment != "" {
		values.Set("file_comment", string(item.Comment))
	}
	response := &SlackResponse{}
	if err := post("pins.add", values, response, api.debug); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// RemovePin un-pins an item from a channel
func (api *Client) RemovePin(channel string, item ItemRef) error {
	values := url.Values{
		"channel": {channel},
		"token":   {api.config.token},
	}
	if item.Timestamp != "" {
		values.Set("timestamp", string(item.Timestamp))
	}
	if item.File != "" {
		values.Set("file", string(item.File))
	}
	if item.Comment != "" {
		values.Set("file_comment", string(item.Comment))
	}
	response := &SlackResponse{}
	if err := post("pins.remove", values, response, api.debug); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// ListPins returns information about the items a user reacted to.
func (api *Client) ListPins(channel string) ([]Item, *Paging, error) {
	values := url.Values{
		"channel": {channel},
		"token":   {api.config.token},
	}
	response := &listPinsResponseFull{}
	err := post("pins.list", values, response, api.debug)
	if err != nil {
		return nil, nil, err
	}
	if !response.Ok {
		return nil, nil, errors.New(response.Error)
	}
	return response.Items, &response.Paging, nil
}
