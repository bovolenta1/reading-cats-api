package httpapi

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Router struct {
	me *MeHandler
}

func NewRouter(me *MeHandler) *Router {
	return &Router{me: me}
}

func (r *Router) Route(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if event.RequestContext.HTTP.Method == http.MethodGet && event.RawPath == "/v1/me" {
		return r.me.Handle(ctx, event)
	}
	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNotFound}, nil
}
