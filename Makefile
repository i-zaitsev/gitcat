BINARY_NAME=gitcat
GO=go
# golangci-lint cannot type-check with a newer local Go than the one it was built
# for, so the linter runs with the toolchain that go.mod declares.
LINT_TOOLCHAIN ?= go$(shell $(GO) list -m -f '{{.GoVersion}}')
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS = -X github.com/i-zaitsev/gitcat/internal/version.version=$(VERSION)

.PHONY: all build install test test-v lint fmt vet ci snapshot clean

all: lint test build

build:
	$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/gitcat

install:
	$(GO) install -ldflags "$(LDFLAGS)" ./cmd/gitcat

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

# snapshot builds the release archives into dist/ without publishing anything.
snapshot:
	goreleaser release --snapshot --clean --skip=publish

clean:
	rm -rf bin dist
