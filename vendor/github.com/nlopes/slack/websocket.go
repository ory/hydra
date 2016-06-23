package slack

import (
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/net/websocket"
)

const (
	// MaxMessageTextLength is the current maximum message length in number of characters as defined here
	// https://api.slack.com/rtm#limits
	MaxMessageTextLength = 4000
)

// RTM represents a managed websocket connection. It also supports
// all the methods of the `Client` type.
//
// Create this element with Client's NewRTM().
type RTM struct {
	idGen IDGenerator
	pings map[int]time.Time

	// Connection life-cycle
	conn             *websocket.Conn
	IncomingEvents   chan RTMEvent
	outgoingMessages chan OutgoingMessage
	killChannel      chan bool
	forcePing        chan bool
	rawEvents        chan json.RawMessage
	wasIntentional   bool
	isConnected      bool

	// Client is the main API, embedded
	Client
	websocketURL string

	// UserDetails upon connection
	info *Info
}

// NewRTM returns a RTM, which provides a fully managed connection to
// Slack's websocket-based Real-Time Messaging protocol.
func newRTM(api *Client) *RTM {
	return &RTM{
		Client:           *api,
		IncomingEvents:   make(chan RTMEvent, 50),
		outgoingMessages: make(chan OutgoingMessage, 20),
		pings:            make(map[int]time.Time),
		isConnected:      false,
		wasIntentional:   true,
		killChannel:      make(chan bool),
		forcePing:        make(chan bool),
		rawEvents:        make(chan json.RawMessage),
		idGen:            NewSafeID(1),
	}
}

// Disconnect and wait, blocking until a successful disconnection.
func (rtm *RTM) Disconnect() error {
	if !rtm.isConnected {
		return errors.New("Invalid call to Disconnect - Slack API is already disconnected")
	}
	rtm.killChannel <- true
	return nil
}

// Reconnect only makes sense if you've successfully disconnectd with Disconnect().
func (rtm *RTM) Reconnect() error {
	logger.Println("RTM::Reconnect not implemented!")
	return nil
}

// GetInfo returns the info structure received when calling
// "startrtm", holding all channels, groups and other metadata needed
// to implement a full chat client. It will be non-nil after a call to
// StartRTM().
func (rtm *RTM) GetInfo() *Info {
	return rtm.info
}

// SendMessage submits a simple message through the websocket.  For
// more complicated messages, use `rtm.PostMessage` with a complete
// struct describing your attachments and all.
func (rtm *RTM) SendMessage(msg *OutgoingMessage) {
	if msg == nil {
		rtm.Debugln("Error: Attempted to SendMessage(nil)")
		return
	}

	rtm.outgoingMessages <- *msg
}
