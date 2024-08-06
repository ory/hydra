#!/bin/bash

set -euxo pipefail
cd "$( dirname "${BASH_SOURCE[0]}" )/../.."

docker compose -f quickstart.yml -f quickstart-postgres.yml -f test/conformance/docker-compose.yml down -v
