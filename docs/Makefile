.PHONY: gen

gen:
ifeq ($(shell command -v swag 2> /dev/null),)
	@echo "Installing Swagger..."
	@go get github.com/swaggo/swag/cmd/swag
	@go install github.com/swaggo/swag/cmd/swag
endif
	@swag init -g ../cmd/main.go -d ../api,../config,../drivers/explorer,../server -o ./
