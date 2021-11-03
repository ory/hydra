#!/bin/bash

# shellcheck disable=SC1090,SC1091
source "$HOME"/.profile

go get -u -d github.com/ory/hydra
go get -d -u github.com/devopsfaith/krakend-examples/gin
(cd "$HOME"/hydra-login-consent-node || exit; git pull -ff; npm i)
cd "$HOME" || exit
go install github.com/ory/hydra
go install github.com/devopsfaith/krakend-examples/gin
