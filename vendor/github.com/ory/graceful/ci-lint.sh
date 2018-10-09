#!/usr/bin/env bash

set -e

directories=$(glide novendor)
for i in $directories
do
  if [[ "$i" == "." ]]; then
    continue
  fi
  go vet $i
  golint $i
  goimports -d $(dirname $i)
done
