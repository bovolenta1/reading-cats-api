package httpapi

import (
	app "reading-cats-api/internal/application/reading"

	"github.com/aws/aws-lambda-go/events"
)

func BuildGetReadingProgressInput(event events.APIGatewayV2HTTPRequest) (app.GetReadingProgressInput, error) {
	// Extract Claims from event
	claims, err := ExtractClaims(event)
	if err != nil {
		return app.GetReadingProgressInput{}, err
	}

	return app.GetReadingProgressInput{
		Claims: claims,
	}, nil
}
