
GOPATH:=$(shell go env GOPATH)
.PHONY: proto test docker
build:
	go build -a -installsuffix cgo -ldflags '-w' -i -o web main.go
test:
	go test -v ./... -cover
docker:
	docker build . -t web:latest
