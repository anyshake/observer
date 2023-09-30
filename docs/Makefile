.PHONY: build

build:
ifeq ($(shell command -v swag 2> /dev/null),)
	$(error "Swagger is not installed. Install swag from github.com/swaggo/swag and try again.")
endif
	@swag init -g ../cmd/main.go -d ../app,../server,../publisher -o ./
