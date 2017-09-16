#!/usr/bin/env bash

set -euo pipefail

swagger generate spec -m -o ./docs/api.swagger.json
