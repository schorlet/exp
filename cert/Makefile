GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# debug symbols:
# -gcflags "-N -l"

## default: build the application
$(APP)-$(GOOS)-$(GOARCH): $(wildcard *.go **/*.go)
	goimports -l -w .
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go install \
	-v -trimpath

.PHONY: clean
## clean: clean generated files
clean:
	go clean -i -x
	rm -f $(APP)-$(GOOS)-* cover.out

.PHONY: lint
## lint: lint
lint: $(APP)-$(GOOS)-$(GOARCH)
	gofmt -e -d -s $(shell find . -name '*.go')
	go vet
	golangci-lint run .


.PHONY: test
## test: run test
test: $(APP)-$(GOOS)-$(GOARCH)
	goimports -l -w .
	go test -v .

## cover.out: test coverage
cover.out: $(APP)-$(GOOS)-$(GOARCH)
	goimports -l -w .
	go test -v -coverprofile=cover.out
	go tool cover -func=cover.out

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':'
