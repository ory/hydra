package slack

import (
	"errors"
	"net/url"
	"strconv"
)

type channelResponseFull struct {
	Channel      Channel   `json:"channel"`
	Channels     []Channel `json:"channels"`
	Purpose      string    `json:"purpose"`
	Topic        string    `json:"topic"`
	NotInChannel bool      `json:"not_in_channel"`
	History
	SlackResponse
}

// Channel contains information about the channel
type Channel struct {
	groupConversation
	IsChannel bool `json:"is_channel"`
	IsGeneral bool `json:"is_general"`
	IsMember  bool `json:"is_member"`
}

func channelRequest(path string, values url.Values, debug bool) (*channelResponseFull, error) {
	response := &channelResponseFull{}
	err := post(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// ArchiveChannel archives the given channel
func (api *Client) ArchiveChannel(channel string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
	}
	_, err := channelRequest("channels.archive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// UnarchiveChannel unarchives the given channel
func (api *Client) UnarchiveChannel(channel string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
	}
	_, err := channelRequest("channels.unarchive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// CreateChannel creates a channel with the given name and returns a *Channel
func (api *Client) CreateChannel(channel string) (*Channel, error) {
	values := url.Values{
		"token": {api.config.token},
		"name":  {channel},
	}
	response, err := channelRequest("channels.create", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// GetChannelHistory retrieves the channel history
func (api *Client) GetChannelHistory(channel string, params HistoryParameters) (*History, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
	}
	if params.Latest != DEFAULT_HISTORY_LATEST {
		values.Add("latest", params.Latest)
	}
	if params.Oldest != DEFAULT_HISTORY_OLDEST {
		values.Add("oldest", params.Oldest)
	}
	if params.Count != DEFAULT_HISTORY_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Inclusive != DEFAULT_HISTORY_INCLUSIVE {
		if params.Inclusive {
			values.Add("inclusive", "1")
		} else {
			values.Add("inclusive", "0")
		}
	}
	if params.Unreads != DEFAULT_HISTORY_UNREADS {
		if params.Unreads {
			values.Add("unreads", "1")
		} else {
			values.Add("unreads", "0")
		}
	}
	response, err := channelRequest("channels.history", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.History, nil
}

// GetChannelInfo retrieves the given channel
func (api *Client) GetChannelInfo(channel string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
	}
	response, err := channelRequest("channels.info", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// InviteUserToChannel invites a user to a given channel and returns a *Channel
func (api *Client) InviteUserToChannel(channel, user string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
		"user":    {user},
	}
	response, err := channelRequest("channels.invite", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// JoinChannel joins the currently authenticated user to a channel
func (api *Client) JoinChannel(channel string) (*Channel, error) {
	values := url.Values{
		"token": {api.config.token},
		"name":  {channel},
	}
	response, err := channelRequest("channels.join", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil
}

// LeaveChannel makes the authenticated user leave the given channel
func (api *Client) LeaveChannel(channel string) (bool, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
	}
	response, err := channelRequest("channels.leave", values, api.debug)
	if err != nil {
		return false, err
	}
	if response.NotInChannel {
		return response.NotInChannel, nil
	}
	return false, nil
}

// KickUserFromChannel kicks a user from a given channel
func (api *Client) KickUserFromChannel(channel, user string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
		"user":    {user},
	}
	_, err := channelRequest("channels.kick", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// GetChannels retrieves all the channels
func (api *Client) GetChannels(excludeArchived bool) ([]Channel, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if excludeArchived {
		values.Add("exclude_archived", "1")
	}
	response, err := channelRequest("channels.list", values, api.debug)
	if err != nil {
		return nil, err
	}
	return response.Channels, nil
}

// SetChannelReadMark sets the read mark of a given channel to a specific point
// Clients should try to avoid making this call too often. When needing to mark a read position, a client should set a
// timer before making the call. In this way, any further updates needed during the timeout will not generate extra calls
// (just one per channel). This is useful for when reading scroll-back history, or following a busy live channel. A
// timeout of 5 seconds is a good starting point. Be sure to flush these calls on shutdown/logout.
func (api *Client) SetChannelReadMark(channel, ts string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
		"ts":      {ts},
	}
	_, err := channelRequest("channels.mark", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// RenameChannel renames a given channel
func (api *Client) RenameChannel(channel, name string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
		"name":    {name},
	}
	// XXX: the created entry in this call returns a string instead of a number
	// so I may have to do some workaround to solve it.
	response, err := channelRequest("channels.rename", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil

}

// SetChannelPurpose sets the channel purpose and returns the purpose that was
// successfully set
func (api *Client) SetChannelPurpose(channel, purpose string) (string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
		"purpose": {purpose},
	}
	response, err := channelRequest("channels.setPurpose", values, api.debug)
	if err != nil {
		return "", err
	}
	return response.Purpose, nil
}

// SetChannelTopic sets the channel topic and returns the topic that was successfully set
func (api *Client) SetChannelTopic(channel, topic string) (string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
		"topic":   {topic},
	}
	response, err := channelRequest("channels.setTopic", values, api.debug)
	if err != nil {
		return "", err
	}
	return response.Topic, nil
}
