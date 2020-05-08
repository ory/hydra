#!/bin/bash
# Exit script if you try to use an uninitialized variable.
set -o nounset

# Exit script if a statement returns a non-true return value.
set -o errexit

# Use the error status of the first failure, rather than that of the last item in a pipeline.
set -o pipefail

# The CLIDOC assignment is BASH specific i.e. not portable. Consider using this alternative assignment if $HOME does not expand properly
# CLIDOC=~/"docs/docs/reference/hydra_cli-reference.md"
# Add touch filename if .md does not exist

CLIDOC="$HOME/docs/docs/reference/hydra_cli-reference.md"

echo "Creating CLI Reference..."

# https://linux.die.net/abs-guide/here-docs.html
# turn this in an in mem here-doc to avoid /tmp out
# exec 9<<EOF
# cat <&9 >$OUT

cat <<EOF >$CLIDOC
---
id: hydra_cli-reference
title: Ory Hydra Command Line Interface Reference
---

Ory CLI Reference as of $(date)
# Command overview
$(hydra help)

# Command Reference clients
$(hydra help clients)

# Command Reference keys
$(hydra help keys)

# Command Reference migrate
$(hydra help migrate)

# Command Reference serve
$(hydra help serve)

# Command Reference token
$(hydra help token)

# Command Reference version
$(hydra help version)
EOF

echo "writing CLI Reference to $CLIDOC..."
