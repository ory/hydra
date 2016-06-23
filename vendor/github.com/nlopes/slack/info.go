package slack

import (
	"fmt"
	"time"
)

// UserPrefs needs to be implemented
type UserPrefs struct {
	// "highlight_words":"",
	// "user_colors":"",
	// "color_names_in_list":true,
	// "growls_enabled":true,
	// "tz":"Europe\/London",
	// "push_dm_alert":true,
	// "push_mention_alert":true,
	// "push_everything":true,
	// "push_idle_wait":2,
	// "push_sound":"b2.mp3",
	// "push_loud_channels":"",
	// "push_mention_channels":"",
	// "push_loud_channels_set":"",
	// "email_alerts":"instant",
	// "email_alerts_sleep_until":0,
	// "email_misc":false,
	// "email_weekly":true,
	// "welcome_message_hidden":false,
	// "all_channels_loud":true,
	// "loud_channels":"",
	// "never_channels":"",
	// "loud_channels_set":"",
	// "show_member_presence":true,
	// "search_sort":"timestamp",
	// "expand_inline_imgs":true,
	// "expand_internal_inline_imgs":true,
	// "expand_snippets":false,
	// "posts_formatting_guide":true,
	// "seen_welcome_2":true,
	// "seen_ssb_prompt":false,
	// "search_only_my_channels":false,
	// "emoji_mode":"default",
	// "has_invited":true,
	// "has_uploaded":false,
	// "has_created_channel":true,
	// "search_exclude_channels":"",
	// "messages_theme":"default",
	// "webapp_spellcheck":true,
	// "no_joined_overlays":false,
	// "no_created_overlays":true,
	// "dropbox_enabled":false,
	// "seen_user_menu_tip_card":true,
	// "seen_team_menu_tip_card":true,
	// "seen_channel_menu_tip_card":true,
	// "seen_message_input_tip_card":true,
	// "seen_channels_tip_card":true,
	// "seen_domain_invite_reminder":false,
	// "seen_member_invite_reminder":false,
	// "seen_flexpane_tip_card":true,
	// "seen_search_input_tip_card":true,
	// "mute_sounds":false,
	// "arrow_history":false,
	// "tab_ui_return_selects":true,
	// "obey_inline_img_limit":true,
	// "new_msg_snd":"knock_brush.mp3",
	// "collapsible":false,
	// "collapsible_by_click":true,
	// "require_at":false,
	// "mac_ssb_bounce":"",
	// "mac_ssb_bullet":true,
	// "win_ssb_bullet":true,
	// "expand_non_media_attachments":true,
	// "show_typing":true,
	// "pagekeys_handled":true,
	// "last_snippet_type":"",
	// "display_real_names_override":0,
	// "time24":false,
	// "enter_is_special_in_tbt":false,
	// "graphic_emoticons":false,
	// "convert_emoticons":true,
	// "autoplay_chat_sounds":true,
	// "ss_emojis":true,
	// "sidebar_behavior":"",
	// "mark_msgs_read_immediately":true,
	// "start_scroll_at_oldest":true,
	// "snippet_editor_wrap_long_lines":false,
	// "ls_disabled":false,
	// "sidebar_theme":"default",
	// "sidebar_theme_custom_values":"",
	// "f_key_search":false,
	// "k_key_omnibox":true,
	// "speak_growls":false,
	// "mac_speak_voice":"com.apple.speech.synthesis.voice.Alex",
	// "mac_speak_speed":250,
	// "comma_key_prefs":false,
	// "at_channel_suppressed_channels":"",
	// "push_at_channel_suppressed_channels":"",
	// "prompted_for_email_disabling":false,
	// "full_text_extracts":false,
	// "no_text_in_notifications":false,
	// "muted_channels":"",
	// "no_macssb1_banner":false,
	// "privacy_policy_seen":true,
	// "search_exclude_bots":false,
	// "fuzzy_matching":false
}

// UserDetails contains user details coming in the initial response from StartRTM
type UserDetails struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Created        JSONTime  `json:"created"`
	ManualPresence string    `json:"manual_presence"`
	Prefs          UserPrefs `json:"prefs"`
}

// JSONTime exists so that we can have a String method converting the date
type JSONTime int64

// String converts the unix timestamp into a string
func (t JSONTime) String() string {
	tm := t.Time()
	return fmt.Sprintf("\"%s\"", tm.Format("Mon Jan _2"))
}

// Time returns a `time.Time` representation of this value.
func (t JSONTime) Time() time.Time {
	return time.Unix(int64(t), 0)
}

// Team contains details about a team
type Team struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

// Icons XXX: needs further investigation
type Icons struct {
	Image48 string `json:"image_48"`
}

// Bot contains information about a bot
type Bot struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Deleted bool   `json:"deleted"`
	Icons   Icons  `json:"icons"`
}

// Info contains various details about Users, Channels, Bots and the authenticated user.
// It is returned by StartRTM or included in the "ConnectedEvent" RTM event.
type Info struct {
	URL      string       `json:"url,omitempty"`
	User     *UserDetails `json:"self,omitempty"`
	Team     *Team        `json:"team,omitempty"`
	Users    []User       `json:"users,omitempty"`
	Channels []Channel    `json:"channels,omitempty"`
	Groups   []Group      `json:"groups,omitempty"`
	Bots     []Bot        `json:"bots,omitempty"`
	IMs      []IM         `json:"ims,omitempty"`
}

type infoResponseFull struct {
	Info
	WebResponse
}

// GetBotByID returns a bot given a bot id
func (info Info) GetBotByID(botID string) *Bot {
	for _, bot := range info.Bots {
		if bot.ID == botID {
			return &bot
		}
	}
	return nil
}

// GetUserByID returns a user given a user id
func (info Info) GetUserByID(userID string) *User {
	for _, user := range info.Users {
		if user.ID == userID {
			return &user
		}
	}
	return nil
}

// GetChannelByID returns a channel given a channel id
func (info Info) GetChannelByID(channelID string) *Channel {
	for _, channel := range info.Channels {
		if channel.ID == channelID {
			return &channel
		}
	}
	return nil
}

// GetGroupByID returns a group given a group id
func (info Info) GetGroupByID(groupID string) *Group {
	for _, group := range info.Groups {
		if group.ID == groupID {
			return &group
		}
	}
	return nil
}
