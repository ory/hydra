#!/bin/bash

set -euo pipefail

if [[ $(cat package.json | grep $(git describe --tag)) ]]; then
  echo "Seems like deploy script ran already."
  exit 0;
else
  echo "Running deploy script..."
fi

gox -ldflags "-X github.com/ory/hydra/cmd.Version=`git describe --tags` -X github.com/ory/hydra/cmd.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X github.com/ory/hydra/cmd.GitHash=`git rev-parse HEAD`" -output "dist/{{.Dir}}-{{.OS}}-{{.Arch}}"
npm version -f --no-git-tag-version $(git describe --tag)
