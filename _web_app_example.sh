docker-compose exec hydra \
    /app/hydra token user \
    --client-id auth-code-client2 \
    --client-secret secret \
    --endpoint http://192.168.99.101:4444/ \
    --port 5555 \
    --scope openid,offline