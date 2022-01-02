GO ?= go

build:
	$(GO) build -v -o bin/traceable ./cmd/traceable