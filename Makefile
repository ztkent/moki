.PHONY: build test clean

# The name of the output binary
BINARY_NAME=bashgpt

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
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/bashgpt

install:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/bashgpt
	mv bashgpt $(GOPATH)/bin

all: clean deps test build