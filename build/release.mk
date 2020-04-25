RELEASE_SCRIPT ?= ./scripts/release.sh

GOTOOLS += github.com/goreleaser/goreleaser

REL_CMD ?= goreleaser
DIST_DIR ?= ./dist

ARCHIVE_URL       ?= https://github.com/xaque208/$(strip $(PROJECT_NAME))/archive/v$(strip $(PROJECT_VER_TAGGED)).tar.gz

# Example usage: make release version=0.11.0
release: build
	@echo "=== $(PROJECT_NAME) === [ release          ]: Generating release."
	$(RELEASE_SCRIPT) $(version)

release-clean:
	@echo "=== $(PROJECT_NAME) === [ release-clean    ]: distribution files..."
	@rm -rfv $(DIST_DIR)/*

release-publish: clean tools docker-login
	@echo "=== $(PROJECT_NAME) === [ release-publish  ]: publishing release via $(REL_CMD)"
	$(REL_CMD)

