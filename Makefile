VERSION=`git rev-parse --short HEAD`
#flags=-ldflags="-s -w -X main.version=${VERSION}"
OS := $(shell uname)
ifeq ($(OS),Darwin)
flags=-ldflags="-s -w -X main.version=${VERSION}"
else
flags=-ldflags="-s -w -X main.version=${VERSION} -extldflags -static"
endif

all: build

vet:
	go vet .

build:
	go clean; rm -rf pkg; CGO_ENABLED=0 go build -o jy ${flags}

build_amd64: build_linux

build_darwin:
	go clean; rm -rf pkg jy; GOOS=darwin CGO_ENABLED=0 go build -o jy ${flags}

build_linux:
	go clean; rm -rf pkg jy; GOOS=linux CGO_ENABLED=0 go build -o jy ${flags}

build_power8:
	go clean; rm -rf pkg jy; GOARCH=ppc64le GOOS=linux CGO_ENABLED=0 go build -o jy ${flags}

build_arm64:
	go clean; rm -rf pkg jy; GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -o jy ${flags}

build_windows:
	go clean; rm -rf pkg jy; GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build -o jy ${flags}

test : test1

test1:
	go test -v -bench=.
