.PHONY: build test install uninstall all

BINARY_NAME=moki

# The Go path
GOPATH=$(shell go env GOPATH)
GOBUILD=go build
GOTEST=go test
GOCLEAN=go clean

test:
	$(GOTEST) -v ./...

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/moki

install: build
	mv moki $(GOPATH)/bin

clean:
	$(GOCLEAN)
	rm $(GOPATH)/bin/$(BINARY_NAME)

all: test install