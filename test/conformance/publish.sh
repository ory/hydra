#!/bin/bash

set -euxo pipefail
cd "$( dirname "${BASH_SOURCE[0]}" )"


docker buildx build --output type=docker --platform linux/amd64 -t oryd/hydra-oidc-server:latest .
docker buildx build --output type=docker --platform linux/amd64 -t oryd/hydra-oidc-httpd:latest -f httpd/Dockerfile .

docker push oryd/hydra-oidc-server:latest
docker push oryd/hydra-oidc-httpd:latest
