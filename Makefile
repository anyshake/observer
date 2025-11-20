VERSION_MAJOR=4
VERSION_MINOR=2
VERSION_PATCH=1

# Caution: software versioning mechanism depends on format of first three lines in this file

.PHONY: build digest clean run gen sync

GO ?= go

ASSETS_DIR=./build/assets
SRC_DIR=./cmd/observer
DIST_DIR=./build/dist

BINARY=observer
ifeq (${GOOS}, windows)
    BINARY := $(BINARY).exe
endif

COMMIT=$(shell git rev-parse --short HEAD)
TIMESTAMP=$(shell date +%s)

BUILD_FLAGS=-s -w \
	-X main.versionMajor=$(VERSION_MAJOR) \
	-X main.versionMinor=$(VERSION_MINOR) \
	-X main.versionPatch=$(VERSION_PATCH) \
	-X main.buildTimestamp=$(TIMESTAMP) \
	-X main.buildToolchain=${BUILD_TOOLCHAIN} \
	-X main.buildChannel=${BUILD_CHANNEL} \
	-X main.buildCommit=$(COMMIT)
BUILD_ARGS=-v -trimpath

build:
	@echo "[Info] Building project, output file path: $(DIST_DIR)/$(BINARY)"
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} GOMIPS=${GOMIPS} \
		$(GO) build -ldflags="$(BUILD_FLAGS)" $(BUILD_ARGS) -o $(DIST_DIR)/$(BINARY) $(SRC_DIR)/*.go
	@cp -r $(ASSETS_DIR) $(DIST_DIR)
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
	$(GO) run -gcflags="all=-N -l" -race $(SRC_DIR)/*.go --config $(DIST_DIR)/config.json.local

clean:
	@echo "[Warn] Cleaning up project..."
	@rm -rf $(DIST_DIR)/*

gen:
ifeq ($(shell command -v gqlgen 2> /dev/null),)
	@echo "[Info] Installing gqlgen..."
	@$(GO) get github.com/99designs/gqlgen
	@$(GO) install github.com/99designs/gqlgen
endif
	@echo "[Info] Generating GraphQL code..."
	@gqlgen generate

sync:
	@echo "not implemented"
