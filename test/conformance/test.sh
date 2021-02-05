#!/bin/bash

set -euxo pipefail
cd "$( dirname "${BASH_SOURCE[0]}" )"

go test -tags conformity -test.timeout 60m -failfast "$@" .
