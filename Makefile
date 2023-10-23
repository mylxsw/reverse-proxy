
build:
	go build -race -o build/debug/reverse-proxy main.go

build-release:
	GOOS=linux GOARCH=amd64 go build  -o build/release/reverse-proxy main.go

.PHONY: build build-release 

