#!/usr/bin/env bash

set -Eeuxo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

killall hydra || true

export HYDRA_URL=http://127.0.0.1:4444/
export HYDRA_ADMIN_URL=http://127.0.0.1:4445/
export OAUTH2_CLIENT_ID=foobar-$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)
export OAUTH2_CLIENT_SECRET=bazbar
export OAUTH2_ISSUER_URL=http://127.0.0.1:4444/
export LOG_LEVEL=debug
export REDIRECT_URL=http://127.0.0.1:4445/callback
export OAUTH2_SCOPE=openid,offline

go install .
go build -buildmode=plugin -o memtest.so ./test/plugin

DATABASE_URL=memtest:// \
    DATABASE_PLUGIN=memtest.so \
    OAUTH2_CONSENT_URL=http://127.0.0.1:3000/consent \
    OAUTH2_LOGIN_URL=http://127.0.0.1:3000/login \
    OAUTH2_ERROR_URL=http://127.0.0.1:3000/error \
    OAUTH2_SHARE_ERROR_DEBUG=true \
    OAUTH2_ACCESS_TOKEN_STRATEGY=jwt \
    hydra serve all --dangerous-force-http --disable-telemetry &

while ! echo exit | nc 127.0.0.1 4444; do sleep 1; done
while ! echo exit | nc 127.0.0.1 4445; do sleep 1; done

hydra clients create \
    --endpoint http://127.0.0.1:4445 \
    --id $OAUTH2_CLIENT_ID \
    --secret $OAUTH2_CLIENT_SECRET \
    --response-types token,code,id_token \
    --grant-types refresh_token,authorization_code,client_credentials \
    --scope openid,offline \
    --callbacks http://127.0.0.1:5555/callback

token=$(hydra token client)

hydra token introspect "$token"

hydra clients delete $OAUTH2_CLIENT_ID

kill %1

rm memtest.so

exit 0

sleep 5
