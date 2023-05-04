
.PHONY: build

VERSION=1.0.5d
ENVIRONMENT=develop
BINARY=observer
RELEASE=`date +%Y%m%d%H%M%S`
COMMIT=`git rev-parse --short HEAD`

build: armv7l aarch64 x86_64 i386 mipsle

armv7l:
	@mkdir -p ./release/armv7l
	@rm -rf  ./release/armv7l/*
	@env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w -X main.version=v${VERSION} -X main.release=#0~${ENVIRONMENT}${RELEASE}-${COMMIT}" -o ./release/armv7l/${BINARY} ../cmd/*.go
	@cp ./config.json ./release/armv7l
	@cp ./g-observer.service ./release/armv7l

aarch64:
	@mkdir -p ./release/aarch64
	@rm -rf  ./release/aarch64/*
	@env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.version=v${VERSION} -X main.release=#0~${ENVIRONMENT}${RELEASE}-${COMMIT}" -o ./release/aarch64/${BINARY} ../cmd/*.go
	@cp ./config.json ./release/aarch64
	@cp ./g-observer.service ./release/aarch64

x86_64:
	@mkdir -p ./release/x86_64
	@rm -rf  ./release/x86_64/*
	@env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=v${VERSION} -X main.release=#0~${ENVIRONMENT}${RELEASE}-${COMMIT}" -o ./release/x86_64/${BINARY} ../cmd/*.go
	@cp ./config.json ./release/x86_64
	@cp ./g-observer.service ./release/x86_64

i386:
	@mkdir -p ./release/i386
	@rm -rf  ./release/i386/*
	@env GOOS=linux GOARCH=386 go build -ldflags="-s -w -X main.version=v${VERSION} -X main.release=#0~${ENVIRONMENT}${RELEASE}-${COMMIT}" -o ./release/i386/${BINARY} ../cmd/*.go
	@cp ./config.json ./release/i386
	@cp ./g-observer.service ./release/i386

mipsle:
	@mkdir -p ./release/mipsle
	@rm -rf  ./release/mipsle/*
	@env GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags="-s -w -X main.version=v${VERSION} -X main.release=#0~${ENVIRONMENT}${RELEASE}-${COMMIT}" -o ./release/mipsle/${BINARY} ../cmd/*.go
	@cp ./config.json ./release/mipsle
	@cp ./g-observer.service ./release/mipsle
