#!/bin/bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

scripts/run-genswag.sh
git add -A
git diff --exit-code

./scripts/run-gensdk.sh
git add -A
git diff --exit-code
