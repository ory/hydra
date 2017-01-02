#!/bin/bash

export PATH=PATH:$GOPATH/bin

# Test clients
hydra clients create --id foobar
hydra clients delete foobar

# Test token on arbitrary endpoints
curl --header "Authorization: bearer $(hydra token client)" http://localhost:4444/clients

# Test token validation
hydra token validate $(hydra token client)