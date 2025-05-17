package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/softorwhite/lambda-o11y/app/handler"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
)

func main() {

	ctx := context.Background()

	tp, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		fmt.Printf("error creating tracer provider: %v", err)
	}

	defer func(ctx context.Context) {
		err := tp.Shutdown(ctx)
		if err != nil {
			fmt.Printf("error shutting down tracer provider: %v", err)
		}
	}(ctx)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	h := handler.NewHandler(ctx)

	lambda.Start(otellambda.InstrumentHandler(h.HandleRequest, xrayconfig.WithRecommendedOptions(tp)...))
}
