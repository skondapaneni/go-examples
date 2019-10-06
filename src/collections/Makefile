PACKAGES=$(go list ./... | grep -v vendor)
prefix=/Users/skondapaneni/go-client
exec_prefix=$(prefix)
bindir=$(exec_prefix)/bin
datarootdir=$(prefix)/share
datadir=$(datarootdir)
mandir=$(datarootdir)/man

.PHONY: all deps build 

all: build 

deps:
	go get golang.org/x/lint/golint
	go get github.com/stretchr/testify/assert
	go get golang.org/x/tools/cmd/cover
	go get
#	go get github.com/golang/lint/golint

build: deps
	go install collections

