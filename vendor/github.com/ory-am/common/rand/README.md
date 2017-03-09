# rand
A library based on crypto/rand to create random sequences, which are cryptographically strong. See: [crypto/rand](http://golang.org/pkg/crypto/rand/)

## Install

Run `go get github.com/ory-am/common/rand`

## Usage

### Create a random integer

Create a random integer using [crypto/rand.Read](http://golang.org/pkg/crypto/rand/#Read):

```
import "github.com/ory-am/common/rand/numeric"
import "fmt"

func main() {
    fmt.Printf("%d", numeric.Int64())
    fmt.Printf("%d", numeric.UInt64())
    fmt.Printf("%d", numeric.Int32())
    fmt.Printf("%d", numeric.UInt32())
}
```

### Create a random rune sequence / string

Create a random string using [crypto/rand.Read](http://golang.org/pkg/crypto/rand/#Read):

```
import "github.com/ory-am/common/rand/sequence"
import "fmt"

func main() {
    allowed := []rune("abcdefghijklmnopqrstuvwxyz")
    length := 10
    seq, err := sequence.RuneSequence(length, allowed)

    fmt.Printf("%s", seq)
    fmt.Printf("%s", string(seq))
}
```