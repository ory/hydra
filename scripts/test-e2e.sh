#!/usr/bin/env bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

DATABASE_URL=memory hydra host --dangerous-auto-logon --dangerous-force-http --disable-telemetry &
while ! echo exit | nc 127.0.0.1 4444; do sleep 1; done

hydra clients create --id foobar
hydra clients delete foobar
curl --header "Authorization: bearer $(hydra token client)" http://localhost:4444/clients
hydra token validate $(hydra token client)
