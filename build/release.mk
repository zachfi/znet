GOTOOLS += github.com/goreleaser/goreleaser

REL_CMD ?= goreleaser
DIST_DIR ?= ./dist

ARCHIVE_URL       ?= https://github.com/xaque208/$(strip $(PROJECT_NAME))/archive/v$(strip $(PROJECT_VER_TAGGED)).tar.gz

# Example usage: make release version=0.11.0
release: build release-publish

release-clean:
	@echo "=== $(PROJECT_NAME) === [ release-clean    ]: distribution files..."
	@rm -rfv $(DIST_DIR)/*

release-publish: clean tools
	@echo "=== $(PROJECT_NAME) === [ release-publish  ]: publishing release via $(REL_CMD)"
	$(REL_CMD)

