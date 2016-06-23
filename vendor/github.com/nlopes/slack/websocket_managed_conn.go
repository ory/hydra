package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"time"

	"golang.org/x/net/websocket"
)

// ManageConnection can be called on a Slack RTM instance returned by the
// NewRTM method. It will connect to the slack RTM API and handle all incoming
// and outgoing events. If a connection fails then it will attempt to reconnect
// and will notify any listeners through an error event on the IncomingEvents
// channel.
//
// If the connection ends and the disconnect was unintentional then this will
// attempt to reconnect.
//
// This should only be called once per slack API! Otherwise expect undefined
// behavior.
//
// The defined error events are located in websocket_internals.go.
func (rtm *RTM) ManageConnection() {
	var connectionCount int
	for {
		connectionCount++
		// start trying to connect
		// the returned err is already passed onto the IncomingEvents channel
		info, conn, err := rtm.connect(connectionCount)
		// if err != nil then the connection is sucessful - otherwise it is
		// fatal
		if err != nil {
			return
		}
		rtm.info = info
		rtm.IncomingEvents <- RTMEvent{"connected", &ConnectedEvent{
			ConnectionCount: connectionCount,
			Info:            info,
		}}

		rtm.conn = conn
		rtm.isConnected = true

		keepRunning := make(chan bool)
		// we're now connected (or have failed fatally) so we can set up
		// listeners
		go rtm.handleIncomingEvents(keepRunning)

		// this should be a blocking call until the connection has ended
		rtm.handleEvents(keepRunning, 30*time.Second)

		// after being disconnected we need to check if it was intentional
		// if not then we should try to reconnect
		if rtm.wasIntentional {
			return
		}
		// else continue and run the loop again to connect
	}
}

// connect attempts to connect to the slack websocket API. It handles any
// errors that occur while connecting and will return once a connection
// has been successfully opened.
func (rtm *RTM) connect(connectionCount int) (*Info, *websocket.Conn, error) {
	// used to provide exponential backoff wait time with jitter before trying
	// to connect to slack again
	boff := &backoff{
		Min:    100 * time.Millisecond,
		Max:    5 * time.Minute,
		Factor: 2,
		Jitter: true,
	}

	for {
		// send connecting event
		rtm.IncomingEvents <- RTMEvent{"connecting", &ConnectingEvent{
			Attempt:         boff.attempts + 1,
			ConnectionCount: connectionCount,
		}}
		// attempt to start the connection
		info, conn, err := rtm.startRTMAndDial()
		if err == nil {
			return info, conn, nil
		}
		// check for fatal errors - currently only invalid_auth
		if sErr, ok := err.(*WebError); ok && (sErr.Error() == "invalid_auth" || sErr.Error() == "account_inactive") {
			rtm.IncomingEvents <- RTMEvent{"invalid_auth", &InvalidAuthEvent{}}
			return nil, nil, sErr
		}
		// any other errors are treated as recoverable and we try again after
		// sending the event along the IncomingEvents channel
		rtm.IncomingEvents <- RTMEvent{"connection_error", &ConnectionErrorEvent{
			Attempt:  boff.attempts,
			ErrorObj: err,
		}}
		// get time we should wait before attempting to connect again
		dur := boff.Duration()
		rtm.Debugf("reconnection %d failed: %s", boff.attempts+1, err)
		rtm.Debugln(" -> reconnecting in", dur)
		time.Sleep(dur)
	}
}

// startRTMAndDial attemps to connect to the slack websocket. It returns the
// full information returned by the "rtm.start" method on the slack API.
func (rtm *RTM) startRTMAndDial() (*Info, *websocket.Conn, error) {
	info, url, err := rtm.StartRTM()
	if err != nil {
		return nil, nil, err
	}

	conn, err := websocketProxyDial(url, "http://api.slack.com")
	if err != nil {
		return nil, nil, err
	}
	return info, conn, err
}

