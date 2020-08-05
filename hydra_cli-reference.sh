#!/bin/bash
# Exit script if you try to use an uninitialized variable.
set -o nounset

# Exit script if a statement returns a non-true return value.
set -o errexit

# Use the error status of the first failure, rather than that of the last item in a pipeline.
#set -o pipefail

# The CLIDOC assignment is BASH specific i.e. not portable. Consider using this alternative assignment if $HOME does not expand properly
# CLIDOC=~/"docs/docs/reference/hydra_cli-reference.md"
# Add touch filename if .md does not exist

CLIDOC="./docs/docs/reference/hydra_cli-reference.md"

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
---
# Command overview
---
$(hydra help)



---
# Command Reference clients
---
$(hydra help clients)
---
## Command Reference clients create
---
$(hydra help clients create)   # Create a new OAuth 2.0 Client
---
## Command Reference clients delete
---
$(hydra help clients delete)   # Delete an OAuth 2.0 Client
---
## Command Reference clients get
---
$(hydra help clients get)      # Get an OAuth 2.0 Client
---
## Command Reference clients import
---
$(hydra help clients import)   # Import OAuth 2.0 Clients from one or more JSON files
---
## Command Reference clients list
---
$(hydra help clients list)     # List OAuth 2.0 Clients



---
# Command Reference keys
---
$(hydra help keys)
---
## Command Reference keys create
---
$(hydra help keys create)      # Create a new JSON Web Key Set
---
## Command Reference keys delete
---
$(hydra help keys delete)      # Delete a new JSON Web Key Set
---
## Command Reference keys get
---
$(hydra help keys get)         # Get a new JSON Web Key Set
---
## Command Reference keys get import
---
$(hydra help keys import)      # Imports cryptographic keys of any format to the JSON Web Key Store




---
# Command Reference migrate
---
$(hydra help migrate)
---
## Command Reference migrate sql
---$(hydra help migrate sql)      # Create SQL schemas and apply migration plans



---
# Command Reference serve
---
$(hydra help serve)
---
## Command Reference serve admin
---
$(hydra help serve admin)      # Serves Administrative HTTP/2 APIs
---
## Command Reference serve all
---
$(hydra help serve all)        # Serves both public and administrative HTTP/2 APIs
---
## Command Reference serve public
---
$(hydra help serve public)     # Serves Public HTTP/2 APIs



---
# Command Reference token
---
$(hydra help token)
---
## Command Reference token client
---
$(hydra help token client)    # An exemplary OAuth 2.0 Client performing the OAuth 2.0 Client Credentials Flow
---
## Command Reference token flush
---
$(hydra help token flush)     # Removes inactive access tokens from the database
---
## Command Reference token introspect
---
$(hydra help token introspect)# Introspect an access or refresh token
---
## Command Reference token revoke
---
$(hydra help token revoke)    # Revoke an access or refresh token
---
## Command Reference token user
---
$(hydra help token user)      # An exemplary OAuth 2.0 Client performing the OAuth 2.0 Authorize Code Flow



---
# Command Reference version
---
$(hydra help version)         # Display this binary's version, build time and git hash of this build
    
EOF

echo "writing CLI Reference to $CLIDOC..."
