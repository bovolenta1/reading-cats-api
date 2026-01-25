package httpapi

import (
	"context"
	"log"
	"net/http"

	app "reading-cats-api/internal/application/reading"

	"github.com/aws/aws-lambda-go/events"
)

type GetReadingProgressHandler struct {
	uc *app.GetReadingProgressUseCase
}

func NewGetReadingProgressHandler(uc *app.GetReadingProgressUseCase) *GetReadingProgressHandler {
	return &GetReadingProgressHandler{uc: uc}
}

func (h *GetReadingProgressHandler) Handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	in, err := BuildGetReadingProgressInput(event)
	if err != nil {
		return Error(http.StatusUnauthorized, err.Error()), nil
	}

	out, err := h.uc.Execute(ctx, in)
	if err != nil {
		log.Printf("reading.progress error request_id=%s err=%v", event.RequestContext.RequestID, err)
		return Error(http.StatusInternalServerError, "internal error"), nil
	}

	return JSON(http.StatusOK, map[string]any{"progress": out.Progress}), nil
}
