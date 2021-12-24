GO ?= go

DATE := $(shell date -u '+%FT%T%z')
GITHUB_SHA ?= $(shell git rev-parse HEAD)
GITHUB_REF ?= local

build:
	$(GO) build -v -o bin/traceable -ldflags='-X "main.version=example" -X "main.commit=example" -X "main.date=example" -X "main.builtBy=example"' ./cmd/traceable