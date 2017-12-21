APP?=kube
PORT?=8000

RELEASE=0.0.1
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%F %T %Z')

$(APP): $(wildcard *.go)
	go build \
	-ldflags "-s -w \
	-X main.Release=$(RELEASE) \
	-X main.Commit=$(COMMIT) \
	-X 'main.BuildTime=$(BUILD_TIME)'" \
	-o $(APP)

clean:
	rm -f $(APP)

build: $(APP)

run: build
	PORT=$(PORT) ./$(APP)

test:
	go test -v
