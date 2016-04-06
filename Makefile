
NAME := s3rw
ARCH := amd64
VERSION := 1.0
DATE := $(shell date)
COMMIT_ID := $(shell git rev-parse --short HEAD)
SDK_INFO := $(shell go version)
LD_FLAGS := '-X "main.buildInfo=Version: $(VERSION), commitID: $(COMMIT_ID), build date: $(DATE), SDK: $(SDK_INFO)"'

all: binaries 

gb: 
	@which gb > /dev/null || go get github.com/constabulary/gb/...

binaries:  gb
	gb build -ldflags $(LD_FLAGS) 

clean: 
	rm -rf bin/*

phony: .all .binaries clean gb
