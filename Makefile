# Makefile
include setup.mk
VERSION=$(shell git describe --tags --abbrev=0)

.PHONY: api
# generate apis by compiling proto files and generating relevant golang/python code
api:
	make -C ./apis api

.PHONY: build

.PHONY: build_dev
# build for development 
build_dev:
	make -C ./app/src/domain build
	make -C ./app/src/platform build
	make -C ./app/src/webapp build
	make -C ./app/src/vapusai build
.PHONY: release_dev
# release worker images for development testing
release_dev:
	make -C ./app/src/domain release
	make -C ./app/src/platform release
	make -C ./app/src/webapp release
	make -C ./app/src/vapusai release
# show help for managing the project in your local environment
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.PHONY: check-golangci-lint-tool
check-golangci-lint-tool:
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint is not installed. Please run \"make init\" or install the tool manually."; \
		exit 1; \
	fi

.PHONY: check-buf-tool
check-buf-tool:
	@if ! command -v buf >/dev/null 2>&1; then \
		echo "buf is not installed. Please run \"make init\" or install the tool manually."; \
		exit 1; \
	fi

.PHONY: check-goimports-tool
check-goimports-tool:
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "goimports is not installed. Please run \"make init\" or install the tool manually."; \
		exit 1; \
	fi
.PHONY: check-poetry-tool
check-poetry-tool:
	@if ! command -v poetry >/dev/null 2>&1; then \
		echo "poetry is not installed. Please install poetry (https://python-poetry.org/docs/#installation)"; \
		exit 1; \
	fi