#!/bin/sh

/usr/bin/assert_openssl_fips.sh

# NOTE: Never use TEST_SAND_MIGRATE_DB in production
if [ "$TEST_SAND_MIGRATE_DB" == 'true' ]; then
    echo "migrating db..."
    LD_PRELOAD=/usr/lib/libfipsify.so PROCESS_NAME="sand migrate" /usr/bin/sand migrate sql $DATABASE_URL
fi
echo "Starting sand..."
# NOTE: Never use TEST_SAND_DANGEROUS_FORCE_HTTP in production
if [ "$TEST_SAND_DANGEROUS_FORCE_HTTP" == 'true' ]; then
    LD_PRELOAD=/usr/lib/libfipsify.so PROCESS_NAME="sand server" /usr/bin/sand host --dangerous-force-http
else
    LD_PRELOAD=/usr/lib/libfipsify.so PROCESS_NAME="sand server" /usr/bin/sand host
fi
