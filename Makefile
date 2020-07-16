SHELL := /bin/bash
NAME := chatops
GO := GO111MODULE=on GO15VENDOREXPERIMENT=1 go
GO_NOMOD := GO111MODULE=off go
PACKAGE_NAME := github.com/cagiti/chatops
ROOT_PACKAGE := github.com/cagiti/chatops
ORG := cagiti

GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/ | grep -v e2e)
GO_DEPENDENCIES := $(shell find . -type f -name '*.go')

CGO_ENABLED = 0
BUILDFLAGS :=
BUILDTAGS :=

GOPATH1=$(firstword $(subst :, ,$(GOPATH)))

export PATH := $(PATH):$(GOPATH1)/bin

CLIENTSET_NAME_VERSIONED := v0.15.11

build: $(GO_DEPENDENCIES)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(BUILDTAGS) $(BUILDFLAGS) -o build/$(NAME) main.go

docker:
	docker build -t $(ORG)/$(NAME) .

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
	rm -rf build

modtidy:
	$(GO) mod tidy

test:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test -coverprofile=coverage.out $(PACKAGE_DIRS)
