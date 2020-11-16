#!/bin/bash

set -euxo pipefail

go test -tags conformity -test.timeout 120m -failfast .
