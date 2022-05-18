#!/bin/bash

set -Eeuo pipefail

# This script generates down migrations templates and is part of the
# SQL migration pipeline.

for f in $(find . -name "*.up.sql"); do
	base=$(basename $f)
	dir=$(dirname $f)
	migra_name=$(echo $base | sed -e "s/\..*\.up\.sql//" | sed -e "s/\.up\.sql//")
	echo $migra_name
	if compgen -G "$dir/$migra_name*.down.sql" > /dev/null; then
		echo "Down migration exists"
	else
		echo "Down migration does not exist"
		touch $dir/$migra_name.down.sql
	fi
done
