package httpapi

import (
	"context"
	"net/http"

	app "reading-cats-api/internal/application/reading"

	"github.com/aws/aws-lambda-go/events"
)

type ChangeGoalHandler struct {
	uc *app.ChangeGoalUseCase
}

func NewChangeGoalHandler(uc *app.ChangeGoalUseCase) *ChangeGoalHandler {
	return &ChangeGoalHandler{uc: uc}
}

func (h *ChangeGoalHandler) Handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	in, err := BuildChangeGoalInput(event)
	if err != nil {
		if err == ErrUnauthorized {
			return Error(event, http.StatusUnauthorized, err.Error()), nil
		}
		return Error(event, http.StatusBadRequest, err.Error()), nil
	}

	out, err := h.uc.Execute(ctx, in)
	if err != nil {
		return Error(event, http.StatusInternalServerError, err.Error()), nil
	}

	return JSON(http.StatusOK, out), nil
}
