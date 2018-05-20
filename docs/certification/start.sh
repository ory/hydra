#!/bin/bash

source $HOME/.profile

domain=oidc-certification.ory.sh:8443
hydraport=9000
idport=9001

(cd ./hydra-login-consent-node; HYDRA_URL=http://localhost:$hydraport PORT=$idport npm start &)

PORT=$hydraport \
    OAUTH2_CONSENT_URL=https://$domain/consent \
    OAUTH2_LOGIN_URL=https://$domain/login \
    OAUTH2_ISSUER_URL=https://$domain/ \
    OAUTH2_SHARE_ERROR_DEBUG=1 \
    LOG_LEVEL=debug \
    DATABASE_URL=memory \
    hydra serve --dangerous-force-http

kill %1
