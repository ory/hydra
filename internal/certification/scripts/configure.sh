#!/bin/bash

# shellcheck disable=SC1090,SC1091
source "$HOME"/.profile

hydra clients delete \
    --endpoint http://localhost:9000 \
    test-client

hydra clients create \
    --endpoint http://localhost:9000 \
    --id test-client \
    --secret test-secret \
    --response-types code,id_token,token \
    --grant-types implicit,refresh_token,authorization_code \
    --scope openid,offline,offline_access,profile,email,address,phone \
    --callbacks https://op.certification.openid.net:60848/authz_cb,https://op.certification.openid.net:61290/authz_cb,https://op.certification.openid.net:61291/authz_cb,https://op.certification.openid.net:61302/authz_cb,https://op.certification.openid.net:61303/authz_cb,https://op.certification.openid.net:61304/authz_cb
