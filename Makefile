VERSION          := $(shell git describe --tags --always --dirty="-dev")
COMMIT           := $(shell git rev-parse --short HEAD)
DATE             := $(shell date -u '+%Y-%m-%d-%H%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.version=$(VERSION)" -X "main.commit=$(COMMIT)" -X "main.date=$(DATE)"'

all: build
build: main.go
	go build -o bin/lsaddr $(VERSION_FLAGS)
install: build
	mv bin/lsaddr $(GOBIN)
test:
	go test ./...
fmt:
	gofmt -s -w .
