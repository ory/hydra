package cmd

import (
	"fmt"
	"os"
)

func fatal(message string, args ...interface{}) {
	fmt.Printf(message+"\n", args...)
	os.Exit(1)
}
