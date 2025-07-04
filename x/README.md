# ory/x

[![GoDoc reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/ory/x)
[![tests](https://github.com/ory/x/actions/workflows/test.yml/badge.svg)](https://github.com/ory/x/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/ory/x/badge.svg?branch=master)](https://coveralls.io/github/ory/x?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ory/x)](https://goreportcard.com/report/github.com/ory/x)

Shared libraries used in the ORY ecosystem. Use at your own risk. Breaking
changes should be anticipated.

## Run tests under Wine

Install [Wine](https://www.winehq.org/) and then for a given package e.g.
`./jsonnetsecure`:

```sh
# Need to compile the jsonnet program for Windows since it is required by some tests.
$ GOOS=windows GOARCH=amd64 go build -o ./jsonnet.exe github.com/ory/x/jsonnetsecure/cmd
$ GOOS=windows GOARCH=amd64 go test -c ./jsonnetsecure
$ ORY_JSONNET_PATH=$PWD/jsonnet.exe WINEDEBUG=-all wine  $PWD/jsonnetsecure.test.exe
```

_Note: Wine only emulates Windows amd64 so it requires Rosetta on aarch64
macOS._
