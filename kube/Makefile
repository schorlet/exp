APP?=kube
PORT?=8000


$(APP): $(wildcard *.go)
	go build -o $(APP)

clean:
	rm -f $(APP)

build: $(APP)

run: build
	PORT=$(PORT) ./$(APP)

test:
	go test -v
