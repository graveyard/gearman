include golang.mk
.DEFAULT_GOAL := test

.PHONY: all build clean test

SHELL := /bin/bash
PKG := gopkg.in/Clever/gearman.v2
PKGS := $(shell go list ./... | grep -v vendor)
EXECUTABLE := $(shell basename $(PKG))
$(eval $(call golang-version-check,1.9))

all: test build

build:
	go build -o bin/$(EXECUTABLE) $(PKG)

test: $(PKGS)
$(PKGS): golang-test-all-deps
	$(call golang-test-all,$@)

install_deps: golang-dep-vendor-deps
	$(call golang-dep-vendor)
