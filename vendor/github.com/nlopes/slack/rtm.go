package slack

import (
	"fmt"
	"net/url"
)

// StartRTM calls the "rtm.start" endpoint and returns the provided URL and the full Info
// block.
//
// To have a fully managed Websocket connection, use `NewRTM`, and call `ManageConnection()`
// on it.
func (api *Client) StartRTM() (info *Info, websocketURL string, err error) {
	response := &infoResponseFull{}
	err = post("rtm.start", url.Values{"token": {api.config.token}}, response, api.debug)
	if err != nil {
		return nil, "", fmt.Errorf("post: %s", err)
	}
	if !response.Ok {
		return nil, "", response.Error
	}

	// websocket.Dial does not accept url without the port (yet)
	// Fixed by: https://github.com/golang/net/commit/5058c78c3627b31e484a81463acd51c7cecc06f3
	// but slack returns the address with no port, so we have to fix it
	api.Debugln("Using URL:", response.Info.URL)
	websocketURL, err = websocketizeURLPort(response.Info.URL)
	if err != nil {
		return nil, "", fmt.Errorf("parsing response URL: %s", err)
	}

	return &response.Info, websocketURL, nil
}

// NewRTM returns a RTM, which provides a fully managed connection to
// Slack's websocket-based Real-Time Messaging protocol./
func (api *Client) NewRTM() *RTM {
	return newRTM(api)
}
