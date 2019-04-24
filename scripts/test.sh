#!/bin/bash

echo "Initializing Docker Environment"

docker-compose -f quickstart.yml \
    -f quickstart-tracing.yml \
    -f quickstart-postgres.yml \
    -f quickstart-mysql.yml \
    kill
docker-compose -f quickstart.yml \
    -f quickstart-tracing.yml \
    -f quickstart-postgres.yml \
    -f quickstart-mysql.yml \
    rm -f
make docker
docker-compose -f quickstart.yml \
    -f quickstart-cors.yml \
    -f quickstart-debug.yml \
    -f quickstart-tracing.yml \
    -f quickstart-postgres.yml \
    up --build -d

echo "Starting Cypress"
npm run test
