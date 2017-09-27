#!/bin/bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

java -jar scripts/swagger-codegen-cli-2.2.3.jar generate -i ./docs/api.swagger.json -l go -o ./sdk/go/swagger
java -jar scripts/swagger-codegen-cli-2.2.3.jar generate -i ./docs/api.swagger.json -l javascript -o ./sdk/js/swagger

scripts/run-format.sh