PROJECT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
SHELL = /bin/bash
ifeq (${GOPATH},)
        GOPATH := ${HOME}/go
endif

build:
	@cd ${PROJECT_DIR} && \
	mkdir -p dist/ && \
	go build -o dist/runinit cmd/runinit.go

test:
	@cd ${PROJECT_DIR} && \
	go test ./...

.PHONY: clean
clean:
	@cd ${PROJECT_DIR} && \
	rm -rf dist/
