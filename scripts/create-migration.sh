#!/bin/bash

# create-migration.sh is a helper script to generate sql migration files.
# Migration file names follow the date-time naming convention established by Hydra.

set -Eeo pipefail

if [ -n "${DEBUG}" ]; then
  set -x
fi

migration_dir="$( dirname "${BASH_SOURCE[0]}" )/../persistence/sql/migrations/"
supported_dialects=(cockroach mysql postgres sqlite)

if [ ! -d "$migration_dir" ]; then
  echo "expected $migration_dir to exist"
  exit 1
fi

function usage() {
cat << EOF

Create up and down migration files for SQLite, CockroachDB, PostgreSQL
  $0 <migration-name> --all-dialects

  --all-dialects          create a file for each supported dialect (${supported_dialects[*]})

  e.g.
    $0 add_column_to_table
    $0 add_column_to_table --all-dialects
EOF
}

if [ -z "${1}" ]; then
  usage
  exit 1
fi

if [[ ! "${1}" =~ ^[A-Za-z_]+$ ]]; then
  echo "invalid migration name: ${1} must be only A-Z, a-z, or _"
  usage
  exit 1
fi

timestamp=$(date -u +%Y%m%d%H%M%S%5N)
echo $timestamp

file_prefix="${timestamp}_${1}"

function create_files() {
   dialect=$1
   filename="${file_prefix}"

   if [ ! -z "${dialect}" ]; then
     filename="${filename}.${dialect}."
   fi

   touch "${migration_dir}${filename}.up.sql"
   touch "${migration_dir}${filename}.down.sql"
}

if [ "${2}" == "--all-dialects" ]; then
  for dialect in "${supported_dialects[@]}"
  do
    create_files $dialect
  done
else
  create_files
fi
