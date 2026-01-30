package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type Router struct {
	me                 *MeHandler
	registerReading    *RegisterReadingHandler
	getReadingProgress *GetReadingProgressHandler
	changeGoal         *ChangeGoalHandler
	createGroup        *CreateGroupHandler
	createSeason       *CreateSeasonHandler
}

func NewRouter(me *MeHandler, readingHandler *RegisterReadingHandler, getReadingProgress *GetReadingProgressHandler, changeGoal *ChangeGoalHandler, createGroup *CreateGroupHandler, createSeason *CreateSeasonHandler) *Router {
	return &Router{
		me:                 me,
		registerReading:    readingHandler,
		getReadingProgress: getReadingProgress,
		changeGoal:         changeGoal,
		createGroup:        createGroup,
		createSeason:       createSeason,
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

	if event.RequestContext.HTTP.Method == http.MethodPost && event.RawPath == "/v1/groups" {
		return r.createGroup.Handle(ctx, event)
	}

	// POST /v1/groups/{groupId}/seasons
	if event.RequestContext.HTTP.Method == http.MethodPost && strings.HasPrefix(event.RawPath, "/v1/groups/") && strings.HasSuffix(event.RawPath, "/seasons") {
		return r.createSeason.Handle(ctx, event)
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNotFound}, nil
}
