.PHONY: build clean run

VERSION=$(shell cat ../../VERSION)
RELEASE=$(shell date +%Y%m%d%H%M%S)
COMMIT=$(shell git rev-parse --short HEAD)

DIST_DIR=../dist
BUILD_DIR=./build
GLOBAL_MAKEFILE=../../build/Makefile

build: clean pre-build
	@GENERATE_SOURCEMAP=false npm run build
	@rm -rf $(DIST_DIR)
	@mv $(BUILD_DIR) $(DIST_DIR)

run: pre-build
	@BROWSER=none npm run start

clean:
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@mkdir -p $(DIST_DIR)

pre-build: 
	@echo REACT_APP_VERSION=$(VERSION) > .env
	@echo REACT_APP_RELEASE=$(COMMIT)-$(RELEASE) >> .env
