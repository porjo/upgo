# Upgo

## Usage

See example utility [./cmd/upgo/main.go](cmd/upgo/main.go)

## OpenAPI

Generate library from the [Up Bank OpenAPI spec](https://github.com/up-banking/api)

```
$ oapi-codegen -generate=models,client -package upgo oapi/openapi.json > oapi/oapi.go
```
