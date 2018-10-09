# graceful

[![Build Status](https://travis-ci.org/ory/graceful.svg?branch=master)](https://travis-ci.org/ory/graceful)
[![Coverage Status](https://coveralls.io/repos/github/ory/graceful/badge.svg?branch=master)](https://coveralls.io/github/ory/graceful?branch=master)
[![Docs: GoDoc](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/ory/graceful)

Best practice http server configurations and helpers for Go 1.8's http graceful shutdown feature. Currently supports
best practice configurations by:

* [Cloudflare](https://blog.cloudflare.com/exposing-go-on-the-internet/)

## Usage

To install this library, do:

```sh
go get github.com/ory/graceful
```

### Running Cloudflare Config with Graceful Shutdown

```go
package main

import (
    "net/http"
    "log"

    "github.com/ory/graceful"
)

func main() {
    server := graceful.WithDefaults(&http.Server{
        Addr: ":54932",
        // Handler: someHandler,
    })

    log.Println("main: Starting the server")
    if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
        log.Fatalln("main: Failed to gracefully shutdown")
    }
    log.Println("main: Server was shutdown gracefully")
}
```
