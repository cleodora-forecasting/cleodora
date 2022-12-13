################### Start helper #########################
# Helper for nice Makefile documentation as seen here:
# https://github.com/audreyr/cookiecutter-pypackage
.DEFAULT_GOAL := help
define BROWSER_PYSCRIPT
import os, webbrowser, sys
try:
    from urllib import pathname2url
except:
    from urllib.request import pathname2url

webbrowser.open("file://" + pathname2url(os.path.abspath(sys.argv[1])))
endef
export BROWSER_PYSCRIPT

define PRINT_HELP_PYSCRIPT
from __future__ import print_function
import re, sys

result = []
for line in sys.stdin:
    match = re.match(r'^([a-zA-Z_-]+):.*?## (.*)$$', line)
    if match:
        target, help = match.groups()
        result.append("%-20s %s" % (target, help))
result.sort()
print(*result, sep='\n')
endef
export PRINT_HELP_PYSCRIPT
BROWSER := python3 -c "$$BROWSER_PYSCRIPT"

.PHONY: help
help:
	@python3 -c "$$PRINT_HELP_PYSCRIPT" < $(MAKEFILE_LIST)

################### End helper ###########################

SHELL=/bin/bash

EMBEDDED_FRONTEND_DIR=cleosrv/cleosrv/frontend_build

.PHONY: build
build: ## Build Cleodora binary
	@rm -rf frontend/build
	@cd frontend; npm run build
	@rm -rf $(EMBEDDED_FRONTEND_DIR)
	@cp -r frontend/build $(EMBEDDED_FRONTEND_DIR)
	@mkdir -p build
	@go build \
		-ldflags "-X github.com/cleodora-forecasting/cleodora/cleoutils.Version=`git describe --always --dirty`" \
		-tags production \
		-o build/cleosrv \
		github.com/cleodora-forecasting/cleodora/cleosrv
	@go build \
		-ldflags "-X github.com/cleodora-forecasting/cleodora/cleoutils.Version=`git describe --always --dirty`" \
		-tags production \
		-o build/cleoc \
		github.com/cleodora-forecasting/cleodora/cleoc
	@rm -rf $(EMBEDDED_FRONTEND_DIR)


.PHONY: lint
lint: ## Run linters with auto fix
	@./bin/golangci-lint run --fix
	@cd frontend && npm run lint


.PHONY: generate
generate: ## Generate code
	@go generate ./...
	@cd frontend && npm run generate
