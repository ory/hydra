#!/usr/bin/env bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

hydra clients create --id foobar
hydra clients delete foobar
curl --header "Authorization: bearer $(hydra token client)" http://localhost:4444/clients
hydra token validate $(hydra token client)
