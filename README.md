# Twelve-Factor Golang Template
An example HTTP server following the [Twelve-Factor App](https://12factor.net/).

## Course of action
- `go mod init github.com/psokdhio/12factor`.
- `cobra-cli init --viper`.
- `cobra-cli add serve`.
- Add OpenAPI3 spec to `/api`. Use ping example from [github.com/oapi-codegen/oapi-codegen](https://github.com/oapi-codegen/oapi-codegen/).
- Define the `/tool.go` and `/api/generate.go`.
- `go get` and `go generate ./...`.
- Implement the server.
- Wire up the server the `serve` command.

## How to Build
Considering `go` is installed
```sh
make
```

Look inside the Makefile.

