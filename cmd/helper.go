package cmd

import (
	"fmt"
	"os"
	"github.com/ory-am/common/rand/sequence"
)

func fatal(message string, args ...interface{}) {
	fmt.Printf(message + "\n", args...)
	os.Exit(1)
}

var secretCharSet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-.,:;$%!&/()=?+*#<>")

func generateSecret(length int) []byte {
	secret, err := sequence.RuneSequence(length, secretCharSet)
	if err != nil {
		fatal("Could not generated random secret because %s", err)
		return []byte{}
	}
	return []byte(string(secret))
}