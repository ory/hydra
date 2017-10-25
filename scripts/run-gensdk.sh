#!/bin/bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

scripts/run-genswag.sh

rm -rf ./sdk/go/hydra/swagger
rm -rf ./sdk/js/swagger

# curl -O scripts/swagger-codegen-cli-2.2.3.jar http://central.maven.org/maven2/io/swagger/swagger-codegen-cli/2.2.3/swagger-codegen-cli-2.2.3.jar

java -jar scripts/swagger-codegen-cli-2.2.3.jar generate -i ./docs/api.swagger.json -l go -o ./sdk/go/hydra/swagger
java -jar scripts/swagger-codegen-cli-2.2.3.jar generate -i ./docs/api.swagger.json -l javascript -o ./sdk/js/swagger

scripts/run-format.sh

git checkout HEAD -- sdk/go/hydra/swagger/configuration.go
git checkout HEAD -- sdk/go/hydra/swagger/api_client.go
rm -f ./sdk/js/swagger/package.json
rm -rf ./sdk/js/swagger/test

npm run prettier
