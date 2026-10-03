.PHONY: build test vet install uninstall run clean all

BINARY_NAME=moki

# Version injected into the binary, derived from the latest git tag.
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS=-X github.com/ztkent/moki/internal/app.Version=$(VERSION)

# The Go path
GOPATH=$(shell go env GOPATH)
GOBUILD=go build
GOTEST=go test
GOVET=go vet
GOCLEAN=go clean

test:
	$(GOTEST) -race ./...

vet:
	$(GOVET) ./...

build:
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) -v ./cmd/moki

run: build
	./$(BINARY_NAME)

install: build
	mv $(BINARY_NAME) $(GOPATH)/bin

clean:
	$(GOCLEAN)
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

all: vet test install