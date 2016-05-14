package cmd

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/cmd/server"
	"github.com/ory-am/hydra/pkg"
	"github.com/spf13/cobra"
	"github.com/Sirupsen/logrus"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Start the hydra host service",
	Run:   runHostCmd,
}

func init() {
	RootCmd.AddCommand(hostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	hostCmd.Flags().Bool("dangerous-auto-logon", false, "Stores the root credentials in ~/.hydra.yml. Do not use in production.")
}

func runHostCmd(cmd *cobra.Command, args []string) {
	router := httprouter.New()
	serverHandler := &server.Handler{}
	serverHandler.Start(c, router)

	if ok, _ := cmd.Flags().GetBool("dangerous-auto-logon"); ok {
		logrus.Warnln("Do not use flag --dangerous-auto-logon in production.")
		err :=c.Persist()
		pkg.Must(err, "Could not write configuration file: ", err)
	}

	logrus.Infof("Starting server on %s", c.GetAddress())
	err := http.ListenAndServe(c.GetAddress(), router)
	pkg.Must(err, "Could not start server because %s.", err)
}
