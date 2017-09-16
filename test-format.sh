#!/bin/bash

set -euo pipefail

toformat = $(goimports -l $(go list -f {{.Dir}} ./... | grep -v vendor | grep -v hydra$))
[ -n "$toformat" ] && echo $toformat && exit 1