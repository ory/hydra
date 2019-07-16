docker-compose -f quickstart.yml exec hydra \
    hydra clients create \
    --endpoint http://192.168.99.101:4445 \
    --id auth-code-client2 \
    --secret secret \
    --grant-types authorization_code,refresh_token \
    --response-types code,id_token \
    --scope openid,offline \
    --callbacks http://192.168.99.101:5555/callback