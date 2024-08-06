#!/bin/bash

DB=${DB:-postgres}
TRACING=${TRACING:-false}
PROMETHEUS=${PROMETHEUS:-false}

DC="docker compose -f quickstart.yml"
if [[ $DB == "mysql" ]]; then
    DC+=" -f quickstart-mysql.yml"
fi
if [[ $DB == "postgres" ]]; then
    DC+=" -f quickstart-postgres.yml"
fi
if [[ $TRACING == true ]]; then
    DC+=" -f quickstart-tracing.yml"
fi
if [[ $PROMETHEUS == true ]]; then
    DC+=" -f quickstart-prometheus.yml"
fi
DC+=" up --build"

$DC
