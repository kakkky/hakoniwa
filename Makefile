.PHONY: wire build run mock

wire:
	go run github.com/google/wire/cmd/wire@latest ./cmd

mock:
	go generate ./domain/...

build:
	go build -o ./bin/hakoniwa ./cmd

run:
	go run ./cmd
