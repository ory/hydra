#!/bin/bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

go fmt $(go list ./... | grep -v 'vendor')
go fmt .

git add -A
