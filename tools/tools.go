//go:build tools

package tools

import (
	// build/test.mk
	_ "github.com/stretchr/testify/assert"

	// build/lint.mk
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/psampaz/go-mod-outdated"
	_ "golang.org/x/tools/cmd/goimports"

	// build/document.mk
	_ "github.com/git-chglog/git-chglog/cmd/git-chglog"
	_ "golang.org/x/tools/cmd/godoc"

	// build/test.mk
	_ "gotest.tools/gotestsum"

	// build/release.mk
	_ "github.com/goreleaser/goreleaser"

	// build/compile.mk
	_ "github.com/Masterminds/sprig/v3"
	_ "moul.io/protoc-gen-gotemplate"
)
