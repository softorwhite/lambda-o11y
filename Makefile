
.PHONY: build-lambda-o11y zip-lambda-o11y 

LAMBDA_CMD_DIR=cmd/lambda-o11y

build-lambda-o11y:
	cd app && GOOS=linux GOARCH=amd64 go build -o $(LAMBDA_CMD_DIR)/build/bootstrap $(LAMBDA_CMD_DIR)/main.go

zip-lambda-o11y: build-lambda-o11y
	zip -j app/$(LAMBDA_CMD_DIR)/build/bootstrap.zip app/$(LAMBDA_CMD_DIR)/build/bootstrap
