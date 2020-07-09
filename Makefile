SHELL := /bin/bash
NAME := chatops-actions
GO := GO111MODULE=on GO15VENDOREXPERIMENT=1 go
GO_NOMOD := GO111MODULE=off go
PACKAGE_NAME := github.com/plumming/chatops-actions
ROOT_PACKAGE := github.com/plumming/chatops-actions
ORG := plumming

# set dev version unless VERSION is explicitly set via environment
VERSION ?= $(shell echo "$$(git describe --abbrev=0 --tags 2>/dev/null)-dev+$(REV)" | sed 's/^v//')

GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
GO_DEPENDENCIES := $(shell find . -type f -name '*.go')

REV        := $(shell git rev-parse --short HEAD 2> /dev/null || echo 'unknown')
SHA1       := $(shell git rev-parse HEAD 2> /dev/null || echo 'unknown')
BRANCH     := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown')
BUILD_DATE := $(shell date +%Y%m%d-%H:%M:%S)
BUILDFLAGS :=
CGO_ENABLED = 0
BUILDTAGS :=

GOPATH1=$(firstword $(subst :, ,$(GOPATH)))

export PATH := $(PATH):$(GOPATH1)/bin

CLIENTSET_NAME_VERSIONED := v0.15.11

build: $(GO_DEPENDENCIES)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(BUILDTAGS) $(BUILDFLAGS) main.go

all: version check

check: fmt build

version:
	echo "Go version: $(GO_VERSION)"

get-fmt-deps: ## Install goimports
	$(GO_NOMOD) get golang.org/x/tools/cmd/goimports

importfmt: get-fmt-deps
	@echo "Formatting the imports..."
	goimports -w $(GO_DEPENDENCIES)

fmt: importfmt
	@FORMATTED=`$(GO) fmt $(PACKAGE_DIRS)`
	@([[ ! -z "$(FORMATTED)" ]] && printf "Fixed unformatted files:\n$(FORMATTED)") || true

clean:
	rm -rf build release

modtidy:
	$(GO) mod tidy
