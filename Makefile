.PHONY: build run test test-integration lint clean release-snapshot release-check

BINARY := mynow
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/mynow

run: build
	./bin/$(BINARY)

test:
	go test ./... -v -count=1

test-integration:
	MYN_INTEGRATION_TEST=1 go test ./test/integration/... -v -count=1 -timeout 120s

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ dist/

release-snapshot:
	goreleaser release --snapshot --clean

release-check:
	goreleaser check
