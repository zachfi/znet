#
# Makefile fragment for installing deps
#

GO           ?= go

deps: tools deps-only

deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Installing package dependencies required by the project..."
	${GO} mod tidy
	${GO} mod vendor

.PHONY: deps deps-only
