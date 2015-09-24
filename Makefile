
NAME := app
ARCH := amd64
VERSION := 1.0
DATE := $(shell date)
COMMIT_ID := $(shell git rev-parse --short HEAD)
SDK_INFO := $(shell go version)
LD_FLAGS := -X main.buildInfo 'Version: $(VERSION), commitID: $(COMMIT_ID), build date: $(DATE), SDK: $(SDK_INFO)'

all: clean binaries 

test:
	go test

binaries: test 
	go build -ldflags "$(LD_FLAGS)" -o $(NAME)-darwin-$(ARCH)

clean: 
	go clean