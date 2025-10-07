#!/usr/bin/env bash

set -euxo pipefail

schema_version="${1:-$(git rev-parse --short HEAD)}"

sed "s!ory://tracing-config!https://raw.githubusercontent.com/ory/hydra/$schema_version/oryx/otelx/config.schema.json!g;" spec/config.json > .schema/config.schema.json

git commit --author="ory-bot <60093411+ory-bot@users.noreply.github.com>" -m "autogen: render config schema" .schema/config.schema.json || true
