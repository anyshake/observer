.PHONY: build digest clean run docs gen

BINARY=observer
ifeq (${GOOS}, windows)
    BINARY := $(BINARY).exe
endif

SRC_DIR=./cmd
DIST_DIR=./build/dist
ASSETS_DIR=./build/assets
LICENSE_PATH=./LICENSE

COMMIT=$(shell git rev-parse --short HEAD)
VERSION=$(shell cat ./VERSION)

BUILD_FLAGS=-s -w -X main.version=$(VERSION) -X main.tag=$(COMMIT)-$(shell date +%s)
BUILD_ARGS=-v -trimpath

build:
	@echo "[Info] Building project, output file path: $(DIST_DIR)/$(BINARY)"
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} GOMIPS=${GOMIPS} \
		go build -ldflags="$(BUILD_FLAGS)" $(BUILD_ARGS) -o $(DIST_DIR)/$(BINARY) $(SRC_DIR)/*.go
	@cp -r $(ASSETS_DIR) $(DIST_DIR)
	@cp $(LICENSE_PATH) $(DIST_DIR)
	@echo "[Info] Build completed."

digest:
ifneq ($(wildcard $(DIST_DIR)/$(BINARY)),)
	@openssl dgst -md5 $(DIST_DIR)/$(BINARY)* | awk '{print "MD5=" $$2}'
	@openssl dgst -sha1 $(DIST_DIR)/$(BINARY)* | awk '{print "SHA1=" $$2}'
	@openssl dgst -sha256 $(DIST_DIR)/$(BINARY)* | awk '{print "SHA2-256=" $$2}'
	@openssl dgst -sha512 $(DIST_DIR)/$(BINARY)* | awk '{print "SHA2-512=" $$2}'
else
	@echo "[Error] Binary $(DIST_DIR)/$(BINARY) not found, please build first."
	@exit 1
endif

run:
	@mkdir -p $(DIST_DIR)
ifeq ($(wildcard $(DIST_DIR)/config.json.local),)
	@cp $(ASSETS_DIR)/config.json $(DIST_DIR)/config.json.local
endif
	@echo "[Info] Running project..."
	go run $(SRC_DIR)/*.go --config $(DIST_DIR)/config.json.local

clean:
	@echo "[Warn] Cleaning up project..."
	@rm -rf $(DIST_DIR)/*

docs:
ifeq ($(shell command -v swag 2> /dev/null),)
	@echo "Installing Swagger..."
	@go get github.com/swaggo/swag/cmd/swag
	@go install github.com/swaggo/swag/cmd/swag
endif
	@swag init -d ./ -o ./docs -g ./cmd/main.go

gen:
ifeq ($(shell command -v gqlgen 2> /dev/null),)
	@echo "[Info] Installing gqlgen..."
	@go get github.com/99designs/gqlgen
	@go install github.com/99designs/gqlgen
endif
	@echo "[Info] Generating GraphQL code..."
	@gqlgen generate
