#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

# Script to update all dependencies website dependencies, in particular the
# theme. Execute it from the 'website' directory.  hugo needs to be installed.
#
# GitHub dependabot does not support hugo and they are currently not adding new
# ecosystems
# (https://github.com/dependabot/dependabot-core/blob/06702c83e5d964173504ddba6838c37fe37207ba/CONTRIBUTING.md)

hugo mod tidy
hugo mod get -u ./...
