#!/bin/bash
set -e

# Add hydra-host as command if needed
if [[ "$1" == -* ]]; then
	set -- app "$@"
fi

# Rename hydra url to db url
if [ -z "$DATABASE_URL" ]; then
	export DATABASE_URL=$HYDRA_DATABASE_URL
fi

# Auto generate database url by linked containers
if [ -z "$DATABASE_URL" ]; then
	if [ -n "$DB_PORT_5432_TCP_PORT" ]; then
		export DATABASE_URL="postgres://postgres:$HYDRA_DATABASE_PASSWORD@db:5432/postgres?sslmode=disable"
	elif [ -n "$DB_PORT_29015_TCP_PORT"]; then
		export DATABASE_URL="rethinkdb://db:29015"
	fi
fi

# Start host
exec "$@"
