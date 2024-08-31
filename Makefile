.PHONY: build clean gen docs version run windows

BINARY=observer
VERSION=$(shell cat ./VERSION)
RELEASE=$(shell date +%Y%m%d%H%M%S)
COMMIT=$(shell git rev-parse --short HEAD)

SRC_DIR=./cmd
DIST_DIR=./build/dist
ASSETS_DIR=./build/assets

BUILD_ARCH=arm arm64 386 amd64 ppc64le riscv64 \
	mips mips64le mipsle loong64 s390x
BUILD_FLAGS=-s -w -X main.version=$(VERSION) \
	-X main.release=$(COMMIT)-$(RELEASE)
BUILD_ARGS=-trimpath

build: clean windows $(BUILD_ARCH)
$(BUILD_ARCH):
	@echo "Building Linux $@ ..."
	@mkdir -p $(DIST_DIR)/$@
	@rm -rf $(DIST_DIR)/$@/*
	@CGO_ENABLED=0 GOOS=linux GOARCH=$@ go build -ldflags="$(BUILD_FLAGS)" \
		$(BUILD_ARGS) -o $(DIST_DIR)/$@/$(BINARY) $(SRC_DIR)/*.go
	@cp -r $(ASSETS_DIR) $(DIST_DIR)/$@

windows:
	@echo "Building Windows 32-bit, 64-bit, arm, arm64 ..."
	@mkdir -p $(DIST_DIR)/win32 $(DIST_DIR)/win64 $(DIST_DIR)/winarm $(DIST_DIR)/winarm64
	@rm -rf $(DIST_DIR)/win32/* $(DIST_DIR)/win64/* $(DIST_DIR)/winarm/* $(DIST_DIR)/winarm64/*
	@CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="$(BUILD_FLAGS)" \
		$(BUILD_ARGS) -o $(DIST_DIR)/win32/$(BINARY).exe $(SRC_DIR)/*.go
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(BUILD_FLAGS)" \
		$(BUILD_ARGS) -o $(DIST_DIR)/win64/$(BINARY).exe $(SRC_DIR)/*.go
	@CGO_ENABLED=0 GOOS=windows GOARCH=arm go build -ldflags="$(BUILD_FLAGS)" \
		$(BUILD_ARGS) -o $(DIST_DIR)/winarm/$(BINARY).exe $(SRC_DIR)/*.go
	@CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="$(BUILD_FLAGS)" \
		$(BUILD_ARGS) -o $(DIST_DIR)/winarm64/$(BINARY).exe $(SRC_DIR)/*.go
	@cp -r $(ASSETS_DIR) $(DIST_DIR)/win32
	@cp -r $(ASSETS_DIR) $(DIST_DIR)/win64
	@cp -r $(ASSETS_DIR) $(DIST_DIR)/winarm
	@cp -r $(ASSETS_DIR) $(DIST_DIR)/winarm64

gen:
ifeq ($(shell command -v gqlgen 2> /dev/null),)
	@echo "Installing gqlgen..."
	@go get github.com/99designs/gqlgen
	@go install github.com/99designs/gqlgen
endif
	@gqlgen generate

version:
	@go run $(SRC_DIR)/*.go --version

run:
	@go run $(SRC_DIR)/*.go --config $(ASSETS_DIR)/config.json

clean:
	@rm -rf $(DIST_DIR)/*
