# pagination

A simple helper for dealing with pagination.

```
go get github.com/ory/pagination
```

## Example

```go
package main

import (
	"github.com/ory/pagination"
    "net/http"
    "net/url"
    "fmt"
)

func main() {
	u, _ := url.Parse("http://localhost/foo?offset=0&limit=10")
    limit, offset := pagination.Parse(&http.Request{URL: u}, 5, 5, 10)

    items := []string{"a", "b", "c", "d"}
    start, end := pagination.Index(limit, offset, len(items))
    fmt.Printf("Got items: %v", items[start:end])
}
```
