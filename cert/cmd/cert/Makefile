APP ?= cert

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# debug symbols:
# -gcflags "-N -l"

## default: build the application
$(APP)-$(GOOS)-$(GOARCH): $(wildcard *.go)
	goimports -l -w .
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
	-v -trimpath \
	-o $(APP)-$(GOOS)-$(GOARCH)

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


.PHONY: run
## run: run the cert command
run: $(APP)-$(GOOS)-$(GOARCH)
	./$(APP)-$(GOOS)-$(GOARCH) -v -validity 720h

.PHONY: debug-compile
## debug-compile: compile without optimisation
debug-compile:
	goimports -l -w .
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
	-gcflags "-N -l" \
	-v \
	-o $(APP)-$(GOOS)-$(GOARCH)

.PHONY: attach
## attach: debug the server
attach:
# 	dlv debug github.com/schorlet/exp/cert
# 	dlv exec ./$(APP)-$(GOOS)-$(GOARCH)
	dlv attach $(shell pgrep $(APP)-$(GOOS)-$(GOARCH))


.PHONY: test
## test: run test
test: $(APP)-$(GOOS)-$(GOARCH)
	goimports -l -w .
	go test -v ./...

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
