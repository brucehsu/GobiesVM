BINARY := bin/gobiesvm
PACKAGE := ./cmd/gobiesvm

.PHONY: all build test race vet vuln check clean

all: build

build:
	mkdir -p bin
	go build -trimpath -o $(BINARY) $(PACKAGE)

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

vuln:
	go tool govulncheck ./...

check: vet test race vuln build

clean:
	rm -f $(BINARY)
