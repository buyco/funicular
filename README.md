# Funicular
[![GoDoc](https://godoc.org/github.com/buyco/funicular?status.svg)](http://godoc.org/github.com/buyco/funicular) [![Build Status](https://github.com/buyco/funicular/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/buyco/funicular/actions/workflows/test.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/buyco/funicular)](https://goreportcard.com/report/github.com/buyco/funicular) [![license](https://img.shields.io/github/license/buyco/funicular.svg)](https://github.com/buyco/funicular/LICENSE)

###### 01000110 01010101 01001110 01001001 01000011 01010101 01001100 01000001 01010010

Simple facades to create commands.

## Important information

This package use `retract` directive in `go.mod` file and is now only compatible with Go 1.16 from v1.0.0

## How to install Go ?

#### Debian / Ubuntu:
```bash
$ sudo apt update
$ sudo apt install golang
```

#### Arch:
```bash
$ sudo pacman -Sy go
```

#### Mac OS X:
```bash
$ brew update
$ brew install golang
```

#### Last release from script:
See: https://github.com/udhos/update-golang

#### From tarballs:
See: https://golang.org/doc/install


## Check golang version
```bash
$ go version
```

## Install from Makefile

### Commands:
To list available commands:
```bash
$ make help
```

To compile examples:
```bash
$ make build
```

To run tests:
```bash
$ go install github.com/joho/godotenv/cmd/godotenv
$ godotenv -f <env_file> make test
```
