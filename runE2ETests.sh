#!/bin/bash
set -o errexit
set -o pipefail
set -o nounset

CLEOSRV_PATH=./dist/cleosrv_linux_amd64_v1/cleosrv
DB_PATH=./e2e_tests.db

# Execute from the top level directory of the repository.
# Build the app, run it, and execute end to end tests against it.

trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

rm -rf "${DB_PATH}"
make clean
make build
"${CLEOSRV_PATH}" --database "${DB_PATH}" &
CLEOSRV_PID=$!

cd e2e_tests
node_modules/.bin/mocha --timeout 15000 frontPageTest.spec.js

kill "${CLEOSRV_PID}"
