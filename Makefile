CURRENT_VERSION_MAJOR = 4
CURRENT_VERSION_MINOR = 3
CURRENT_VERSION_PATCH = 3

REQUIRED_VERSION_MAJOR = 4
REQUIRED_VERSION_MINOR = 3
REQUIRED_VERSION_PATCH = 0

# Caution: software versioning mechanism depends on format of above lines in this file

.PHONY: build digest clean run gen version

GO ?= go

ASSETS_DIR = ./build/assets
SRC_DIR = ./cmd/observer
DIST_DIR = ./build/dist

BINARY = observer
ifeq (${GOOS}, windows)
    BINARY := $(BINARY).exe
endif

TIMESTAMP = $(shell date +%s)
COMMIT = $(shell git rev-parse --short HEAD)

BUILD_FLAGS = -s -w \
	-X main.versionMajor=$(CURRENT_VERSION_MAJOR) \
	-X main.versionMinor=$(CURRENT_VERSION_MINOR) \
	-X main.versionPatch=$(CURRENT_VERSION_PATCH) \
	-X main.versionPreRelease=${VERSION_PRE_RELEASE} \
	-X main.buildToolchain=${BUILD_TOOLCHAIN} \
	-X main.buildChannel=${BUILD_CHANNEL} \
	-X main.buildTimestamp=$(TIMESTAMP) \
	-X main.buildCommit=$(COMMIT)
BUILD_ARGS = -v -trimpath

build:
	@echo "[Info] Building project, output file path: $(DIST_DIR)/$(BINARY)"
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} GOMIPS=${GOMIPS} \
		$(GO) build -ldflags="$(BUILD_FLAGS)" $(BUILD_ARGS) -o $(DIST_DIR)/$(BINARY) $(SRC_DIR)
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

version:
	@echo -n 'latest_major=$(CURRENT_VERSION_MAJOR);'
	@echo -n 'latest_minor=$(CURRENT_VERSION_MINOR);'
	@echo -n 'latest_patch=$(CURRENT_VERSION_PATCH);'
	@echo -n 'required_major=$(REQUIRED_VERSION_MAJOR);'
	@echo -n 'required_minor=$(REQUIRED_VERSION_MINOR);'
	@echo -n 'required_patch=$(REQUIRED_VERSION_PATCH)'
