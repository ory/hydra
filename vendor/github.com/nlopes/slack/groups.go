package slack

import (
	"errors"
	"net/url"
	"strconv"
)

// Group contains all the information for a group
type Group struct {
	groupConversation
	IsGroup bool `json:"is_group"`
}

type groupResponseFull struct {
	Group          Group   `json:"group"`
	Groups         []Group `json:"groups"`
	Purpose        string  `json:"purpose"`
	Topic          string  `json:"topic"`
	NotInGroup     bool    `json:"not_in_group"`
	NoOp           bool    `json:"no_op"`
	AlreadyClosed  bool    `json:"already_closed"`
	AlreadyOpen    bool    `json:"already_open"`
	AlreadyInGroup bool    `json:"already_in_group"`
	Channel        Channel `json:"channel"`
	History
	SlackResponse
}

func groupRequest(path string, values url.Values, debug bool) (*groupResponseFull, error) {
	response := &groupResponseFull{}
	err := post(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// ArchiveGroup archives a private group
func (api *Client) ArchiveGroup(group string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
	}
	_, err := groupRequest("groups.archive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// UnarchiveGroup unarchives a private group
func (api *Client) UnarchiveGroup(group string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
	}
	_, err := groupRequest("groups.unarchive", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// CreateGroup creates a private group
func (api *Client) CreateGroup(group string) (*Group, error) {
	values := url.Values{
		"token": {api.config.token},
		"name":  {group},
	}
	response, err := groupRequest("groups.create", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Group, nil
}

// CreateChildGroup creates a new private group archiving the old one
// This method takes an existing private group and performs the following steps:
//   1. Renames the existing group (from "example" to "example-archived").
//   2. Archives the existing group.
//   3. Creates a new group with the name of the existing group.
//   4. Adds all members of the existing group to the new group.
func (api *Client) CreateChildGroup(group string) (*Group, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
	}
	response, err := groupRequest("groups.createChild", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Group, nil
}

// CloseGroup closes a private group
func (api *Client) CloseGroup(group string) (bool, bool, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
	}
	response, err := imRequest("groups.close", values, api.debug)
	if err != nil {
		return false, false, err
	}
	return response.NoOp, response.AlreadyClosed, nil
}

// GetGroupHistory fetches all the history for a private group
func (api *Client) GetGroupHistory(group string, params HistoryParameters) (*History, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
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
	response, err := groupRequest("groups.history", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.History, nil
}

// InviteUserToGroup invites a specific user to a private group
func (api *Client) InviteUserToGroup(group, user string) (*Group, bool, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
		"user":    {user},
	}
	response, err := groupRequest("groups.invite", values, api.debug)
	if err != nil {
		return nil, false, err
	}
	return &response.Group, response.AlreadyInGroup, nil
}

// LeaveGroup makes authenticated user leave the group
func (api *Client) LeaveGroup(group string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
	}
	_, err := groupRequest("groups.leave", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// KickUserFromGroup kicks a user from a group
func (api *Client) KickUserFromGroup(group, user string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
		"user":    {user},
	}
	_, err := groupRequest("groups.kick", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// GetGroups retrieves all groups
func (api *Client) GetGroups(excludeArchived bool) ([]Group, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if excludeArchived {
		values.Add("exclude_archived", "1")
	}
	response, err := groupRequest("groups.list", values, api.debug)
	if err != nil {
		return nil, err
	}
	return response.Groups, nil
}

// GetGroupInfo retrieves the given group
func (api *Client) GetGroupInfo(group string) (*Group, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
	}
	response, err := groupRequest("groups.info", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Group, nil
}

// SetGroupReadMark sets the read mark on a private group
// Clients should try to avoid making this call too often. When needing to mark a read position, a client should set a
// timer before making the call. In this way, any further updates needed during the timeout will not generate extra
// calls (just one per channel). This is useful for when reading scroll-back history, or following a busy live
// channel. A timeout of 5 seconds is a good starting point. Be sure to flush these calls on shutdown/logout.
func (api *Client) SetGroupReadMark(group, ts string) error {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
		"ts":      {ts},
	}
	_, err := groupRequest("groups.mark", values, api.debug)
	if err != nil {
		return err
	}
	return nil
}

// OpenGroup opens a private group
func (api *Client) OpenGroup(group string) (bool, bool, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
	}
	response, err := groupRequest("groups.open", values, api.debug)
	if err != nil {
		return false, false, err
	}
	return response.NoOp, response.AlreadyOpen, nil
}

// RenameGroup renames a group
// XXX: They return a channel, not a group. What is this crap? :(
// Inconsistent api it seems.
func (api *Client) RenameGroup(group, name string) (*Channel, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
		"name":    {name},
	}
	// XXX: the created entry in this call returns a string instead of a number
	// so I may have to do some workaround to solve it.
	response, err := groupRequest("groups.rename", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Channel, nil

}

// SetGroupPurpose sets the group purpose
func (api *Client) SetGroupPurpose(group, purpose string) (string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
		"purpose": {purpose},
	}
	response, err := groupRequest("groups.setPurpose", values, api.debug)
	if err != nil {
		return "", err
	}
	return response.Purpose, nil
}

// SetGroupTopic sets the group topic
func (api *Client) SetGroupTopic(group, topic string) (string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {group},
		"topic":   {topic},
	}
	response, err := groupRequest("groups.setTopic", values, api.debug)
	if err != nil {
		return "", err
	}
	return response.Topic, nil
}
