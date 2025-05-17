package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/softorwhite/lambda-o11y/app/application/usecase"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type HandlerFunc func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Handler struct {
	u func() *usecase.UserUseCase
}

func NewHandler(ctx context.Context) *Handler {
	// Initialize AWS config.
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	otelaws.AppendMiddlewares(&cfg.APIOptions)
	tracer = otel.Tracer("tracename1")

	return &Handler{
		u: func() *usecase.UserUseCase {
			return usecase.NewUserUseCase()
		},
	}
}

var tracer trace.Tracer

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
		),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/open-telemetry/opentelemetry-go/releases/latest", nil)
	if err != nil {
		fmt.Printf("failed to create http request, %v\n", err)
	}
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("failed to make http request, %v\n", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("failed to close http response body, %v\n", err)
		}
	}(res.Body)

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		fmt.Printf("failed to read http response body, %v\n", err)
	}

	fmt.Printf("Latest OTel Go Release is '%s'\n", data["name"])
	_, span := tracer.Start(
		ctx,
		"tracer_name",
	)

	// "Unfortunately there is no equivalent to the OpenTelemetry Event in the X-Ray data model so events are dropped when sending OTel spans to the X-ray backend"
	// https://github.com/aws-observability/aws-otel-collector/issues/821#issuecomment-1020768376
	span.AddEvent(
		"event1",
		trace.WithAttributes(
			attribute.String("key1", "value1"),
		),
	)
	defer span.End()

	user, err := h.u().GetUser("123")
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	jsonData, err := json.Marshal(user)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(jsonData), StatusCode: 200}, nil
}
