.PHONY: all test

all: check test build

check: goimports govet

goimports:
	@echo checking go imports...
	@goimports -d .

govet:
	@echo checking go vet...
	@go tool vet .

test:
	@go test -v ./...

clean:
	@go clean github.com/bruceadowns/gomiruleaker

build:
	@echo build gomiruleaker
	@go build github.com/bruceadowns/gomiruleaker
