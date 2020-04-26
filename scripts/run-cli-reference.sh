#!/bin/bash

CLIDOC=./docs/docs/cli-reference.md

echo "Creating CLI Reference..."

# https://linux.die.net/abs-guide/here-docs.html
# turn this in an in mem here-doc to avoid /tmp out
# exec 9<<EOF
# cat <&9 >$OUT

cat <<EOF >$CLIDOC
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

echo "writting CLI Reference to $CLIDOC...


