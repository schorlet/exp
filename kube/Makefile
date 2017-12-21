APP?=kube
PORT?=8000

RELEASE := 0.0.1
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%F %T %Z')

$(APP): $(wildcard *.go)
	goimports -l -w $?
	gofmt -e -d -s $?
	go tool vet -all $?
	go build \
	-ldflags "-s -w \
	-X main.Release=$(RELEASE) \
	-X main.Commit=$(COMMIT) \
	-X 'main.BuildTime=$(BUILD_TIME)'" \
	-o $(APP)

lint:
	gometalinter $(@D) --deadline=20s --aggregate --sort=linter --format="{{.Linter}}: {{.Path}}:{{.Line}}:{{if .Col}}{{.Col}}{{end}}: {{.Message}}"

clean:
	rm -f $(APP)

run: $(APP)
	PORT=$(PORT) ./$(APP)

test:
	go test -v -coverprofile=c.out
	go tool cover -func=c.out
