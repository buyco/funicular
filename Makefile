# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin

## build: Run build on examples
build: go-get go-build

## test: Run test suites
test: go-get go-test

## check: Format and lint
check: fmt lint go-tidy

## lint: Run go fmt
fmt:
	@echo "  >  Running go formatter..."
	gofmt -w -s ./ 1>&2

## lint: Run linter
lint:
	@echo "  >  Running staticcheck go linter..."
	@GOBIN=$(GOBIN) go install honnef.co/go/tools/cmd/staticcheck@latest
	@$(GOBIN)/staticcheck -checks all ./...

## lint: Run vet
vet:
	@echo "  >  Running go vet..."
	go vet ./...

go-test:
	@echo "  >  Run tests..."
	@GOBIN=$(GOBIN) go install github.com/onsi/ginkgo/ginkgo@v1.16.4
	@$(GOBIN)/ginkgo -r --randomizeAllSpecs --randomizeSuites --race --trace --cover -gcflags="-l" 1>&2

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	go get ./...

go-tidy:
	@echo "  > Running go mod tidy"
	go mod tidy

go-build:
	@echo "  >  Building examples binaries..."
	@GOBIN=$(GOBIN) go build -tags debug $(LDFLAGS) examples/s3_app/s3_app.go
	@GOBIN=$(GOBIN) go build -tags debug $(LDFLAGS) examples/sftp_app/sftp_app.go

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo