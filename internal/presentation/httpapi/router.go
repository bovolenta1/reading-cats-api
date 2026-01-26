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
	changeGoal         *ChangeGoalHandler
}

func NewRouter(me *MeHandler, readingHandler *RegisterReadingHandler, getReadingProgress *GetReadingProgressHandler, changeGoal *ChangeGoalHandler) *Router {
	return &Router{
		me:                 me,
		registerReading:    readingHandler,
		getReadingProgress: getReadingProgress,
		changeGoal:         changeGoal,
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

	if event.RequestContext.HTTP.Method == http.MethodPut && event.RawPath == "/v1/reading/goal" {
		return r.changeGoal.Handle(ctx, event)
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNotFound}, nil
}
