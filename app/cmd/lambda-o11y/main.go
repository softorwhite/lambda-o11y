package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/softorwhite/lambda-o11y/app/handler"
)

func main() {

	h := handler.NewHandler()

	lambda.Start(h.HandleRequest)
}
