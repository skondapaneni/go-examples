PACKAGES=$(go list ./... | grep -v vendor)
prefix=/Users/sriharikondapaneni/go-examples
exec_prefix=$(prefix)
bindir=$(exec_prefix)/bin
datarootdir=$(prefix)/share
datadir=$(datarootdir)
mandir=$(datarootdir)/man

export GOPATH = ${prefix}

GOCMD=go
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

.PHONY: all deps build 

.DEFAULT_GOAL := all 

all: deps build 

deps:
#	$(GOGET) -v ./...
	go get github.com/codegangsta/cli
	go get golang.org/x/lint/golint
	go get github.com/stretchr/testify/assert
	go get golang.org/x/tools/cmd/cover

build: deps
	go install collections
	go install libinfra
	go install node
	go install parser
	go install cli
	go install vsl
	go install main

test:
#	$(GOTEST) -v ./...
	$(GOTEST) testing

clean: 
	$(GOCLEAN)
	rm -f $(bindir)/*

