package pkg

import (
	"fmt"
	"os"
)

func MustArgs(expected, actual int) {
	if expected == actual {
		return
	}
	fmt.Fprintf(os.Stderr, "Invalid number of arguments. Expected %d but got %d.", expected, actual)
	os.Exit(0)
}
