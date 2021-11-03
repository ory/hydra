#!/bin/bash

# shellcheck disable=SC1090,SC1091
source "$HOME"/.profile

domain=oidc-certification.ory.sh:8443
hydraport=9000
idport=9001

docker start kong-database
docker start kong

dockerize -wait http://localhost:8001/ -timeout 30s

ip=$(curl ifconfig.co)

curl -i -X DELETE --url http://localhost:8001/apis/hydra-oauth
curl -i -X DELETE --url http://localhost:8001/apis/login-consent

curl -i -X POST \
  --url http://localhost:8001/apis/ \
  --data 'name=hydra-oauth' \
  --data upstream_url=http://"${ip}":9000/ \
  --data 'uris=/oauth2,/.well-known,/userinfo,/clients' \
  --data 'strip_uri=false' \
  --data 'preserve_host=true'

curl -i -X POST \
  --url http://localhost:8001/apis/ \
  --data 'name=login-consent' \
  --data upstream_url=http://"$ip":9001/ \
  --data 'uris=/login,/consent' \
  --data 'strip_uri=false' \
  --data 'preserve_host=true'


(cd ./hydra-login-consent-node || exit; HYDRA_URL=http://localhost:$hydraport PORT=$idport npm start &)

PORT=$hydraport \
    OAUTH2_CONSENT_URL=https://$domain/consent \
    OAUTH2_LOGIN_URL=https://$domain/login \
    OAUTH2_ISSUER_URL=https://$domain/ \
    OAUTH2_SHARE_ERROR_DEBUG=1 \
    OIDC_DYNAMIC_CLIENT_REGISTRATION_DEFAULT_SCOPE=openid,offline,offline_access,profile,email,address,phone \
    LOG_LEVEL=debug \
    DATABASE_URL=memory \
    hydra serve --dangerous-force-http

kill %1
