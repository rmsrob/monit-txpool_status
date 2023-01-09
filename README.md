[![Go Reference](https://pkg.go.dev/badge/github.com/rrobrms/monit-txpool_status.svg)](https://pkg.go.dev/github.com/rrobrms/monit-txpool_status)
[![Go Report Card](https://goreportcard.com/badge/github.com/rrobrms/monit-txpool_status)](https://goreportcard.com/report/github.com/rrobrms/monit-txpool_status)
[![Coverage Status](https://coveralls.io/repos/github/rrobrms/monit-txpool_status/badge.svg?branch=master)](https://coveralls.io/github/rrobrms/monit-txpool_status?branch=master)

# Monit txpool_status

> Get formatted response for RPC call txpool_status

## Prerequisite
Required:
- [Go](https://go.dev/doc/install)

## Install

```sh
go instal github.com/rrobrms/monit-txpool_status@latest
```

## Usage

> Adjust this variables to your needs in `main.go`
```go
var (
	RPC = "ws://127.0.0.1:8545"
	TICK = 2
)
```

```sh
go run main.go
```
