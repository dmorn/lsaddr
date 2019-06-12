all: build

build: main.go
	go build -o bin/lsaddr

test:
	go test ./...
