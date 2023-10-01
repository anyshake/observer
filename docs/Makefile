.PHONY: build

build:
ifeq ($(shell command -v swag 2> /dev/null),)
	$(error "Swagger is not installed. Run go install github.com/swaggo/swag/cmd/swag@latest first")
endif
	@swag init -g ../cmd/main.go -d ../app,../server,../publisher -o ./
