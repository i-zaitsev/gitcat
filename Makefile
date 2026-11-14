BINARY_NAME=gitcat
GO=go
# golangci-lint cannot type-check with a newer local Go than the one it was built
# for, so the linter runs with the toolchain that go.mod declares.
LINT_TOOLCHAIN ?= go$(shell $(GO) list -m -f '{{.GoVersion}}')

.PHONY: all build install test test-v lint fmt vet ci clean

all: lint test build

build:
	$(GO) build -o bin/$(BINARY_NAME) ./cmd/gitcat

install:
	$(GO) install ./cmd/gitcat

test:
	$(GO) test ./...

test-v:
	$(GO) test -v ./...

lint: fmt vet ci

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

ci:
	GOTOOLCHAIN=$(LINT_TOOLCHAIN) golangci-lint run

clean:
	rm -rf bin
