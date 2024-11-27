#!/bin/bash

set -o nounset
set -o errexit
set -o pipefail

# This script is used to generate and compare the Hydra DDL at different
# versions. This is useful when reviewing changes and troubleshooting
# migrations.
#
# Side effects:
# - Creates a directory at ./output/sql and stores the generated SQL
#
# Arguments:
# $1 - database type (e.g. postgres, mysql, sqlite)
# $2 - commit hash or branch name of the earlier version
# $3 - commit hash or branch name of the later version
#
# Example output: (the script prints the command instead of executing it)
# 'git diff --no-index ./output/sql/f5864b39.sqlite.dump.sql ./output/sql/master.sqlite.dump.sql'
#

# Usage
# ./scripts/db-diff.sh sqlite master 649f56cc
# ./scripts/db-diff.sh postgres HEAD~1 HEAD
if [ "$#" -ne 3 ]; then
	echo "Usage: $0 <sqlite|postgres|cockroach|mysql> <commit-ish> <commit-ish>"
	exit 1
fi

# Exports:
# - DB_DIALECT
# - DB_USER
# - DB_PASSWORD
# - DB_HOST
# - DB_PORT
# - DB_NAME
function hydra::util::parse-connection-url {
	local -r url=$1
	if [[ "${url}" =~ ^(.*)://([^:]*):?(.*)@\(?(.*):([0-9]*)\)?/([^?]*) ]]; then
		export DB_DIALECT="${BASH_REMATCH[1]}"
		export DB_USER="${BASH_REMATCH[2]}"
		export DB_PASSWORD="${BASH_REMATCH[3]}"
		export DB_HOST="${BASH_REMATCH[4]}"
		export DB_PORT="${BASH_REMATCH[5]}"
		export DB_DB="${BASH_REMATCH[6]}"
	else
		echo "Failed to parse URL"
		exit 1
	fi
}

function hydra::util::ensure-sqlite {
	if ! sqlite3 --version > /dev/null 2>&1; then
		echo 'Error: sqlite3 is not installed' >&2
		exit 1
	fi
}

function hydra::util::ensure-pg_dump {
	if ! pg_dump --version > /dev/null 2>&1; then
		echo 'Error: pg_dump is not installed' >&2
		exit 1
	fi
}

function hydra::util::ensure-mysqldump {
	if ! mysqldump --version > /dev/null 2>&1; then
		echo 'Error: mysqldump is not installed' >&2
		exit 1
	fi
}

function dump_pg {
	if test -z $TEST_DATABASE_POSTGRESQL; then
		echo 'Error: TEST_DATABASE_POSTGRESQL is not set; try running "source scripts/test-env.sh"' >&2
		exit 1
	fi

	hydra::util::ensure-pg_dump

	make test-resetdb >/dev/null 2>&1
	sleep 4
  go run . migrate sql up "$TEST_DATABASE_POSTGRESQL" --yes >&2 || true
	sleep 1
	pg_dump -s "$TEST_DATABASE_POSTGRESQL" | sed '/^--/d'
}

function dump_cockroach {
	if test -z $TEST_DATABASE_COCKROACHDB; then
		echo 'Error: TEST_DATABASE_COCKROACHDB is not set; try running "source scripts/test-env.sh"' >&2
		exit 1
	fi

	make test-resetdb >/dev/null 2>&1
	sleep 4
	go run . migrate sql up "$TEST_DATABASE_COCKROACHDB" --yes > /dev/null || true
	hydra::util::parse-connection-url "${TEST_DATABASE_COCKROACHDB}"
	docker run --rm --net=host -it cockroachdb/cockroach:latest-v24.1 dump --dump-all --dump-mode=schema --insecure --user="${DB_USER}" --host="${DB_HOST}" --port="${DB_PORT}"
}

function dump_sqlite {
	if test -z $SQLITE_PATH; then
		SQLITE_PATH="$(mktemp -d)/temp.sqlite"
	fi

	hydra::util::ensure-sqlite

	rm "$SQLITE_PATH" > /dev/null 2>&1 || true
	go run -tags sqlite,sqlite_omit_load_extension . migrate sql up "sqlite://$SQLITE_PATH?_fk=true" --yes > /dev/null 2>&1 || true
	echo '.dump' | sqlite3 "$SQLITE_PATH"
}

function dump_mysql {
	if test -z $TEST_DATABASE_MYSQL; then
		echo 'Error: TEST_DATABASE_MYSQL is not set; try running "source scripts/test-env.sh"' >&2
		exit 1
	fi

	hydra::util::ensure-mysqldump
	make test-resetdb >/dev/null 2>&1
	sleep 10
	go run . migrate sql up "$TEST_DATABASE_MYSQL" --yes > /dev/null || true
	hydra::util::parse-connection-url "${TEST_DATABASE_MYSQL}"
	mysqldump --user="$DB_USER" --password="$DB_PASSWORD" --host="$DB_HOST" --port="$DB_PORT" "$DB_DB" --no-data
}

if ! git diff-index --quiet HEAD --; then
	echo 'Error: working tree is dirty' >&2
	exit 1
fi

case $1 in
	postgres)
		DUMP_CMD=dump_pg
		;;
	cockroach)
		DUMP_CMD=dump_cockroach
		;;
	sqlite)
		DUMP_CMD=dump_sqlite
		;;
	mysql)
		DUMP_CMD=dump_mysql
		;;
	*)
		echo 'Error: unknown database type' >&2
		exit 1
		;;
esac

DIALECT=$1
COMMIT_FROM=$(git rev-parse "$2")
COMMIT_TO=$(git rev-parse "$3")
DDL_FROM="./output/sql/$COMMIT_FROM.$DIALECT.dump.sql"
DDL_TO="./output/sql/$COMMIT_TO.$DIALECT.dump.sql"

mkdir -p ./output/sql/

set -x
# shellcheck disable=SC2064

if git symbolic-ref --quiet HEAD; then
	trap "git checkout -q $(git symbolic-ref HEAD); git symbolic-ref HEAD $(git symbolic-ref HEAD)" EXIT
else
	trap "git checkout $(git rev-parse HEAD)" EXIT
fi

git checkout "$COMMIT_FROM" >/dev/null 2>&1
$DUMP_CMD > "$DDL_FROM"

git checkout "$COMMIT_TO" >/dev/null 2>&1
$DUMP_CMD > "$DDL_TO"

set +x
echo '+--------------------------'
echo '|'
echo '| Use the following command to print the diff:'
echo '| git diff --no-index '"$DDL_FROM"' '"$DDL_TO"
echo '|'
echo '+--------------------------'
set -x