// killConnection stops the websocket connection and signals to all goroutines
// that they should cease listening to the connection for events.
//
// This should not be called directly! Instead a boolean value (true for
// intentional, false otherwise) should be sent to the killChannel on the RTM.
func (rtm *RTM) killConnection(keepRunning chan bool, intentional bool) error {
	rtm.Debugln("killing connection")
	if rtm.isConnected {
		close(keepRunning)
	}
	rtm.isConnected = false
	rtm.wasIntentional = intentional
	err := rtm.conn.Close()
	rtm.IncomingEvents <- RTMEvent{"disconnected", &DisconnectedEvent{intentional}}
	return err
}

// handleEvents is a blocking function that handles all events. This sends
// pings when asked to (on rtm.forcePing) and upon every given elapsed
// interval. This also sends outgoing messages that are received from the RTM's
// outgoingMessages channel. This also handles incoming raw events from the RTM
// rawEvents channel.
func (rtm *RTM) handleEvents(keepRunning chan bool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		// catch "stop" signal on channel close
		case intentional := <-rtm.killChannel:
			_ = rtm.killConnection(keepRunning, intentional)
			return
			// send pings on ticker interval
		case <-ticker.C:
			err := rtm.ping()
			if err != nil {
				_ = rtm.killConnection(keepRunning, false)
				return
			}
		case <-rtm.forcePing:
			err := rtm.ping()
			if err != nil {
				_ = rtm.killConnection(keepRunning, false)
				return
			}
		// listen for messages that need to be sent
		case msg := <-rtm.outgoingMessages:
			rtm.sendOutgoingMessage(msg)
		// listen for incoming messages that need to be parsed
		case rawEvent := <-rtm.rawEvents:
			rtm.handleRawEvent(rawEvent)
		}
	}
}

// handleIncomingEvents monitors the RTM's opened websocket for any incoming
// events. It pushes the raw events onto the RTM channel rawEvents.
//
// This will stop executing once the RTM's keepRunning channel has been closed
// or has anything sent to it.
func (rtm *RTM) handleIncomingEvents(keepRunning <-chan bool) {
	for {
		// non-blocking listen to see if channel is closed
		select {
		// catch "stop" signal on channel close
		case <-keepRunning:
			return
		default:
			rtm.receiveIncomingEvent()
		}
	}
}

// sendOutgoingMessage sends the given OutgoingMessage to the slack websocket.
//
// It does not currently detect if a outgoing message fails due to a disconnect
// and instead lets a future failed 'PING' detect the failed connection.
func (rtm *RTM) sendOutgoingMessage(msg OutgoingMessage) {
	rtm.Debugln("Sending message:", msg)
	if len(msg.Text) > MaxMessageTextLength {
		rtm.IncomingEvents <- RTMEvent{"outgoing_error", &MessageTooLongEvent{
			Message:   msg,
			MaxLength: MaxMessageTextLength,
		}}
		return
	}
	err := websocket.JSON.Send(rtm.conn, msg)
	if err != nil {
		rtm.IncomingEvents <- RTMEvent{"outgoing_error", &OutgoingErrorEvent{
			Message:  msg,
			ErrorObj: err,
		}}
		// TODO force ping?
	}
}

// ping sends a 'PING' message to the RTM's websocket. If the 'PING' message
// fails to send then this returns an error signifying that the connection
// should be considered disconnected.
//
// This does not handle incoming 'PONG' responses but does store the time of
// each successful 'PING' send so latency can be detected upon a 'PONG'
// response.
func (rtm *RTM) ping() error {
	id := rtm.idGen.Next()
	rtm.Debugln("Sending PING ", id)
	rtm.pings[id] = time.Now()

	msg := &Ping{ID: id, Type: "ping"}
	err := websocket.JSON.Send(rtm.conn, msg)
	if err != nil {
		rtm.Debugf("RTM Error sending 'PING %d': %s", id, err.Error())
		return err
	}
	return nil
}

