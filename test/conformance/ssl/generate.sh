#!/bin/bash

set -euxo pipefail
cd "$( dirname "${BASH_SOURCE[0]}" )"

subj="/C=GB/ST=London/L=London/O=Global Security/OU=IT Department/CN=ory.sh.local"
openssl genrsa -out ory-ca.key 2048
openssl req -x509 -new -nodes -key ory-ca.key -sha256 -days 4096 -out ory-ca.pem -subj "$subj"

NAME=ory-conformity
openssl genrsa -out $NAME.key 2048
openssl req -new -key $NAME.key -out $NAME.csr  -subj "$subj"

cat > $NAME.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = httpd
DNS.2 = hydra
DNS.3 = consent
IP.1 = 127.0.0.1
EOF

openssl x509 -req -in $NAME.csr -CA ory-ca.pem -CAkey ory-ca.key -CAcreateserial \
  -out $NAME.crt -days 825 -sha256 -extfile $NAME.ext
