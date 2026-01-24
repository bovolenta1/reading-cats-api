package httpapi

import (
	"context"
	"net/http"

	app "reading-cats-api/internal/application/user"

	"github.com/aws/aws-lambda-go/events"
)

type MeHandler struct {
	uc *app.EnsureMeUseCase
}

func NewMeHandler(uc *app.EnsureMeUseCase) *MeHandler {
	return &MeHandler{uc: uc}
}

func (h *MeHandler) Handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	in, err := BuildEnsureMeInput(event)
	if err != nil {
		// 401 se claims faltando, 400 se formato inválido
		return Error(http.StatusUnauthorized, err.Error()), nil
	}

	me, err := h.uc.Execute(ctx, in)
	if err != nil {
		return Error(http.StatusInternalServerError, err.Error()), nil
	}

	return JSON(http.StatusOK, me), nil
}
