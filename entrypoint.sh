#!/bin/sh
echo "Starting telegraf agent for statsd..."
telegraf &
sleep 3
# NOTE: Never use TEST_SAND_MIGRATE_DB in production
if [ "$TEST_SAND_MIGRATE_DB" == 'true' ]; then
    echo "migrating db..."
    /bin/sand migrate sql $DATABASE_URL
fi
echo "Starting sand..."
# NOTE: Never use TEST_SAND_DANGEROUS_FORCE_HTTP in production
if [ "$TEST_SAND_DANGEROUS_FORCE_HTTP" == 'true' ]; then
    /bin/sand host --dangerous-force-http
else
    /bin/sand host
fi
