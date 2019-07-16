docker-compose -f quickstart.yml exec hydra \
    hydra clients create \
    --endpoint http://192.168.99.101:4445/ \
    --id my-client2 \
    --secret secret \
    -g client_credentials