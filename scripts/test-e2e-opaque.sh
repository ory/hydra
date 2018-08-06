#!/usr/bin/env bash

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

killall hydra || true
killall mock-lcp || true
killall mock-cb || true

export HYDRA_URL=http://127.0.0.1:4444/
export HYDRA_ADMIN_URL=http://127.0.0.1:4445/
export OAUTH2_CLIENT_ID=foobar
export OAUTH2_CLIENT_SECRET=bazbar
export OAUTH2_ISSUER_URL=http://127.0.0.1:4444/
export LOG_LEVEL=debug
export REDIRECT_URL=http://127.0.0.1:5555/callback
export AUTH2_SCOPE=openid,offline

go install .
go install ./test/mock-client
go install ./test/mock-lcp
go install ./test/mock-cb

DATABASE_URL=memory \
    OAUTH2_CONSENT_URL=http://127.0.0.1:3000/consent \
    OAUTH2_LOGIN_URL=http://127.0.0.1:3000/login \
    OAUTH2_ERROR_URL=http://127.0.0.1:3000/error \
    OAUTH2_SHARE_ERROR_DEBUG=true \
    hydra serve all --dangerous-force-http --disable-telemetry &

PORT=3000 mock-lcp &

PORT=5555 mock-cb &

while ! echo exit | nc 127.0.0.1 4444; do sleep 1; done
while ! echo exit | nc 127.0.0.1 4445; do sleep 1; done
while ! echo exit | nc 127.0.0.1 5555; do sleep 1; done
while ! echo exit | nc 127.0.0.1 3000; do sleep 1; done


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

## Authenticate but do not remember user
cookie=$(OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept" \
    mock-client)
export AUTH_COOKIE=$cookie

## Must fail because prompt=none but no session was remembered
if OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept&prompt=none" \
    mock-client; then
    echo "should have failed"
    exit 1
fi

# Authenticate and remember login (but not consent)
cookie=$(OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept&rememberLogin=yes" \
    mock-client)
export AUTH_COOKIE=$cookie

## Must fail because prompt=none but consent was not yet stored
if OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept&prompt=none" \
    mock-client; then
    echo "should have failed"
    exit 1
fi

# Remember consent
cookie=$(OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept&rememberConsent=yes" \
    mock-client)

## Prompt none should work now because cookie was set
OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept&prompt=none" \
    mock-client

hydra clients delete $OAUTH2_CLIENT_ID

kill %1
kill %2
kill %3
exit 0

sleep 5
