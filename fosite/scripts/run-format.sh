#!/bin/bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

goimports -w $(go list -f {{.Dir}} ./... | grep -v vendor | grep -v fosite$)
goimports -w *.go
