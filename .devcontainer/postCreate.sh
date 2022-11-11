#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

go mod download

cd frontend
npm install
cd -

cd e2e_tests
npm install
cd -

cd website
go mod download
cd -
