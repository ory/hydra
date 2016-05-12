package pkg

import (
	"fmt"
	"os"
)

func Must(err error, message string, args ...interface{}) {
	if err != nil {
		return
	}
	fmt.Fprint(os.Stderr, message+"\n", args...)
	os.Exit(0)
}