package cmd

var JanitorCmd = cmdHandler.Janitor.Command()

func init() {
	RootCmd.AddCommand(JanitorCmd)
}
