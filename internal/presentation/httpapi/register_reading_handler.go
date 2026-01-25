package httpapi

import (
	"context"
	"net/http"

	app "reading-cats-api/internal/application/reading"

	"github.com/aws/aws-lambda-go/events"
)

type RegisterReadingHandler struct {
	uc *app.RegisterReadingUseCase
}

func NewRegisterReadingHandler(uc *app.RegisterReadingUseCase) *RegisterReadingHandler {
	return &RegisterReadingHandler{uc: uc}
}

func (h *RegisterReadingHandler) Handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	in, err := BuildRegisterReadingInput(event)
	if err != nil {
		if err == ErrUnauthorized {
			return Error(http.StatusUnauthorized, err.Error()), nil
		}
		return Error(http.StatusBadRequest, err.Error()), nil
	}

	out, err := h.uc.Execute(ctx, in)
	if err != nil {
		return Error(http.StatusBadRequest, err.Error()), nil
	}

	return JSON(http.StatusOK, map[string]any{"progress": out.Progress}), nil
}
