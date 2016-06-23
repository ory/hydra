package main

import (
	"flag"
	"fmt"

	"github.com/nlopes/slack"
)

/*
   WARNING: This example is destructive in the sense that it create a channel called testpinning
*/
func main() {
	var (
		apiToken string
		debug    bool
	)

	flag.StringVar(&apiToken, "token", "YOUR_TOKEN_HERE", "Your Slack API Token")
	flag.BoolVar(&debug, "debug", false, "Show JSON output")
	flag.Parse()

	api := slack.New(apiToken)
	if debug {
		api.SetDebug(true)
	}

	var (
		postAsUserName  string
		postAsUserID    string
		postToChannelID string
	)

	// Find the user to post as.
	authTest, err := api.AuthTest()
	if err != nil {
		fmt.Printf("Error getting channels: %s\n", err)
		return
	}

	channelName := "testpinning"

	// Post as the authenticated user.
	postAsUserName = authTest.User
	postAsUserID = authTest.UserID

	// Create a temporary channel
	channel, err := api.CreateChannel(channelName)

	if err != nil {
		// If the channel exists, that means we just need to unarchive it
		if err.Error() == "name_taken" {
			err = nil
			channels, err := api.GetChannels(false)
			if err != nil {
				fmt.Println("Could not retrieve channels")
				return
			}
			for _, archivedChannel := range channels {
				if archivedChannel.Name == channelName {
					if archivedChannel.IsArchived {
						err = api.UnarchiveChannel(archivedChannel.ID)
						if err != nil {
							fmt.Printf("Could not unarchive %s: %s\n", archivedChannel.ID, err)
							return
						}
					}
					channel = &archivedChannel
					break
				}
			}
		}
		if err != nil {
			fmt.Printf("Error setting test channel for pinning: %s\n", err)
			return
		}
	}
	postToChannelID = channel.ID

	fmt.Printf("Posting as %s (%s) in channel %s\n", postAsUserName, postAsUserID, postToChannelID)

	// Post a message.
	postParams := slack.PostMessageParameters{}
	channelID, timestamp, err := api.PostMessage(postToChannelID, "Is this any good?", postParams)
	if err != nil {
		fmt.Printf("Error posting message: %s\n", err)
		return
	}

	// Grab a reference to the message.
	msgRef := slack.NewRefToMessage(channelID, timestamp)

	// Add message pin to channel
	if err := api.AddPin(channelID, msgRef); err != nil {
		fmt.Printf("Error adding pin: %s\n", err)
		return
	}

	// List all of the users pins.
	listPins, _, err := api.ListPins(channelID)
	if err != nil {
		fmt.Printf("Error listing pins: %s\n", err)
		return
	}
	fmt.Printf("\n")
	fmt.Printf("All pins by %s...\n", authTest.User)
	for _, item := range listPins {
		fmt.Printf(" > Item type: %s\n", item.Type)
	}

	// Remove the pin.
	err = api.RemovePin(channelID, msgRef)
	if err != nil {
		fmt.Printf("Error remove pin: %s\n", err)
		return
	}

	if err = api.ArchiveChannel(channelID); err != nil {
		fmt.Printf("Error archiving channel: %s\n", err)
		return
	}

}
