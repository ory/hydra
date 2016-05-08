package cmd

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite/handler/core"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"github.com/spf13/cobra"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Start the hydra host service",
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
	var c = new(configuration)

	fmt.Println("Connecting to backend...")
	clientStore := newClientStore(c)
	fositeStore := newFositeStore(c, clientStore)
	ladonStore := newLadonStore(c)
	hmacStrategy := newHmacStrategy(c)
	keyManager := newKeyManager(c)
	idStrategy := newIdStrategy(c, keyManager)
	hasher := newHasher(c)
	fosite := newFosite(c, hmacStrategy, idStrategy, fositeStore, hasher)
	fositeHandler := newOAuth2Handler(c, fosite, keyManager)
	fmt.Println("Successfully connected to all backends.")

	if err := createAdminIfNotExists(clientStore, ladonStore, hasher); err != nil {
		fatal("%s", err.Error())
	}

	ladonWarden := &ladon.Ladon{Manager: ladonStore}
	localWarden := &warden.LocalWarden{
		Warden: ladonWarden,
		TokenValidator: &core.CoreValidator{
			AccessTokenStrategy: hmacStrategy,
			AccessTokenStorage:  fositeStore,
		},
		Issuer: c.GetIssuer(),
	}

	fmt.Println("Setting up routes...")
	router := httprouter.New()
	clientHandler := &client.Handler{
		Manager: clientStore,
		H:       &herodot.JSON{},
		W:       localWarden,
	}
	clientHandler.SetRoutes(router)
	fositeHandler.SetRoutes(router)

	fmt.Printf("Starting server on %s\n", c.GetAddress())
	fatal("Could not start server because %s.\n", http.ListenAndServe(c.GetAddress(), router))
}
