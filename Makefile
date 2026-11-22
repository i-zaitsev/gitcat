BINARY_NAME=gitcat
GO=go
# golangci-lint cannot type-check with a newer local Go than the one it was built for:
# the linter runs with the toolchain that go.mod declares
LINT_TOOLCHAIN ?= go$(shell $(GO) list -m -f '{{.GoVersion}}')
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS = -X github.com/i-zaitsev/gitcat/internal/version.version=$(VERSION)

.PHONY: all build install test test-v lint fmt vet ci snapshot release release-patch release-minor release-major clean

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

# builds the release archives into dist/ without publishing anything
snapshot:
	goreleaser release --snapshot --clean --skip=publish

# reminds which release targets exist
release:
	echo "usage: make release-patch | release-minor | release-major"
	exit 1

# starts the Release workflow for the pushed main and waits until it finishes
release-patch release-minor release-major: release-%:
	git fetch --quiet origin main
	[ "$$(git rev-parse HEAD)" = "$$(git rev-parse origin/main)" ] || { echo "HEAD is not origin/main: push first"; exit 1; }
	prev=$$(gh run list --workflow release.yml --limit 1 --json databaseId --jq '.[0].databaseId'); \
	gh workflow run release.yml -f bump=$*; \
	until id=$$(gh run list --workflow release.yml --limit 1 --json databaseId --jq '.[0].databaseId'); [ "$$id" != "$$prev" ]; do sleep 2; done; \
	gh run watch "$$id" --exit-status

clean:
	rm -rf bin dist
