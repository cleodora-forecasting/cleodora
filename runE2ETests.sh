#!/bin/bash
set -o errexit
set -o pipefail
set -o nounset

# Execute from the top level directory of the repository.
# Build the app, run it, and execute end to end tests against it.

make build
./build/cleosrv &
CLEOSRV_PID=$!

cd e2e_tests
node_modules/.bin/mocha --timeout 15000 frontPageTest.spec.js

kill "${CLEOSRV_PID}"
