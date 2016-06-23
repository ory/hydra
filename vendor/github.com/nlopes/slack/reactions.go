package slack

import (
	"errors"
	"net/url"
	"strconv"
)

// ItemReaction is the reactions that have happened on an item.
type ItemReaction struct {
	Name  string   `json:"name"`
	Count int      `json:"count"`
	Users []string `json:"users"`
}

// ReactedItem is an item that was reacted to, and the details of the
// reactions.
type ReactedItem struct {
	Item
	Reactions []ItemReaction
}

// GetReactionsParameters is the inputs to get reactions to an item.
type GetReactionsParameters struct {
	Full bool
}

// NewGetReactionsParameters initializes the inputs to get reactions to an item.
func NewGetReactionsParameters() GetReactionsParameters {
	return GetReactionsParameters{
		Full: false,
	}
}

type getReactionsResponseFull struct {
	Type string
	M    struct {
		Reactions []ItemReaction
	} `json:"message"`
	F struct {
		Reactions []ItemReaction
	} `json:"file"`
	FC struct {
		Reactions []ItemReaction
	} `json:"comment"`
	SlackResponse
}

func (res getReactionsResponseFull) extractReactions() []ItemReaction {
	switch res.Type {
	case "message":
		return res.M.Reactions
	case "file":
		return res.F.Reactions
	case "file_comment":
		return res.FC.Reactions
	}
	return []ItemReaction{}
}

const (
	DEFAULT_REACTIONS_USER  = ""
	DEFAULT_REACTIONS_COUNT = 100
	DEFAULT_REACTIONS_PAGE  = 1
	DEFAULT_REACTIONS_FULL  = false
)

// ListReactionsParameters is the inputs to find all reactions by a user.
type ListReactionsParameters struct {
	User  string
	Count int
	Page  int
	Full  bool
}

// NewListReactionsParameters initializes the inputs to find all reactions
// performed by a user.
func NewListReactionsParameters() ListReactionsParameters {
	return ListReactionsParameters{
		User:  DEFAULT_REACTIONS_USER,
		Count: DEFAULT_REACTIONS_COUNT,
		Page:  DEFAULT_REACTIONS_PAGE,
		Full:  DEFAULT_REACTIONS_FULL,
	}
}

type listReactionsResponseFull struct {
	Items []struct {
		Type    string
		Channel string
		M       struct {
			*Message
		} `json:"message"`
		F struct {
			*File
			Reactions []ItemReaction
		} `json:"file"`
		FC struct {
			*Comment
			Reactions []ItemReaction
		} `json:"comment"`
	}
	Paging `json:"paging"`
	SlackResponse
}

func (res listReactionsResponseFull) extractReactedItems() []ReactedItem {
	items := make([]ReactedItem, len(res.Items))
	for i, input := range res.Items {
		item := ReactedItem{}
		item.Type = input.Type
		switch input.Type {
		case "message":
			item.Channel = input.Channel
			item.Message = input.M.Message
			item.Reactions = input.M.Reactions
		case "file":
			item.File = input.F.File
			item.Reactions = input.F.Reactions
		case "file_comment":
			item.File = input.F.File
			item.Comment = input.FC.Comment
			item.Reactions = input.FC.Reactions
		}
		items[i] = item
	}
	return items
}

// AddReaction adds a reaction emoji to a message, file or file comment.
func (api *Client) AddReaction(name string, item ItemRef) error {
	values := url.Values{
		"token": {api.config.token},
	}
	if name != "" {
		values.Set("name", name)
	}
	if item.Channel != "" {
		values.Set("channel", string(item.Channel))
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
	if err := post("reactions.add", values, response, api.debug); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// RemoveReaction removes a reaction emoji from a message, file or file comment.
func (api *Client) RemoveReaction(name string, item ItemRef) error {
	values := url.Values{
		"token": {api.config.token},
	}
	if name != "" {
		values.Set("name", name)
	}
	if item.Channel != "" {
		values.Set("channel", string(item.Channel))
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
	if err := post("reactions.remove", values, response, api.debug); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// GetReactions returns details about the reactions on an item.
func (api *Client) GetReactions(item ItemRef, params GetReactionsParameters) ([]ItemReaction, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if item.Channel != "" {
		values.Set("channel", string(item.Channel))
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
	if params.Full != DEFAULT_REACTIONS_FULL {
		values.Set("full", strconv.FormatBool(params.Full))
	}
	response := &getReactionsResponseFull{}
	if err := post("reactions.get", values, response, api.debug); err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response.extractReactions(), nil
}

// ListReactions returns information about the items a user reacted to.
func (api *Client) ListReactions(params ListReactionsParameters) ([]ReactedItem, *Paging, error) {
	values := url.Values{
		"token": {api.config.token},
	}
	if params.User != DEFAULT_REACTIONS_USER {
		values.Add("user", params.User)
	}
	if params.Count != DEFAULT_REACTIONS_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Page != DEFAULT_REACTIONS_PAGE {
		values.Add("page", strconv.Itoa(params.Page))
	}
	if params.Full != DEFAULT_REACTIONS_FULL {
		values.Add("full", strconv.FormatBool(params.Full))
	}
	response := &listReactionsResponseFull{}
	err := post("reactions.list", values, response, api.debug)
	if err != nil {
		return nil, nil, err
	}
	if !response.Ok {
		return nil, nil, errors.New(response.Error)
	}
	return response.extractReactedItems(), &response.Paging, nil
}
