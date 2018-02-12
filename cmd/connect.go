// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"net/http"
	"net/url"

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
			fmt.Println("You should not provide secrets using command line flags. The secret might leak to bash history and similar systems.")
			c.ClientSecret = u
		} else if u := input("Client Secret [" + secret + "]: "); u != "" {
			c.ClientSecret = u
		}

		if skipNewsletter, _ := cmd.Flags().GetBool("skip-newsletter"); !c.SignedUpForNewsletter && !skipNewsletter {
			u := "https://ory.us10.list-manage.com/subscribe/post?u=ffb1a878e4ec6c0ed312a3480&id=f605a41b53"
			fmt.Println("You are using the CLI for the first time. It is really important to keep your installation up to date. Because this technology is open source, we have no way of knowing who you are and how to contact you. Subscribe to our release and security announcements, and never miss important patches again:")
			m := input("[Enter Email Address]:")

			v := url.Values{}
			v.Add("EMAIL", m)

			req, err := http.NewRequest("POST", u, strings.NewReader(v.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			if m == "" {
			} else {
				if err != nil {
					fmt.Printf("There was some error: %s\n", err.Error())
				} else {
					resp, err := http.DefaultClient.Do(req)

					if err != nil {
						fmt.Printf("There was some error: %s\n", err.Error())
					} else {
						defer resp.Body.Close()

						fmt.Println("To complete the subscription process, please click the link in the email you just received.")
						c.SignedUpForNewsletter = true
					}
				}
			}
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
	connectCmd.Flags().Bool("skip-newsletter", false, "Skip the newsletter sign up question")
}
