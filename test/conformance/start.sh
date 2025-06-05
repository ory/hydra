#!/bin/bash

set -euxo pipefail
cd "$( dirname "${BASH_SOURCE[0]}" )/../.."

# shellcheck disable=SC2086
docker compose -f quickstart.yml -f quickstart-postgres.yml -f test/conformance/docker-compose.yml up ${1:-} -d --force-recreate --build
docker ps -a
docker images
