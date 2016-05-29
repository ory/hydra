package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect with a cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if u := input("Cluster URL [" + c.ClusterURL + "]: "); u != "" {
			c.ClusterURL = u
		}
		if u := input("Client ID: "); u != "" {
			c.ClientID = u
		}
		if u := input("Client Secret: "); u != "" {
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
		fatal("Could not read user input because %s.", err)
	}
	return strings.TrimSpace(s)
}

func init() {
	RootCmd.AddCommand(connectCmd)
}
