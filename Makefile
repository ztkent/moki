.PHONY: build test vet install uninstall run clean all

BINARY_NAME=moki

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
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/moki

run: build
	./$(BINARY_NAME)

install: build
	mv $(BINARY_NAME) $(GOPATH)/bin

clean:
	$(GOCLEAN)
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

all: vet test install