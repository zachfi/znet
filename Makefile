PROJECT_NAME := $(shell basename $(shell pwd))
PROJECT_VER  := $(shell git describe --tags --always --dirty)
GO_PKGS      := $(shell go list ./... | grep -v -e "/vendor/" -e "/example")
GO_FILES     := $(shell find cmd pkg -type f -name "*.go")
NATIVEOS     := $(shell go version | awk -F '[ /]' '{print $$4}')
NATIVEARCH   := $(shell go version | awk -F '[ /]' '{print $$5}')
SRCDIR       ?= .
BUILD_DIR    := ./bin/
COVERAGE_DIR := ./coverage/
GOTOOLS       = github.com/axw/gocov/gocov \
                github.com/AlekSi/gocov-xml \
                github.com/stretchr/testify/assert \
                github.com/robertkrimen/godocdown/godocdown \
                github.com/golangci/golangci-lint/cmd/golangci-lint \


GO            = go
GODOC         = godocdown
GOLINTER      = golangci-lint
GOVENDOR      = dep

# Determine packages by looking into pkg/*
PACKAGES=$(wildcard ${SRCDIR}/pkg/*)

# Determine commands by looking into cmd/*
COMMANDS=$(wildcard ${SRCDIR}/cmd/*)

# Determine binary names by stripping out the dir names
#BINS=$(foreach cmd,${COMMANDS},$(notdir ${cmd}))
BINS=znet

LDFLAGS='-X main.Version=$(PROJECT_VER)'

all: build

# Humans running make:
build: check-version clean validate test coverage compile

# Build command for CI tooling
build-ci: check-version clean validate test compile-only

clean:
	@echo "=== $(PROJECT_NAME) === [ clean            ]: removing binaries and coverage file..."
	@rm -rfv $(BUILD_DIR)/* $(COVERAGE_DIR)/*

tools: check-version
	@echo "=== $(PROJECT_NAME) === [ tools            ]: Installing tools required by the project..."
	@$(GO) get $(GOTOOLS)

tools-update: check-version
	@echo "=== $(PROJECT_NAME) === [ tools-update     ]: Updating tools required by the project..."
	@$(GO) get -u $(GOTOOLS)

deps: tools deps-only

deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Installing package dependencies required by the project..."
	@$(GOVENDOR) ensure

validate: deps
	@echo "=== $(PROJECT_NAME) === [ validate         ]: Validating source code running $(GOLINTER)..."
	@$(GOLINTER) run ./...

compile-only:
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	$(GO) build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME) .; \
	# @for b in $(BINS); do \
	# 	echo "=== $(PROJECT_NAME) === [ compile          ]:     $$b"; \
	# done

compile: deps compile-only

coverage:
	@echo "=== $(PROJECT_NAME) === [ coverage         ]: generating coverage results..."
	@rm -rf $(COVERAGE_DIR)/*
	@for d in $(GO_PKGS); do \
		pkg=`basename $$d` ;\
		$(GO) test -tags 'unit integration' -coverprofile $(COVERAGE_DIR)/$$pkg.tmp $$d ;\
	done
	@echo 'mode: set' > $(COVERAGE_DIR)/coverage.out
	# || true to ignore grep return code if no matches (i.e. no tests written...)
	@cat $(COVERAGE_DIR)/*.tmp | grep -v 'mode: set' >> $(COVERAGE_DIR)/coverage.out || true
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

test-unit:
	@echo "=== $(PROJECT_NAME) === [ unit-test        ]: running unit tests..."
	@$(GO) test -tags unit $(GO_PKGS)

test-integration:
	@echo "=== $(PROJECT_NAME) === [ integration-test ]: running integrtation tests..."
	@$(GO) test -tags integration $(GO_PKGS)

document:
	@echo "=== $(PROJECT_NAME) === [ documentation    ]: Generating Godoc in Markdown..."
	@for p in $(PACKAGES); do \
		echo "=== $(PROJECT_NAME) === [ documentation    ]:     $$p"; \
		$(GODOC) $$p > $$p/README.md ; \
	done
	@for c in $(COMMANDS); do \
		echo "=== $(PROJECT_NAME) === [ documentation    ]:     $$c"; \
		$(GODOC) $$c > $$c/README.md ; \
	done

test-only: test-unit test-integration
test: test-deps test-only

check-version:
ifdef GOOS
ifneq "$(GOOS)" "$(NATIVEOS)"
	$(error GOOS is not $(NATIVEOS). Cross-compiling is only allowed for 'clean', 'deps-only' and 'compile-only' targets)
endif
endif
ifdef GOARCH
ifneq "$(GOARCH)" "$(NATIVEARCH)"
	$(error GOARCH variable is not $(NATIVEARCH). Cross-compiling is only allowed for 'clean', 'deps-only' and 'compile-only' targets)
endif
endif

.PHONY: all build clean coverage document document-only document-deps fmt lint vet validate-deps validate-only validate compile-deps compile-only compile test-deps test-unit test-integration test-only test
