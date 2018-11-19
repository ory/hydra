#!/bin/bash

set -Eeuxo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

swagger generate spec -m -o ./docs/api.swagger.json
