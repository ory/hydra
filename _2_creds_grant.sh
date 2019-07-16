docker-compose -f quickstart.yml exec hydra \
    hydra token client \
    --endpoint http://192.168.99.101:4444/ \
    --client-id my-client2 \
    --client-secret secret