SHELL := /bin/bash
PKG := github.com/Clever/gearman
SUBPKGS := $(addprefix $(PKG)/, $(shell ls -d */))
PKGS = $(PKG) $(SUBPKGS)

.PHONY: test golint README

test: $(PKGS)

golint:
	@go get github.com/golang/lint/golint

$(PKGS): golint README
	@go get -d -t $@
	@gofmt -w=true $(GOPATH)/src/$@/*.go
ifneq ($(NOLINT),1)
	@echo "LINTING..."
	@PATH=$(PATH):$(GOPATH)/bin golint $(GOPATH)/src/$@/*.go
	@echo ""
endif
ifeq ($(COVERAGE),1)
	@go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	@go tool cover -html=$(GOPATH)/src/$@/c.out
else
	@echo "TESTING..."
	@go test $@ -test.v
	@echo ""
endif
