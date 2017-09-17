#!/bin/bash

set -euo pipefail

goimports -w $(go list -f {{.Dir}} ./... | grep -v vendor | grep -v hydra$)