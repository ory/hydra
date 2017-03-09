package main

import (
	"fmt"
	"os"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/urfave/cli"
)

var startCommand = cli.Command{
	Name:  "start",
	Usage: "executes the user defined process in a created container",
	ArgsUsage: `<container-id> [container-id...]

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The start command executes the user defined process in a created container .`,
	Action: func(context *cli.Context) error {
		if err := checkArgs(context, 1, minArgs); err != nil {
			return err
		}
		hasError := false
		factory, err := loadFactory(context)
		if err != nil {
			return err
		}

		for _, id := range context.Args() {
			container, err := factory.Load(id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "container %s does not exist\n", id)
				hasError = true
				continue
			}
			status, err := container.Status()
			if err != nil {
				fmt.Fprintf(os.Stderr, "status for %s: %v\n", id, err)
				hasError = true
				continue
			}
			switch status {
			case libcontainer.Created:
				if err := container.Exec(); err != nil {
					fmt.Fprintf(os.Stderr, "start for %s failed: %v\n", id, err)
					hasError = true
				}
			case libcontainer.Stopped:
				fmt.Fprintln(os.Stderr, "cannot start a container that has run and stopped")
				hasError = true
			case libcontainer.Running:
				fmt.Fprintln(os.Stderr, "cannot start an already running container")
				hasError = true
			default:
				fmt.Fprintf(os.Stderr, "cannot start a container in the %s state\n", status)
				hasError = true
			}
		}

		if hasError {
			return fmt.Errorf("one or more of container start failed")
		}
		return nil
	},
}
