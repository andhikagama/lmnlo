# Lmnlo

API for User Entity

## Requirement

- [Golang](https://golang.org) - Go Programming Language (v11 and above that supports Go Modules)
- [Echo](https://echo.labstack.com/) - HTTP Framework
- [Go Modules](https://github.com/golang/go/wiki/Modules) - Go Moudules Management
- [Mockery](https://github.com/vektra/mockery) - Mock code autogenerator for golang

## Install Dependecies

Run `go mod tidy` to install dependencies required or

Run `go mod vendor` to install dependencies required using vendor directory (like dep or glide)

Make sure you have supported go version to use modules

## Development server

Run `go run main.go` for a dev server. Navigate to `http://localhost:7723/`.

## Test

Run `make test` to test only.

## Build

Open Makefile then change binary name or operating system for the binary (linux/mac/windows), by default it is compiled for linux

Run `make build` to test and build binary.

## Author

- **Andhika Gama** - _Initial work_ - [AndhikaGama](https://github.com/andhikagama)
