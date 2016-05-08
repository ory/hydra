package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/ladon"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/fosite/handler/core"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"log"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runHostCmd,
}

func init() {
	RootCmd.AddCommand(hostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hostCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func runHostCmd(cmd *cobra.Command, args []string) {
	var c config

	fmt.Println("Connecting to backend...")
	fositeStore := newFositeStore(c)
	fmt.Println("Successfully connected to fosite backend.")
	ladonStore := newLadonStore(c)
	fmt.Println("Successfully connected to ladon backend.")
	clientStore := newClientStore(c)
	fmt.Println("Successfully connected to client backends.")
	hmacStrategy := newHmacStrategy(c)
	fmt.Println("Successfully connected to all backends.")

	ladonWarden := &ladon.Ladon{Manager: ladonStore}
	localWarden := &warden.LocalWarden{
		Warden: ladonWarden,
		TokenValidator: &core.CoreValidator{
			AccessTokenStrategy: hmacStrategy,
			AccessTokenStorage:  fositeStore,
		},
		Issuer: c.Iss(),
	}

	fmt.Println("Setting up routes...")
	router := httprouter.New()
	clientHandler := &client.Handler{
		Manager: clientStore,
		H:       &herodot.JSON{},
		W:       localWarden,
	}
	clientHandler.SetRoutes(router)

	fmt.Printf("Starting server on %s.\n", c.Addr())
	log.Fatalf("Could not start server because %s.\n", http.ListenAndServe(c.Addr(), router))
}
