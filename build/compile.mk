#
# Makefile Fragment for Compiling
#

GO         ?= go
BUILD_DIR  ?= ./bin/
PROJECT_MODULE ?= $(shell $(GO) list -m)
# $b replaced by the binary name in the compile loop, -s/w remove debug symbols
LDFLAGS    ?= "-s -w -X main.Version=$(PROJECT_VER) -X main.appName=$$b"
SRCDIR     ?= .
COMPILE_OS ?= freebsd linux

# Determine commands by looking into cmd/*
COMMANDS   ?= $(wildcard ${SRCDIR}/cmd/*)

# Determine binary names by stripping out the dir names
BINS       := $(foreach cmd,${COMMANDS},$(notdir ${cmd}))


compile-clean:
	@echo "=== $(PROJECT_NAME) === [ compile-clean    ]: removing binaries..."
	@rm -rfv $(BUILD_DIR)/*

compile: deps compile-only

proto: proto-grpc gofmt-fix

proto-grpc:
	@echo "=== $(PROJECT_NAME) === [ proto compile    ]: compiling protobufs:"
	@protoc -I ./ \
		--go_out=./ --go_opt=paths=source_relative \
		--go-grpc_out=./ --go-grpc_opt=paths=source_relative \
		rpc/rpc.proto \
		pkg/iot/iot.proto \
		internal/astro/astro.proto \
		internal/agent/agent.proto \
		modules/lights/lights.proto \
		modules/timer/named/named.proto \
		modules/inventory/inventory.proto \
		modules/telemetry/telemetry.proto
	@protoc -I modules/inventory/ -I ./ \
		--gotemplate_out=template_dir=modules/inventory/templates,debug=false,single-package-mode=true,all=true:modules/inventory \
		modules/inventory/inventory.proto
	@protoc -I modules/inventory/ -I ./ \
		--gotemplate_out=template_dir=cmd/inventory/templates,debug=false,single-package-mode=true,all=true:cmd/inventory \
		modules/inventory/inventory.proto

compile-all: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@mkdir -p $(BUILD_DIR)/$(GOOS)
	@for b in $(BINS); do \
		for os in $(COMPILE_OS); do \
			echo "=== $(PROJECT_NAME) === [ compile          ]:     $(BUILD_DIR)$$os/$$b"; \
			BUILD_FILES=`find $(SRCDIR)/cmd/$$b -type f -name "*.go"` ; \
			GOOS=$$os $(GO) build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$$os/$$b $$BUILD_FILES ; \
		done \
	done

compile-only: deps-only
	@echo "=== $(PROJECT_NAME) === [ compile          ]: building commands:"
	@mkdir -p $(BUILD_DIR)/$(GOOS)
	# CGO_ENABLED=0 GOOS=$(GOOS) $(GO) build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(GOOS)/$(PROJECT_NAME) . ; 
	@for b in $(BINS); do \
		echo "=== $(PROJECT_NAME) === [ compile          ]:     $(BUILD_DIR)$(GOOS)/$$b"; \
		BUILD_FILES=`find $(SRCDIR)/cmd/$$b -type f -name "*.go"` ; \
		GOOS=$(GOOS) $(GO) build -ldflags=$(LDFLAGS) -o $(BUILD_DIR)/$(GOOS)/$$b $$BUILD_FILES ; \
	done

# Override GOOS for these specific targets
compile-darwin: GOOS=darwin
compile-darwin: deps-only compile-only

compile-linux: GOOS=linux
compile-linux: deps-only compile-only

compile-freebsd: GOOS=freebsd
compile-freebsd: deps-only compile-only


.PHONY: clean-compile compile compile-darwin compile-linux compile-only compile-freebsd
