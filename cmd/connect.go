package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect with a cluster",
	Run: func(cmd *cobra.Command, args []string) {
		secret := "*********"
		fmt.Println("To keep the current value, press enter.")

		if u, _ := cmd.Flags().GetString("url"); u != "" {
			c.ClusterURL = u
		} else if u := input("Cluster URL [" + c.ClusterURL + "]: "); u != "" {
			c.ClusterURL = u
		}

		if u, _ := cmd.Flags().GetString("id"); u != "" {
			c.ClientID = u
		} else if u := input("Client ID [" + c.ClientID + "]: "); u != "" {
			c.ClientID = u
		}

		if c.ClientSecret == "" {
			secret = "empty"
		}

		if u, _ := cmd.Flags().GetString("secret"); u != "" {
			logrus.Warn("You should not provide secrets using command line flags. The secret might leak to bash history and similar systems.")
			c.ClientSecret = u
		} else if u := input("Client Secret [" + secret + "]: "); u != "" {
			c.ClientSecret = u
		}
		if err := c.Persist(); err != nil {
			log.Fatalf("Unable to save config file because %s.", err)
		}
		fmt.Println("Done.")
	},
}

func input(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	s, err := reader.ReadString('\n')
	if err != nil {
		fatal(fmt.Sprintf("Could not read user input because %s.", err))
	}
	return strings.TrimSpace(s)
}

func init() {
	RootCmd.AddCommand(connectCmd)
	connectCmd.Flags().String("url", "", "The cluster URL")
	connectCmd.Flags().String("id", "", "The client id")
	connectCmd.Flags().String("secret", "", "The client secret")
}
