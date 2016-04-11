include golang.mk
.DEFAULT_GOAL := test

.PHONY: all build clean test golint README

SHELL := /bin/bash
PKG := gopkg.in/Clever/gearman.v2
PKGS := $(shell go list ./... | grep -v vendor)
EXECUTABLE := $(shell basename $(PKG))
READMES := $(addsuffix README.md, $(SUBPKGSREL))
$(eval $(call golang-version-check,1.6))

export GO15VENDOREXPERIMENT = 1

all: test build

build:
	go build -o bin/$(EXECUTABLE) $(PKG)

test: $(PKGS)
$(PKGS): golang-test-all-deps
			$(call golang-test-all,$@)

docs: $(READMES) README.md
README.md: *.go
	@go get github.com/robertkrimen/godocdown/godocdown
	godocdown $(PKG) > $@
%/README.md: PATH := $(PATH):$(GOPATH)/bin
%/README.md: %/*.go
	@go get github.com/robertkrimen/godocdown/godocdown
	godocdown $(PKG)/$(shell dirname $@) > $@
