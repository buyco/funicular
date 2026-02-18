# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin

## build: Run build on examples
build: go-install go-build

## test: Run test suites
test: go-install go-test

## check: Format and lint
check: fmt lint go-tidy

## lint: Run go fmt
fmt:
	@echo "  >  Running go formatter..."
	gofmt -w -s ./ 1>&2

## lint: Run linter
lint:
	@echo "  >  Running staticcheck go linter..."
	@GOBIN=$(GOBIN) go install honnef.co/go/tools/cmd/staticcheck@v0.6.1
	@$(GOBIN)/staticcheck -checks all ./...

## lint: Run vet
vet:
	@echo "  >  Running go vet..."
	go vet ./...

mod-outdated:
	@echo "  > Looking for outdated dependencies..."
	@GOBIN=$(GOBIN) go install -mod=mod github.com/psampaz/go-mod-outdated@v0.9.0
	@go list -u -m -mod=mod -json all | $(GOBIN)/go-mod-outdated -update -direct

go-test:
	@echo "  >  Run tests..."
	@GOBIN=$(GOBIN) go install github.com/onsi/ginkgo/v2/ginkgo@v2.28.1
	@$(GOBIN)/ginkgo -r --randomize-all --randomize-suites --race --trace -coverprofile=cover.out -gcflags="-l" 1>&2

go-install:
	@echo "  >  Checking if there is any missing dependencies..."
	go mod download

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
