# Upgo

[![Go Reference](https://pkg.go.dev/badge/github.com/porjo/upgo.svg)](https://pkg.go.dev/github.com/porjo/upgo)

Upgo is a API client library for [Up Bank Australia](https://developer.up.com.au/) written in Go. Upgo provides these methods:

- `GetAccounts`
- `GetTransactions`

Upgo is a minimal wrapper around the complete API generated from the OpenAPI spec - see [`./oapi`](./oapi). Users needing more advanced functionality should use `./oapi` directly.

## Usage

See example utility [./cmd/upgo/main.go](cmd/upgo/main.go)