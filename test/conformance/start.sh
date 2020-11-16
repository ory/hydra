#!/bin/bash

set -euxo pipefail
cd "$( dirname "${BASH_SOURCE[0]}" )/../.."

# Used by run_test.go to rotate keys.
go build -o test/conformance/hydra/hydra -tags sqlite .

rm test/conformance/etc/sqlite/*.sqlite || true
docker-compose -f quickstart.yml -f quickstart-postgres.yml -f test/conformance/docker-compose.yml up ${1:-} -d
