#!/bin/bash

set -euxo pipefail
cd "$( dirname "${BASH_SOURCE[0]}" )"

docker build -t oryd/hydra-oidc-server:latest .
docker build -t oryd/hydra-oidc-httpd:latest -f httpd/Dockerfile .

docker push oryd/hydra-oidc-server:latest
docker push oryd/hydra-oidc-httpd:latest
