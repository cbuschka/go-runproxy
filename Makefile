PROJECT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
VERSION ::= $(shell git describe --always --tags --dirty)
BUILD_TIME ::= $(shell date "+%Y-%m-%d_%H:%M:%S%:z")
COMMITISH ::= $(shell git describe --always --dirty)
OS ::= $(shell uname -s)
SHELL = /bin/bash
ifeq (${GOPATH},)
        GOPATH := ${HOME}/go
endif

define build_binary
	echo "Building $(1)/$(2)..."
	CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build \
		-a \
		-ldflags "-X github.com/cbuschka/go-runproxy/internal/build.Version=${VERSION} \
		-X github.com/cbuschka/go-runproxy/internal/build.Timestamp=${BUILD_TIME} \
		-X github.com/cbuschka/go-runproxy/internal/build.Commitish=${COMMITISH} \
		-X github.com/cbuschka/go-runproxy/internal/build.Os=$(1) \
		-X github.com/cbuschka/go-runproxy/internal/build.Arch=$(2) \
		-extldflags \"-static\"" \
		-o dist/runproxy-$(1)_$(2)$(3) \
		cmd/runproxy.go
endef

build:
	@cd ${PROJECT_DIR} && \
	mkdir -p dist/ && \
	$(call build_binary,linux,amd64,)

test:
	@cd ${PROJECT_DIR} && \
	go test ./...

run:
	@cd ${PROJECT_DIR} && \
	go run cmd/runproxy.go -v -l 0.0.0.0:8080 -d 127.0.0.1:8000 -t 5000 \
		--match-line '^Serving HTTP on.*$$' \
		-- bash -c 'python3 -u -m http.server'

.PHONY: clean
clean:
	@cd ${PROJECT_DIR} && \
	rm -rf dist/
