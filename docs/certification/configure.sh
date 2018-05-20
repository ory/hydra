#!/bin/bash

source $HOME/.profile

hydra clients delete \
    --endpoint http://localhost:9000 \
    test-client

hydra clients create \
    --endpoint http://localhost:9000 \
    --id test-client \
    --secret test-secret \
    --response-types code,id_token \
    --grant-types refresh_token,authorization_code \
    --scope openid,offline,offline_access,profile,email,address,phone \
    --callbacks https://op.certification.openid.net:60848/authz_cb
