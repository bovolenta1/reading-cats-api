package httpapi

import (
	"context"
	"log"
	"net/http"

	appGroup "reading-cats-api/internal/application/group"

	"github.com/aws/aws-lambda-go/events"
)

type CreateGroupHandler struct {
	uc *appGroup.CreateGroupUseCase
}

func NewCreateGroupHandler(uc *appGroup.CreateGroupUseCase) *CreateGroupHandler {
	return &CreateGroupHandler{uc: uc}
}

func (h *CreateGroupHandler) Handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	in, err := BuildCreateGroupInput(event)
	if err != nil {
		if err.Error() == "unauthorized" {
			return Error(event, http.StatusUnauthorized, "unauthorized"), nil
		}
		return Error(event, http.StatusBadRequest, err.Error()), nil
	}

	out, err := h.uc.Execute(ctx, in)
	if err != nil {
		log.Printf("[httpapi] CreateGroup error: %v", err)
		return Error(event, http.StatusInternalServerError, err.Error()), nil
	}

	return JSON(http.StatusCreated, out), nil
}
