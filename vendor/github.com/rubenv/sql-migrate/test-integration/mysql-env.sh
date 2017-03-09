#!/bin/bash

# Tweak PATH for Travis
export PATH=$PATH:$HOME/gopath/bin

export MYSQL_USER=root
export DATABASE=test

OPTIONS="-config=test-integration/dbconfig.yml -env mysql_env"

set -ex

sql-migrate status $OPTIONS
sql-migrate up $OPTIONS
sql-migrate down $OPTIONS
sql-migrate redo $OPTIONS
sql-migrate status $OPTIONS
