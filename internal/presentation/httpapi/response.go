package httpapi

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func JSON(status int, body any) events.APIGatewayV2HTTPResponse {
	b, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error":"failed to marshal json"}`,
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(b),
	}
}

func ErrorWithEvent(event events.APIGatewayV2HTTPRequest, status int, msg string) events.APIGatewayV2HTTPResponse {
	reqID := event.RequestContext.RequestID
	method := event.RequestContext.HTTP.Method
	path := event.RawPath

	log.Printf("[httpapi] Error reqId=%s status=%d method=%s path=%s msg=%s", reqID, status, method, path, msg)
	return JSON(status, map[string]string{"error": msg})
}
