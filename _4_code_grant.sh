docker-compose -f quickstart.yml exec hydra \
    hydra clients create \
    --endpoint http://192.168.99.101:4445 \
    --id auth-code-client \
    --secret 7m91vTrREXurMca2xYgOl2uHOzVDJTrATvOR13FOxcM.wb-_gHk0pwalx3xcgkXHTqC1Y3sYSbWO7U5a1IhOEK0 \
    --grant-types authorization_code,refresh_token \
    --response-types code,id_token \
    --scope openid,offline \
    --callbacks http://192.168.99.101:5555/callback

