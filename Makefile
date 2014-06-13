SHELL := /bin/bash
SUDO ?= sudo
DOCKER ?= docker
Q := github.com/modcloth/queued-command-runner
TARGETS := $(Q)/qcr

GOPATH := $(shell echo $${GOPATH%%:*})
GOBIN := $(GOPATH)/bin
PATH := $(GOBIN):$(PATH)

export GOPATH
export GOBIN
export PATH

default: test

.PHONY: all
all: clean build test

.PHONY: clean
clean:
	go clean -i -r $(TARGETS) || true
	rm -f $(GOBIN)/qcr

.PHONY: build
build: deps
	go install $(TARGETS)

.PHONY: test
test: build fmtpolice

.PHONY: godep
godep:
	go get github.com/tools/godep

.PHONY: deps
deps: godep
	@echo "godep restoring..."
	$(GOBIN)/godep restore
	go get github.com/golang/lint/golint

.PHONY: fmtpolice
fmtpolice: deps fmt lint

.PHONY: fmt
fmt:
	@echo "----------"
	@echo "checking fmt"
	@set -e ; \
	  for f in $(shell git ls-files '*.go'); do \
	  gofmt $$f | diff -u $$f - ; \
	  done

.PHONY: linter
linter:
	go get github.com/golang/lint/golint

.PHONY: lint
lint: linter
	@echo "----------"
	@echo "checking lint"
	@for file in $(shell git ls-files '*.go') ; do \
	  if [[ "$$($(GOBIN)/golint $$file)" =~ ^[[:blank:]]*$$ ]] ; then \
	  echo yayyy >/dev/null ; \
	  else $(MAKE) lintv && exit 1 ; fi \
	  done

.PHONY: lintv
lintv:
	@echo "----------"
	@for file in $(shell git ls-files '*.go') ; do $(GOBIN)/golint $$file ; done

.PHONY: save
save:
	godep save -copy=false ./...
