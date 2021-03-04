# Funicular
[![GoDoc](https://godoc.org/github.com/buyco/funicular?status.svg)](http://godoc.org/github.com/buyco/funicular) [![Build Status](https://travis-ci.com/buyco/funicular.svg?branch=master)](https://travis-ci.com/buyco/funicular) [![Go Report Card](https://goreportcard.com/badge/github.com/buyco/funicular)](https://goreportcard.com/report/github.com/buyco/funicular) [![license](https://img.shields.io/github/license/buyco/funicular.svg)](https://github.com/buyco/funicular/LICENSE)

###### 01000110 01010101 01001110 01001001 01000011 01010101 01001100 01000001 01010010

Simple facades to create commands.

## Important information

This package use `retract` directive in `go.mod` file and is now only compatible with Go 1.16 from v1.0.0

## Run commands

```bash
$ export GO111MODULE=on # Optional from Go 1.16.x
$ go get ./...
$ cp .env-example .env
$ cd cmd/<cmd>
$ go build ./<cmd>
```

## Run tests locally

```bash
$ export GO111MODULE=on # Optional from Go 1.16.x
$ go get ./...
$ go get github.com/onsi/gomega
$ go get github.com/onsi/ginkgo/ginkgo
$ GO111MODULE=off go get github.com/joho/godotenv/cmd/godotenv
$ godotenv -f <env_file> ginkgo -r --randomizeAllSpecs --randomizeSuites --race --trace
```
