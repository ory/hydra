// +build linux

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var pauseCommand = cli.Command{
	Name:  "pause",
	Usage: "pause suspends all processes inside the container",
	ArgsUsage: `<container-id> [container-id...]

Where "<container-id>" is the name for the instance of the container to be
paused. `,
	Description: `The pause command suspends all processes in the instance of the container.

Use runc list to identiy instances of containers and their current status.`,
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
			if err := container.Pause(); err != nil {
				fmt.Fprintf(os.Stderr, "pause container %s : %s\n", id, err)
				hasError = true
			}
		}

		if hasError {
			return fmt.Errorf("one or more of container pause failed")
		}
		return nil
	},
}

var resumeCommand = cli.Command{
	Name:  "resume",
	Usage: "resumes all processes that have been previously paused",
	ArgsUsage: `<container-id> [container-id...]

Where "<container-id>" is the name for the instance of the container to be
resumed.`,
	Description: `The resume command resumes all processes in the instance of the container.

Use runc list to identiy instances of containers and their current status.`,
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
			if err := container.Resume(); err != nil {
				fmt.Fprintf(os.Stderr, "resume container %s : %s\n", id, err)
				hasError = true
			}
		}

		if hasError {
			return fmt.Errorf("one or more of container resume failed")
		}
		return nil
	},
}
