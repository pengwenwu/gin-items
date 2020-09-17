.PHONY: build help

all: help

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/main.go

help:
	@echo "make build: go build"