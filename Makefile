.PHONY: build test clean

BINARY_NAME=moki

# The Go path
GOPATH=$(shell go env GOPATH)

# The build commands
GOBUILD=go build
GOTEST=go test
GOCLEAN=go clean
GOGET=go get
GOMODTIDY=go mod tidy
GOMODVENDOR=go mod vendor
GOINSTALL=go install

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/moki

install:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/moki
	mv moki $(GOPATH)/bin

all: clean deps test build