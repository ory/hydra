#!/usr/bin/env bash

set -Eeuxo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."

killall hydra || true
killall mock-lcp || true
killall mock-cb || true

export HYDRA_URL=http://127.0.0.1:4444/
export HYDRA_ADMIN_URL=http://127.0.0.1:4445/
export OAUTH2_CLIENT_ID=foobar-$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)
export OAUTH2_CLIENT_SECRET=bazbar
export OAUTH2_ISSUER_URL=http://127.0.0.1:4444/
export LOG_LEVEL=debug
export REDIRECT_URL=http://127.0.0.1:5555/callback
export OAUTH2_SCOPE=openid,offline

go install .
go install ./test/mock-client
go install ./test/mock-lcp
go install ./test/mock-cb

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

clientToken=$(hydra token client)

if [[ $(hydra token introspect "$(echo $clientToken)") =~ "false" ]]; then
    echo "Token introspection should return true"
    exit 1
fi

## Authenticate but do not remember user
cookie=$(OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept" \
    mock-client -print-cookie)
export AUTH_COOKIE=$cookie

## Must fail because prompt=none but no session was remembered
if OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept&prompt=none" \
    mock-client; then
    echo "should have failed"
    exit 1
fi

# Authenticate and remember login (but not consent)
cookie=$(OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept&rememberLogin=yes" \
    mock-client -print-cookie)
export AUTH_COOKIE=$cookie

## Must fail because prompt=none but consent was not yet stored
if OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept&prompt=none" \
    mock-client; then
    echo "should have failed"
    exit 1
fi

# Remember consent
$(OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept&rememberConsent=yes" \
    mock-client)

## Prompt none should work now because cookie was set
userToken=$(OAUTH2_EXTRA="&mockLogin=accept&mockConsent=accept&prompt=none" \
    mock-client -print-token)

if [[ "$(hydra token introspect $userToken)" =~ "false" ]]; then
    echo "Token introspection should return true for the user token"
    exit 1
fi

if [[ "$(curl http://localhost:4445/oauth2/auth/sessions/consent/the-subject)" =~ "the-subject" ]]; then
    echo "Consent session looks good"
else
    echo "Checking consent session should show the-subject"
    exit 1
fi

curl -X DELETE http://localhost:4445/oauth2/auth/sessions/consent/the-subject

if [[ "$(hydra token introspect $userToken)" =~ "true" ]]; then
    echo "Token introspection should return false because the consent session was revoked"
    exit 1
fi

if [[ "$(hydra token introspect $clientToken)" =~ "false" ]]; then
    echo "Token introspection should return true"
    exit 1
fi

hydra clients delete $OAUTH2_CLIENT_ID

if [[ "$(hydra token introspect $clientToken)" =~ "true" ]]; then
    echo "Token introspection should return false because the client was deleted"
    exit 1
fi

hydra clients create \
    --endpoint http://127.0.0.1:4445 \
    --id $OAUTH2_CLIENT_ID \
    --secret $OAUTH2_CLIENT_SECRET \
    --response-types token,code,id_token \
    --grant-types refresh_token,authorization_code,client_credentials \
    --scope openid,offline \
    --callbacks http://127.0.0.1:5555/callback

if [[ "$(hydra token introspect $clientToken)" =~ "true" ]]; then
    echo "Token introspection should return false even after the client was re-created"
    exit 1
fi

kill %1
kill %2
kill %3
exit 0

sleep 5
