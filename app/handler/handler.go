package handler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/softorwhite/lambda-o11y/app/application/usecase"
)

type HandlerFunc func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Handler struct {
	u func() *usecase.UserUseCase
}

func NewHandler() *Handler {
	return &Handler{
		u: func() *usecase.UserUseCase {
			return usecase.NewUserUseCase()
		},
	}
}
func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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
