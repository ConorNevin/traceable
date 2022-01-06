GO ?= go

build:
	$(GO) build -trimpath -v -o bin/traceable ./cmd/traceable