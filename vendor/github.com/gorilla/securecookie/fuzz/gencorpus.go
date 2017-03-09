package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"testing/quick"

	"github.com/gorilla/securecookie"
)

var hashKey = []byte("very-secret12345")
var blockKey = []byte("a-lot-secret1234")
var s = securecookie.New(hashKey, blockKey)

type Cookie struct {
	B bool
	I int
	S string
}

func main() {
	var c Cookie
	t := reflect.TypeOf(c)
	rnd := rand.New(rand.NewSource(0))
	for i := 0; i < 100; i++ {
		v, ok := quick.Value(t, rnd)
		if !ok {
			panic("couldn't generate value")
		}
		encoded, err := s.Encode("fuzz", v.Interface())
		if err != nil {
			panic(err)
		}
		f, err := os.Create(fmt.Sprintf("corpus/%d.sc", i))
		if err != nil {
			panic(err)
		}
		_, err = io.WriteString(f, encoded)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}
