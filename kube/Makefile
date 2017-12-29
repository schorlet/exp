APP ?= kube
PORT ?= 8000

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

RELEASE := 0.0.1
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%F %T %Z')

$(APP): $(wildcard *.go)
	goimports -l -w $?
	gofmt -e -d -s $?
	go tool vet -all $?
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
	-ldflags "-s -w \
	-X main.Release=$(RELEASE) \
	-X main.Commit=$(COMMIT) \
	-X 'main.BuildTime=$(BUILD_TIME)'" \
	-o $(APP)

clean:
	rm -f $(APP)

run: $(APP)
	PORT=$(PORT) ./$(APP)

docker: $(APP)
	docker build -t $(APP):$(RELEASE) --rm=true .
	docker stop $(APP) || true
	docker rm $(APP) || true
	docker run --rm=true --name $(APP) \
	-p "127.0.0.1:$(PORT):$(PORT)" \
	-e "PORT=$(PORT)" \
	$(APP):$(RELEASE)

lint:
	gometalinter $(@D) --deadline=20s --aggregate --sort=linter --format="{{.Linter}}: {{.Path}}:{{.Line}}:{{if .Col}}{{.Col}}{{end}}: {{.Message}}"

test:
	go test -v -coverprofile=c.out
	go tool cover -func=c.out
