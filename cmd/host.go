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
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/cmd/server"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Start the hydra host service",
	Run:   runHostCmd,
}

var c config.Config

func init() {
	RootCmd.AddCommand(hostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	hostCmd.Flags().BoolP("save-credentials", "s", false, "If no admin exists, hydra creates a new one. Use this option to save the credentials to the hydra config file.")
}

func runHostCmd(cmd *cobra.Command, args []string) {
	router := httprouter.New()
	serverHandler := &server.Handler{}
	serverHandler.Listen(c, router)

	fmt.Println("Connecting to backend...")
	hasher := newHasher(c)
	keyManager := newKeyManager(c)
	ladonStore := newLadonStore(c)
	hmacStrategy := newHmacStrategy(c)
	clientStore := newClientStore(c, hasher)
	fositeStore := newFositeStore(c, clientStore)
	idStrategy := newIdStrategy(c, keyManager)
	fosite := newFosite(c, hmacStrategy, idStrategy, fositeStore, hasher)
	fositeHandler := newOAuth2Handler(c, fosite, keyManager)
	fmt.Println("Successfully connected to all backends.")

	save, _ := cmd.Flags().GetBool("save-credentials")
	if err := createAdminIfNotExists(clientStore, ladonStore, save); err != nil {
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
