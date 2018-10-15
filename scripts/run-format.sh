#!/bin/bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

go get -u golang.org/x/tools/cmd/goimports

files=$(go list -f {{.Dir}} ./... | grep -v vendor | grep -v hydra$)

gofmt -w -s $files
goimports -local github.com/ory/hydra -w $files

git add -A
