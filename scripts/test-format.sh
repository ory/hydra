#!/bin/bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

./scripts/run-format.sh
git diff --exit-code

exit 0