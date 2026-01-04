REPOSITORY=$(shell git remote get-url origin | awk -F':' '{print $$2}' | awk -F'.git' '{print $$1}')
PROJECT=$(notdir $(REPOSITORY))
VERSION=$(shell git describe --tags 2>/dev/null)
BUILD=$(shell git rev-parse --short HEAD)

rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))

# REGISTRY:=registry.gitlab.com
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
CMDS=$(dir $(call rwildcard,.,*/main.go))
BINS=${subst cmd/,bin/,${realpath $*}}


# Settable
GOOS ?= linux
GOARCH ?= amd64
GO_MODULE ?= $(shell git remote get-url origin | sed 's/https:\/\///' | sed 's/.git//')
GO_BUILD_PKGS ?= ./cmd/...
GO_BUILD_OUT ?= ./bin


test:
	go test -v ./... 

clean:
	@rm -rf ${GO_BUILD_OUT}

build: clean 
	mkdir -p $(GO_BUILD_OUT)
	go version
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -installsuffix cgo $(LDFLAGS) -o $(GO_BUILD_OUT) $(GO_BUILD_PKGS)

build-local: clean
	mkdir -p $(GO_BUILD_OUT)
	go version
	CGO_ENABLED=0 go build -a -installsuffix cgo $(LDFLAGS) -o $(GO_BUILD_OUT) $(GO_BUILD_PKGS)

