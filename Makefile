GOLANGCI_LINT_VERSION := v2.14.0

.PHONY: run test lint check

run:
	go run ./cmd/api

test:
	go test ./...

# go run builds the pinned version on first use and caches it, so local and CI use the same linter.
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

check: lint test
	go build ./...
