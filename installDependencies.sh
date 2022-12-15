#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

# Script to install all dependencies. Execute it from the top level directory
# of the Git repository.
# You already need to have 'Go' and 'npm' installed.

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

export GOBIN=$PWD/bin
go install github.com/golangci/golangci-lint/cmd/golangci-lint
go install github.com/goreleaser/goreleaser
