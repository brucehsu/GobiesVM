BUILDDIR          := ${CURDIR}
GOPATH            := ${BUILDDIR}

all: dependency build

build:
	cd bin ; go build -x gobiesvm

goenv:
	wget https://raw.githubusercontent.com/c9s/goenv/master/goenv

dependency:
	go get github.com/kr/try

clean:
	rm goenv
	rm bin/*
	rm -rf vendor
