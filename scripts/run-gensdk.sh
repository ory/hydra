#!/bin/bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

swagger-codegen generate -i ./docs/api.swagger.json -l go -o ./sdk/go/swagger
swagger-codegen generate -i ./docs/api.swagger.json -l javascript -o ./sdk/js/swagger
