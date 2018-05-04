#!/usr/bin/env bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

DATABASE_URL=memory hydra serve --dangerous-force-http --disable-telemetry &
while ! echo exit | nc 127.0.0.1 4444; do sleep 1; done

export HYDRA_URL=http://localhost:4444/
export OAUTH2_CLIENT_ID=foobar
export OAUTH2_CLIENT_SECRET=bazbar

hydra clients create --id $OAUTH2_CLIENT_ID --secret $OAUTH2_CLIENT_SECRET -g client_credentials
token=$(hydra token client)
hydra token introspect $token
hydra clients delete foobar
