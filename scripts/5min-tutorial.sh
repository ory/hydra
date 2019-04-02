#!/bin/bash

DB=${DB:-postgres}
TRACING=${TRACING:-false}
TWOC=${TWOC:-false}

DC="docker-compose -f docker-compose.yml"
if [[ $DB == "mysql" ]]; then
    DC+=" -f docker-compose-mysql.yml"
fi
if [[ $DB == "postgres" ]]; then
    DC+=" -f docker-compose-postgres.yml"
fi
if [[ $TRACING == true ]]; then
    DC+=" -f docker-compose-tracing.yml"
fi
if [[ $TWOC == true ]]; then
    DC+=" -f docker-compose-twoc.yml"
fi
DC+=" up --build"

$DC

