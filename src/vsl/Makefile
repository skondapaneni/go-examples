PACKAGES=$(go list ./... | grep -v vendor)
prefix=/Users/skondapaneni/GitHub/go-examples/
exec_prefix=$(prefix)
bindir=$(exec_prefix)/bin
datarootdir=$(prefix)/share
datadir=$(datarootdir)
mandir=$(datarootdir)/man

.PHONY: all deps build test gox build-all install clean

#all: build test
all: build 

deps:
	go get github.com/golang/lint/golint
	go get github.com/vishvananda/netlink
	go get github.com/stretchr/testify/assert
	go get golang.org/x/tools/cmd/cover
	go get -u golang.org/x/lint/golint
	go get

build: deps
	go build cmd/vsl/vslctl.go
	go build net/test.go 
	go build net/test1.go 

test:
	go test -coverprofile=coverage.out; go tool cover -html=coverage.out -o coverage.html
	go vet $(PACKAGES)
	golint $(PACKAGES)

gox:
	go get github.com/mitchellh/gox
	gox -build-toolchain

build-all: test
	which gox || make gox
	gox -arch="386 amd64 arm" -os="darwin linux windows" go-examples/src/vsl/cmd/vsl

install: vslctl
	mkdir -p $(bindir)
	cp vslctl $(bindir)/vslctl

clean:
	rm -f ./vslctl
