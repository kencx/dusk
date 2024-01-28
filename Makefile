binary = dusk

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
	cd cmd && go build -v -o ${binary} .

## run: run binary
run:
	cd cmd && ./${binary}

## clean: remove binaries, dist
clean:
	if [ -f cmd/${binary} ]; then rm cmd/${binary}; fi
	go clean

## cover: get code coverage
cover:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out

## test: run tests
test:
	go test -race ./...
