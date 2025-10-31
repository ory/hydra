#!/bin/bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

toformat=$(goimports -l $(go list -f {{.Dir}} ./... | grep -v vendor | grep -v 'fosite$'))
[ -z "$toformat" ] && echo "All files are formatted correctly"
[ -n "$toformat" ] && echo "Please use \`goimports\` to format the following files:" && echo $toformat && exit 1

exit 0
