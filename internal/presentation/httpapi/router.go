package httpapi

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Router struct {
	me                 *MeHandler
	registerReading    *RegisterReadingHandler
	getReadingProgress *GetReadingProgressHandler
}

func NewRouter(me *MeHandler, readingHandler *RegisterReadingHandler, getReadingProgress *GetReadingProgressHandler) *Router {
	return &Router{
		me:                 me,
		registerReading:    readingHandler,
		getReadingProgress: getReadingProgress,
	}
}

func (r *Router) Route(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if event.RequestContext.HTTP.Method == http.MethodGet && event.RawPath == "/v1/me" {
		return r.me.Handle(ctx, event)
	}

	if event.RequestContext.HTTP.Method == http.MethodPost && event.RawPath == "/v1/reading/logs" {
		return r.registerReading.Handle(ctx, event)
	}

	if event.RequestContext.HTTP.Method == http.MethodGet && event.RawPath == "/v1/reading/progress" {
		return r.getReadingProgress.Handle(ctx, event)
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNotFound}, nil
}
