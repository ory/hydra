#!/bin/bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

goimports -w $(go list -f {{.Dir}} ./... | grep -v vendor | grep -v "ory.hydra$")
goimports -w *.go

git add -A
