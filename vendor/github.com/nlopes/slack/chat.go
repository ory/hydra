package slack

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

const (
	DEFAULT_MESSAGE_USERNAME     = ""
	DEFAULT_MESSAGE_ASUSER       = false
	DEFAULT_MESSAGE_PARSE        = ""
	DEFAULT_MESSAGE_LINK_NAMES   = 0
	DEFAULT_MESSAGE_UNFURL_LINKS = false
	DEFAULT_MESSAGE_UNFURL_MEDIA = true
	DEFAULT_MESSAGE_ICON_URL     = ""
	DEFAULT_MESSAGE_ICON_EMOJI   = ""
	DEFAULT_MESSAGE_MARKDOWN     = true
	DEFAULT_MESSAGE_ESCAPE_TEXT  = true
)

type chatResponseFull struct {
	Channel   string `json:"channel"`
	Timestamp string `json:"ts"`
	Text      string `json:"text"`
	SlackResponse
}

// PostMessageParameters contains all the parameters necessary (including the optional ones) for a PostMessage() request
type PostMessageParameters struct {
	Text        string
	Username    string
	AsUser      bool
	Parse       string
	LinkNames   int
	Attachments []Attachment
	UnfurlLinks bool
	UnfurlMedia bool
	IconURL     string
	IconEmoji   string
	Markdown    bool `json:"mrkdwn,omitempty"`
	EscapeText  bool
}

// NewPostMessageParameters provides an instance of PostMessageParameters with all the sane default values set
func NewPostMessageParameters() PostMessageParameters {
	return PostMessageParameters{
		Username:    DEFAULT_MESSAGE_USERNAME,
		AsUser:      DEFAULT_MESSAGE_ASUSER,
		Parse:       DEFAULT_MESSAGE_PARSE,
		LinkNames:   DEFAULT_MESSAGE_LINK_NAMES,
		Attachments: nil,
		UnfurlLinks: DEFAULT_MESSAGE_UNFURL_LINKS,
		UnfurlMedia: DEFAULT_MESSAGE_UNFURL_MEDIA,
		IconURL:     DEFAULT_MESSAGE_ICON_URL,
		IconEmoji:   DEFAULT_MESSAGE_ICON_EMOJI,
		Markdown:    DEFAULT_MESSAGE_MARKDOWN,
		EscapeText:  DEFAULT_MESSAGE_ESCAPE_TEXT,
	}
}

func chatRequest(path string, values url.Values, debug bool) (*chatResponseFull, error) {
	response := &chatResponseFull{}
	err := post(path, values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

// DeleteMessage deletes a message in a channel
func (api *Client) DeleteMessage(channel, messageTimestamp string) (string, string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
		"ts":      {messageTimestamp},
	}
	response, err := chatRequest("chat.delete", values, api.debug)
	if err != nil {
		return "", "", err
	}
	return response.Channel, response.Timestamp, nil
}

func escapeMessage(message string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return replacer.Replace(message)
}

// PostMessage sends a message to a channel.
// Message is escaped by default according to https://api.slack.com/docs/formatting
// Use http://davestevens.github.io/slack-message-builder/ to help crafting your message.
func (api *Client) PostMessage(channel, text string, params PostMessageParameters) (string, string, error) {
	if params.EscapeText {
		text = escapeMessage(text)
	}
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
		"text":    {text},
	}
	if params.Username != DEFAULT_MESSAGE_USERNAME {
		values.Set("username", string(params.Username))
	}
	if params.AsUser != DEFAULT_MESSAGE_ASUSER {
		values.Set("as_user", "true")
	}
	if params.Parse != DEFAULT_MESSAGE_PARSE {
		values.Set("parse", string(params.Parse))
	}
	if params.LinkNames != DEFAULT_MESSAGE_LINK_NAMES {
		values.Set("link_names", "1")
	}
	if params.Attachments != nil {
		attachments, err := json.Marshal(params.Attachments)
		if err != nil {
			return "", "", err
		}
		values.Set("attachments", string(attachments))
	}
	if params.UnfurlLinks != DEFAULT_MESSAGE_UNFURL_LINKS {
		values.Set("unfurl_links", "true")
	}
	// I want to send a message with explicit `as_user` `true` and `unfurl_links` `false` in request.
	// Because setting `as_user` to `true` will change the default value for `unfurl_links` to `true` on Slack API side.
	if params.AsUser != DEFAULT_MESSAGE_ASUSER && params.UnfurlLinks == DEFAULT_MESSAGE_UNFURL_LINKS {
		values.Set("unfurl_links", "false")
	}
	if params.UnfurlMedia != DEFAULT_MESSAGE_UNFURL_MEDIA {
		values.Set("unfurl_media", "false")
	}
	if params.IconURL != DEFAULT_MESSAGE_ICON_URL {
		values.Set("icon_url", params.IconURL)
	}
	if params.IconEmoji != DEFAULT_MESSAGE_ICON_EMOJI {
		values.Set("icon_emoji", params.IconEmoji)
	}
	if params.Markdown != DEFAULT_MESSAGE_MARKDOWN {
		values.Set("mrkdwn", "false")
	}

	response, err := chatRequest("chat.postMessage", values, api.debug)
	if err != nil {
		return "", "", err
	}
	return response.Channel, response.Timestamp, nil
}

// UpdateMessage updates a message in a channel
func (api *Client) UpdateMessage(channel, timestamp, text string) (string, string, string, error) {
	values := url.Values{
		"token":   {api.config.token},
		"channel": {channel},
		"text":    {escapeMessage(text)},
		"ts":      {timestamp},
	}
	response, err := chatRequest("chat.update", values, api.debug)
	if err != nil {
		return "", "", "", err
	}
	return response.Channel, response.Timestamp, response.Text, nil
}
