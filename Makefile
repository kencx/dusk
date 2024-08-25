binary = dusk
version = $(shell git describe --always --dirty --tags)
ldflags = -ldflags "-s -w -X main.version=${version}"

.PHONY: help air build run clean cover test

default: help

help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## air: build and run with air
air:
	air -c air.toml

## build: build binary
build:
	templ generate -path ui
	cd cmd && CGO_ENABLED=1 go build ${ldflags} -tags "linux libsqlite3 fts5" -o ${binary} .

## run: run binary
run:
	cd cmd && ./${binary}

## clean: remove binaries, dist
clean:
	if [ -f cmd/${binary} ]; then rm cmd/${binary}; fi
	if [ -d cmd/lib ]; then rm -r cmd/lib; fi
	go clean

## cover: get code coverage
cover:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out

## test: run tests
test:
	go test -race ./...