// receiveIncomingEvent attempts to receive an event from the RTM's websocket.
// This will block until a frame is available from the websocket.
func (rtm *RTM) receiveIncomingEvent() {
	event := json.RawMessage{}
	err := websocket.JSON.Receive(rtm.conn, &event)
	if err == io.EOF {
		// EOF's don't seem to signify a failed connection so instead we ignore
		// them here and detect a failed connection upon attempting to send a
		// 'PING' message

		// trigger a 'PING' to detect pontential websocket disconnect
		rtm.forcePing <- true
		return
	} else if err != nil {
		rtm.IncomingEvents <- RTMEvent{"incoming_error", &IncomingEventError{
			ErrorObj: err,
		}}
		// force a ping here too?
		return
	} else if len(event) == 0 {
		rtm.Debugln("Received empty event")
		return
	}
	rtm.Debugln("Incoming Event:", string(event[:]))
	rtm.rawEvents <- event
}

// handleRawEvent takes a raw JSON message received from the slack websocket
// and handles the encoded event.
func (rtm *RTM) handleRawEvent(rawEvent json.RawMessage) {
	event := &Event{}
	err := json.Unmarshal(rawEvent, event)
	if err != nil {
		rtm.IncomingEvents <- RTMEvent{"unmarshalling_error", &UnmarshallingErrorEvent{err}}
		return
	}
	switch event.Type {
	case "":
		rtm.handleAck(rawEvent)
	case "hello":
		rtm.IncomingEvents <- RTMEvent{"hello", &HelloEvent{}}
	case "pong":
		rtm.handlePong(rawEvent)
	default:
		rtm.handleEvent(event.Type, rawEvent)
	}
}

// handleAck handles an incoming 'ACK' message.
func (rtm *RTM) handleAck(event json.RawMessage) {
	ack := &AckMessage{}
	if err := json.Unmarshal(event, ack); err != nil {
		rtm.Debugln("RTM Error unmarshalling 'ack' event:", err)
		rtm.Debugln(" -> Erroneous 'ack' event:", string(event))
		return
	}
	if ack.Ok {
		rtm.IncomingEvents <- RTMEvent{"ack", ack}
	} else {
		rtm.IncomingEvents <- RTMEvent{"ack_error", &AckErrorEvent{ack.Error}}
	}
}

// handlePong handles an incoming 'PONG' message which should be in response to
// a previously sent 'PING' message. This is then used to compute the
// connection's latency.
func (rtm *RTM) handlePong(event json.RawMessage) {
	pong := &Pong{}
	if err := json.Unmarshal(event, pong); err != nil {
		rtm.Debugln("RTM Error unmarshalling 'pong' event:", err)
		rtm.Debugln(" -> Erroneous 'ping' event:", string(event))
		return
	}
	if pingTime, exists := rtm.pings[pong.ReplyTo]; exists {
		latency := time.Since(pingTime)
		rtm.IncomingEvents <- RTMEvent{"latency_report", &LatencyReport{Value: latency}}
		delete(rtm.pings, pong.ReplyTo)
	} else {
		rtm.Debugln("RTM Error - unmatched 'pong' event:", string(event))
	}
}

// handleEvent is the "default" response to an event that does not have a
// special case. It matches the command's name to a mapping of defined events
// and then sends the corresponding event struct to the IncomingEvents channel.
// If the event type is not found or the event cannot be unmarshalled into the
// correct struct then this sends an UnmarshallingErrorEvent to the
// IncomingEvents channel.
func (rtm *RTM) handleEvent(typeStr string, event json.RawMessage) {
	v, exists := eventMapping[typeStr]
	if !exists {
		rtm.Debugf("RTM Error, received unmapped event %q: %s\n", typeStr, string(event))
		err := fmt.Errorf("RTM Error: Received unmapped event %q: %s\n", typeStr, string(event))
		rtm.IncomingEvents <- RTMEvent{"unmarshalling_error", &UnmarshallingErrorEvent{err}}
		return
	}
	t := reflect.TypeOf(v)
	recvEvent := reflect.New(t).Interface()
	err := json.Unmarshal(event, recvEvent)
	if err != nil {
		rtm.Debugf("RTM Error, could not unmarshall event %q: %s\n", typeStr, string(event))
		err := fmt.Errorf("RTM Error: Could not unmarshall event %q: %s\n", typeStr, string(event))
		rtm.IncomingEvents <- RTMEvent{"unmarshalling_error", &UnmarshallingErrorEvent{err}}
		return
	}
	rtm.IncomingEvents <- RTMEvent{typeStr, recvEvent}
}

