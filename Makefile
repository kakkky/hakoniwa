.PHONY: wire build run

wire:
	go run github.com/google/wire/cmd/wire@latest ./cmd

build:
	go build -o ./bin/hakoniwa ./cmd

run:
	go run ./cmd
