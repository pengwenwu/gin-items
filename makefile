.PHONY: build build_local restart help

all: help

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service_item_go cmd/main.go

build_local:
	go build -o service_item_go cmd/main.go

restart:
	sh restart.sh

help:
	@echo "make build: go build linux"
	@echo "make build_local: go build local"
	@echo "make restart: restart http server"