SHELL := /bin/bash
PKG = github.com/Clever/gearman
SUBPKGSREL = $(shell ls -d */ | grep -v bin | grep -v deb)
SUBPKGS = $(addprefix $(PKG)/,$(SUBPKGSREL))
PKGS = $(PKG) $(SUBPKGS)
.PHONY: test $(PKGS) $(SUBPKGSREL)

REDIS_URL ?= localhost:6379

test: $(PKGS)

$(PKGS):
ifeq ($(LINT),1)
	golint $(GOPATH)/src/$@*/**.go
endif
	go get -d -t $@
ifeq ($(COVERAGE),1)
	REDIS_URL=$(REDIS_URL) go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	go tool cover -html=$(GOPATH)/src/$@/c.out
else
	REDIS_URL=$(REDIS_URL) go test $@ -test.v
endif

$(SUBPKGSREL): %: $(addprefix $(PKG)/, %)
