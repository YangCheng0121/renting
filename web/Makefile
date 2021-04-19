
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


build:

	CGO_ENABLED=0 GOOS=linux  /usr/local/go/bin/go build -a -installsuffix cgo -ldflags '-w' -i -o homeweb-web main.go

test:
	go test -v ./... -cover

docker:
	docker build . -t homeweb-web:latest