// eventMapping holds a mapping of event names to their corresponding struct
// implementations. The structs should be instances of the unmarshalling
// target for the matching event type.
var eventMapping = map[string]interface{}{
	"message":         MessageEvent{},
	"presence_change": PresenceChangeEvent{},
	"user_typing":     UserTypingEvent{},

	"channel_marked":          ChannelMarkedEvent{},
	"channel_created":         ChannelCreatedEvent{},
	"channel_joined":          ChannelJoinedEvent{},
	"channel_left":            ChannelLeftEvent{},
	"channel_deleted":         ChannelDeletedEvent{},
	"channel_rename":          ChannelRenameEvent{},
	"channel_archive":         ChannelArchiveEvent{},
	"channel_unarchive":       ChannelUnarchiveEvent{},
	"channel_history_changed": ChannelHistoryChangedEvent{},

	"dnd_updated":      DNDUpdatedEvent{},
	"dnd_updated_user": DNDUpdatedEvent{},

	"im_created":         IMCreatedEvent{},
	"im_open":            IMOpenEvent{},
	"im_close":           IMCloseEvent{},
	"im_marked":          IMMarkedEvent{},
	"im_history_changed": IMHistoryChangedEvent{},

	"group_marked":          GroupMarkedEvent{},
	"group_open":            GroupOpenEvent{},
	"group_joined":          GroupJoinedEvent{},
	"group_left":            GroupLeftEvent{},
	"group_close":           GroupCloseEvent{},
	"group_rename":          GroupRenameEvent{},
	"group_archive":         GroupArchiveEvent{},
	"group_unarchive":       GroupUnarchiveEvent{},
	"group_history_changed": GroupHistoryChangedEvent{},

	"file_created":         FileCreatedEvent{},
	"file_shared":          FileSharedEvent{},
	"file_unshared":        FileUnsharedEvent{},
	"file_public":          FilePublicEvent{},
	"file_private":         FilePrivateEvent{},
	"file_change":          FileChangeEvent{},
	"file_deleted":         FileDeletedEvent{},
	"file_comment_added":   FileCommentAddedEvent{},
	"file_comment_edited":  FileCommentEditedEvent{},
	"file_comment_deleted": FileCommentDeletedEvent{},

	"pin_added":   PinAddedEvent{},
	"pin_removed": PinRemovedEvent{},

	"star_added":   StarAddedEvent{},
	"star_removed": StarRemovedEvent{},

	"reaction_added":   ReactionAddedEvent{},
	"reaction_removed": ReactionRemovedEvent{},

	"pref_change": PrefChangeEvent{},

	"team_join":              TeamJoinEvent{},
	"team_rename":            TeamRenameEvent{},
	"team_pref_change":       TeamPrefChangeEvent{},
	"team_domain_change":     TeamDomainChangeEvent{},
	"team_migration_started": TeamMigrationStartedEvent{},

	"manual_presence_change": ManualPresenceChangeEvent{},

	"user_change": UserChangeEvent{},

	"emoji_changed": EmojiChangedEvent{},

	"commands_changed": CommandsChangedEvent{},

	"email_domain_changed": EmailDomainChangedEvent{},

	"bot_added":   BotAddedEvent{},
	"bot_changed": BotChangedEvent{},

	"accounts_changed": AccountsChangedEvent{},

	"reconnect_url": ReconnectUrlEvent{},
}
