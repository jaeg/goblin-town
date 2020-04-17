.DELETE_ON_ERROR:

GOVERSION=1.14

compile=CGO_ENABLED=1 GO111MODULE=auto CC=$(1) GOOS=$(2) GOARCH=$(3) go build -tags static -ldflags "-s -w" -o=$(4) *.go
gofiles=$(wildcard *.go)
pkgdir=pkg
name=goblin-town
binlin64=$(pkgdir)/$(name)-linux-amd64
darwin64=$(pkgdir)/$(name)-darwin-amd64

uid=$(shell id -u)
gid=$(shell id -g)
pwd=$(shell pwd)

.PHONY: all
all: clean lint build test image

###############
# build targets

$(pkgdir):
	mkdir -p $@

.PHONY: clean
clean:
	rm -rf $(pkgdir)

.PHONY: lint
lint: GOPATH=$(shell go env GOPATH)
lint: $(GOPATH)/bin/golint
	$(GOPATH)/bin/golint ./... | grep -v '^vendor\/' | sed 's/^/golint: /' || true
	go fmt ./...
	go vet ./... 2>&1 | grep -v '^vendor\/' | grep -v '^exit\ status\ 1' | sed 's/^/go vet: /' || true

.PHONY: build
build: lin64

.PHONY: lin64
lin64: $(binlin64)
$(binlin64): $(gofiles)
	$(call compile,gcc,linux,amd64,$@)

.PHONY: darwin64
darwin64: $(darwin64)
$(darwin64): $(gofiles)
	$(call compile,gcc,darwin,amd64,$@)

$(GOPATH)/bin/golint:
	go get -u golang.org/x/lint/golint
