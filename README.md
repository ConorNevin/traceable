# traceable 

[![CircleCI](https://circleci.com/gh/ConorNevin/traceable/tree/main.svg?style=svg)](https://circleci.com/gh/ConorNevin/traceable/tree/main)
[![Coverage Status](https://coveralls.io/repos/github/ConorNevin/traceable/badge.svg?branch=main)](https://coveralls.io/github/ConorNevin/traceable?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/ConorNevin/traceable)](https://goreportcard.com/report/github.com/ConorNevin/traceable)
[![GoDoc](https://godoc.org/github.com/ConorNevin/traceable?status.svg)](https://godoc.org/github.com/ConorNevin/traceable)

A Tool that generates an instrumented implementation of an interface that wraps functions calls with an OpenTracing span.

## Installation

`traceable` requires a working Go installation (Go 1.14+)
```bash
go install github.com/ConorNevin/traceable@latest
```

## Usage

### Using Go Generate

1. Add a go:generate directive to a file in the same package as the target interface: `go:generate traceable -types IFACE -output traced/iface.go`
2. Run go generate on the directory

### Download binary from GitHub release

```bash
    curl -fsSL "https://github.com/ConorNevin/traceable/releases/download/$(VERSION)/traceable_$(uname -s)_$(uname -m)" -o traceable
```