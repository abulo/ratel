PROJECT_NAME := "ratel"
PKG := "github.com/abulo/ratel/v3"
PKG_LIST := $(shell go list ${PKG}/... | grep -v 'history')
GO_FILES := $(shell find . -name '*.go' | grep -v 'vendor' | grep -v 'history'| grep -v _test.go)
.DEFAULT_GOAL := default
.PHONY: all test lint fmt fmtcheck cmt errcheck license

all: fmt errcheck lint build

########################################################
fmt: ## Format the files
	@gofmt -l -w $(GO_FILES)

########################################################
fmtcheck: ## Check and format the files
	@gofmt -l -s $(GO_FILES) | read; if [ $$? == 0 ]; then echo "gofmt check failed for:"; gofmt -l -s $(GO_FILES); fi

########################################################
lint:  ## lint check
	@hash revive 2>&- || go install github.com/mgechev/revive@latest
	@revive -formatter stylish ./...

########################################################
cmt: ## auto comment exported Function
	@hash gocmt 2>&- || go install github.com/Gnouc/gocmt@latest
	@gocmt -d ./... -i

########################################################
errcheck: ## check error
	@hash errcheck 2>&- || go install github.com/kisielk/errcheck@latest
	@errcheck ./...

########################################################
test: ## Run unittests
	@go test -short ${PKG_LIST}

########################################################
race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

########################################################
msan: dep ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

########################################################
dep: ## Get the dependencies
	@go get -v -d ./...

########################################################
version: ## Print git revision info
	@echo $(expr substr $(git rev-parse HEAD) 1 8)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

default: help