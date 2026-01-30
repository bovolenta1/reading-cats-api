package httpapi

import (
	"context"
	"log"
	"net/http"

	appSeason "reading-cats-api/internal/application/season"

	"github.com/aws/aws-lambda-go/events"
)

type CreateSeasonHandler struct {
	uc *appSeason.CreateSeasonUseCase
}

func NewCreateSeasonHandler(uc *appSeason.CreateSeasonUseCase) *CreateSeasonHandler {
	return &CreateSeasonHandler{uc: uc}
}

func (h *CreateSeasonHandler) Handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	in, err := BuildCreateSeasonInput(event)
	if err != nil {
		if err.Error() == "unauthorized" {
			return Error(event, http.StatusUnauthorized, "unauthorized"), nil
		}
		return Error(event, http.StatusBadRequest, err.Error()), nil
	}

	out, err := h.uc.Execute(ctx, in)
	if err != nil {
		log.Printf("[httpapi] CreateSeason error: %v", err)
		return Error(event, http.StatusInternalServerError, err.Error()), nil
	}

	return JSON(http.StatusCreated, out), nil
}
