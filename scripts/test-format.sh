#!/bin/bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

gometalinter --disable-all --enable=gosec --enable=gofmt --vendor ./...

exit 0