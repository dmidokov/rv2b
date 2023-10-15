include .env
export
all: build-and-run
build-and-run:
	/home/noname/go/go1.21.1/bin/go run ./main.go
