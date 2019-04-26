#!/bin/bash

set -euxo pipefail

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

killall hydra || true
killall node || true

# Install ORY Hydra
export GO111MODULE=on
if [[ ! -d "../../node_modules/" ]]; then
    (cd ../..; npm ci)
fi
go install github.com/ory/hydra
go build -buildmode=plugin -o ./memtest.so ./plugin/memtest.go

# Install oauth2-client
if [[ ! -d "./oauth2-client/node_modules/" ]]; then
    (cd oauth2-client; npm ci)
fi
(cd oauth2-client; ADMIN_URL=http://127.0.0.1:5001 PUBLIC_URL=http://127.0.0.1:5000 PORT=5003 npm run start >> ../oauth2-client.e2e.log 2>&1 &)

# Install consent app
(cd oauth2-client; PORT=5002 HYDRA_ADMIN_URL=http://127.0.0.1:5001 npm run consent >> ../login-consent-logout.e2e.log 2>&1 &)

source ./circle-ci.env.bash

case "$1" in
        memory)
            DSN=memory \
                hydra serve all --dangerous-force-http --disable-telemetry >> ./hydra.e2e.log 2>&1 &
            export CYPRESS_jwt_enabled=false
            ;;

        memory-jwt)
            DSN=memory \
                OAUTH2_ACCESS_TOKEN_STRATEGY=jwt \
                OIDC_SUBJECT_IDENTIFIERS_ENABLED=public \
                hydra serve all --dangerous-force-http --disable-telemetry >> ./hydra.e2e.log 2>&1 &
            export CYPRESS_jwt_enabled=true
            ;;

        postgres)
            hydra migrate sql $TEST_DATABASE_POSTGRESQL
            DSN=$TEST_DATABASE_POSTGRESQL \
                hydra serve all --dangerous-force-http --disable-telemetry >> ./hydra.e2e.log 2>&1 &
            export CYPRESS_jwt_enabled=false
            ;;

        postgres-jwt)
            hydra migrate sql $TEST_DATABASE_POSTGRESQL
            DSN=$TEST_DATABASE_POSTGRESQL \
                OAUTH2_ACCESS_TOKEN_STRATEGY=jwt \
                OIDC_SUBJECT_IDENTIFIERS_ENABLED=public \
                hydra serve all --dangerous-force-http --disable-telemetry >> ./hydra.e2e.log 2>&1 &
            export CYPRESS_jwt_enabled=true
            ;;

        mysql)
            hydra migrate sql $TEST_DATABASE_MYSQL
            DSN=$TEST_DATABASE_MYSQL \
                hydra serve all --dangerous-force-http --disable-telemetry >> ./hydra.e2e.log 2>&1 &
            export CYPRESS_jwt_enabled=false
            ;;

        mysql-jwt)
            DSN=$TEST_DATABASE_MYSQL \
                OAUTH2_ACCESS_TOKEN_STRATEGY=jwt \
                OIDC_SUBJECT_IDENTIFIERS_ENABLED=public \
                hydra serve all --dangerous-force-http --disable-telemetry >> ./hydra.e2e.log 2>&1 &
            export CYPRESS_jwt_enabled=true
            ;;

        plugin)
            DSN=plugin://./memtest.so \
                hydra serve all --dangerous-force-http --disable-telemetry >> ./hydra.e2e.log 2>&1 &
            export CYPRESS_jwt_enabled=false
            ;;

        plugin-jwt)
            DSN=plugin://./memtest.so \
                OAUTH2_ACCESS_TOKEN_STRATEGY=jwt \
                OIDC_SUBJECT_IDENTIFIERS_ENABLED=public \
                hydra serve all --dangerous-force-http --disable-telemetry >> ./hydra.e2e.log 2>&1 &
            export CYPRESS_jwt_enabled=true
            ;;
        *)
            echo $"Usage: $0 {memory|memory-jwt}"
            exit 1
esac

npm run wait-on -- -t 6000000 http-get://localhost:5000/health/ready http-get://localhost:5001/health/ready http-get://localhost:5002/ http-get://localhost:5003/oauth2/callback

(cd ../..; npm run test)

kill %1 || true # This is oauth2-client
kill %2 || true # This is the login-consent-logout
kill %3 || true # This is the hydra

cat ./oauth2-client.e2e.log
cat ./login-consent-logout.e2e.log
cat ./hydra.e2e.log

rm ./oauth2-client.e2e.log
rm ./login-consent-logout.e2e.log
rm ./hydra.e2e.log

exit 0
